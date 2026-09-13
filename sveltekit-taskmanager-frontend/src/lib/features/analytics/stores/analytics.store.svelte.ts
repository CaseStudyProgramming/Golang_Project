/**
 * Analytics store using Svelte 5 runes
 * Calculates task statistics and analytics
 */

import { taskStore } from '$lib/features/tasks'
import type { Task } from '$lib/features/tasks/types/task.types'

import type {
	AnalyticsData,
	AnalyticsState,
	CategoryDistribution,
	PriorityDistribution,
	ProductivityInsights,
	TaskStatistics,
	TimeBasedData,
	TimePeriod,
} from '../types/analytics.types'

/**
 * Create analytics store with Svelte 5 runes
 */
function createAnalyticsStore() {
	const state = $state<AnalyticsState>({
		data: null,
		error: null,
		isLoading: false,
		selectedPeriod: 'weekly',
	})

	/**
	 * Calculate task statistics
	 */
	function calculateStatistics(tasks: Task[]): TaskStatistics {
		const total = tasks.length
		const completed = tasks.filter((t) => t.status === 'completed').length
		const inProgress = tasks.filter((t) => t.status === 'in_progress').length
		const todo = tasks.filter((t) => t.status === 'todo').length
		const overdue = tasks.filter(
			(t) => t.dueDate && new Date(t.dueDate) < new Date() && t.status !== 'completed'
		).length
		const cancelled = tasks.filter((t) => t.status === 'cancelled').length
		const completionRate = total > 0 ? Math.round((completed / total) * 100) : 0

		return {
			cancelled,
			completed,
			completionRate,
			inProgress,
			overdue,
			todo,
			total,
		}
	}

	/**
	 * Calculate priority distribution
	 */
	function calculatePriorityDistribution(tasks: Task[]): PriorityDistribution {
		return {
			high: tasks.filter((t) => t.priority === 'high').length,
			low: tasks.filter((t) => t.priority === 'low').length,
			medium: tasks.filter((t) => t.priority === 'medium').length,
			urgent: tasks.filter((t) => t.priority === 'urgent').length,
		}
	}

	/**
	 * Calculate category distribution
	 */
	function calculateCategoryDistribution(tasks: Task[]): CategoryDistribution[] {
		const categoryMap = new Map<string, { completed: number; count: number }>()

		tasks.forEach((task) => {
			if (task.categoryId) {
				const existing = categoryMap.get(task.categoryId) || { completed: 0, count: 0 }
				categoryMap.set(task.categoryId, {
					completed: existing.completed + (task.status === 'completed' ? 1 : 0),
					count: existing.count + 1,
				})
			}
		})

		return Array.from(categoryMap.entries()).map(([categoryId, data]) => ({
			categoryId,
			categoryName: `Category ${categoryId}`, // Would be replaced with actual category name
			completed: data.completed,
			count: data.count,
		}))
	}

	/**
	 * Calculate time-based analytics
	 */
	function calculateTimeBasedData(tasks: Task[], period: TimePeriod): TimeBasedData[] {
		const now = new Date()
		const data: TimeBasedData[] = []
		const daysMap = new Map<string, { completed: number; created: number; total: number }>()

		tasks.forEach((task) => {
			const taskDate = new Date(task.createdAt)
			const dateKey = taskDate.toISOString().split('T')[0]

			const existing = daysMap.get(dateKey) || { completed: 0, created: 0, total: 0 }
			daysMap.set(dateKey, {
				completed: existing.completed + (task.status === 'completed' ? 1 : 0),
				created: existing.created + 1,
				total: existing.total + 1,
			})
		})

		// Generate data points based on period
		const daysToInclude = period === 'daily' ? 7 : period === 'weekly' ? 30 : 90

		for (let i = daysToInclude - 1; i >= 0; i--) {
			const date = new Date(now)
			date.setDate(date.getDate() - i)
			const dateKey = date.toISOString().split('T')[0]
			const dayData = daysMap.get(dateKey) || { completed: 0, created: 0, total: 0 }

			data.push({
				completed: dayData.completed,
				created: dayData.created,
				date: dateKey,
				total: dayData.total,
			})
		}

		return data
	}

	/**
	 * Calculate productivity insights
	 */
	function calculateProductivityInsights(tasks: Task[]): ProductivityInsights {
		const completedTasks = tasks.filter((t) => t.status === 'completed' && t.completedAt)
		const totalTasksCompleted = completedTasks.length

		// Calculate average completion time
		let totalCompletionTime = 0
		completedTasks.forEach((task) => {
			if (task.completedAt) {
				const created = new Date(task.createdAt).getTime()
				const completed = new Date(task.completedAt).getTime()
				totalCompletionTime += (completed - created) / (1000 * 60 * 60) // Convert to hours
			}
		})
		const averageCompletionTime =
			totalTasksCompleted > 0 ? totalCompletionTime / totalTasksCompleted : 0

		// Find most productive day
		const dayMap = new Map<string, number>()
		completedTasks.forEach((task) => {
			if (task.completedAt) {
				const day = new Date(task.completedAt).toLocaleDateString('en-US', { weekday: 'long' })
				dayMap.set(day, (dayMap.get(day) || 0) + 1)
			}
		})

		let mostProductiveDay = 'N/A'
		let maxCompletions = 0
		dayMap.forEach((count, day) => {
			if (count > maxCompletions) {
				maxCompletions = count
				mostProductiveDay = day
			}
		})

		// Calculate streak days
		const sortedCompletedTasks = completedTasks
			.filter((t) => t.completedAt)
			.sort((a, b) => new Date(b.completedAt!).getTime() - new Date(a.completedAt!).getTime())

		let streakDays = 0
		let currentDate = new Date()
		currentDate.setHours(0, 0, 0, 0)

		for (const task of sortedCompletedTasks) {
			if (task.completedAt) {
				const taskDate = new Date(task.completedAt)
				taskDate.setHours(0, 0, 0, 0)

				const diffDays = Math.floor(
					(currentDate.getTime() - taskDate.getTime()) / (1000 * 60 * 60 * 24)
				)

				if (diffDays === streakDays) {
					streakDays++
					currentDate = taskDate
				} else {
					break
				}
			}
		}

		// Calculate tasks per day
		const uniqueDays = new Set(completedTasks.map((t) => new Date(t.completedAt!).toDateString()))
			.size
		const tasksPerDay = uniqueDays > 0 ? totalTasksCompleted / uniqueDays : 0

		// Calculate on-time completion rate
		const onTimeCompleted = completedTasks.filter((task) => {
			if (task.dueDate && task.completedAt) {
				return new Date(task.completedAt) <= new Date(task.dueDate)
			}
			return true // Tasks without due dates are considered on-time
		}).length
		const onTimeCompletionRate =
			totalTasksCompleted > 0 ? Math.round((onTimeCompleted / totalTasksCompleted) * 100) : 0

		return {
			averageCompletionTime: Math.round(averageCompletionTime * 10) / 10,
			mostProductiveDay,
			onTimeCompletionRate,
			streakDays,
			tasksPerDay: Math.round(tasksPerDay * 10) / 10,
			totalTasksCompleted,
		}
	}

	/**
	 * Calculate all analytics data
	 */
	function calculateAnalytics(period: TimePeriod = 'weekly'): AnalyticsData {
		const tasks = taskStore.state.tasks

		return {
			categoryDistribution: calculateCategoryDistribution(tasks),
			period,
			priorityDistribution: calculatePriorityDistribution(tasks),
			productivityInsights: calculateProductivityInsights(tasks),
			statistics: calculateStatistics(tasks),
			timeBasedData: calculateTimeBasedData(tasks, period),
		}
	}

	/**
	 * Refresh analytics data
	 */
	async function refreshAnalytics(): Promise<void> {
		state.isLoading = true
		state.error = null

		try {
			// Ensure tasks are loaded
			if (taskStore.state.tasks.length === 0) {
				await taskStore.fetchTasks()
			}

			state.data = calculateAnalytics(state.selectedPeriod)
		} catch (error) {
			state.error = error instanceof Error ? error.message : 'Failed to load analytics'
			throw error
		} finally {
			state.isLoading = false
		}
	}

	/**
	 * Set time period
	 */
	function setPeriod(period: TimePeriod): void {
		state.selectedPeriod = period
		if (state.data) {
			state.data = calculateAnalytics(period)
		}
	}

	/**
	 * Reset store state
	 */
	function reset(): void {
		state.data = null
		state.isLoading = false
		state.error = null
		state.selectedPeriod = 'weekly'
	}

	return {
		refreshAnalytics,
		reset,
		setPeriod,
		get state() {
			return state
		},
	}
}

/**
 * Export analytics store instance
 * Only create store instance on client side to avoid SSR issues
 */
let analyticsStoreInstance: null | ReturnType<typeof createAnalyticsStore> = null

export const analyticsStore = new Proxy({} as ReturnType<typeof createAnalyticsStore>, {
	get(_target, prop) {
		if (!analyticsStoreInstance) {
			if (typeof window === 'undefined') {
				throw new Error('analyticsStore can only be accessed on the client side')
			}
			analyticsStoreInstance = createAnalyticsStore()
		}
		return analyticsStoreInstance[prop as keyof ReturnType<typeof createAnalyticsStore>]
	},
})

/**
 * Export store creator for testing
 */
export { createAnalyticsStore }

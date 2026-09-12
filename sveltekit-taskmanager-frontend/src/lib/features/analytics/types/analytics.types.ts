/**
 * Analytics feature types
 */

/**
 * Analytics data aggregate
 */
export interface AnalyticsData {
	statistics: TaskStatistics
	priorityDistribution: PriorityDistribution
	categoryDistribution: CategoryDistribution[]
	timeBasedData: TimeBasedData[]
	productivityInsights: ProductivityInsights
	period: TimePeriod
}

/**
 * Analytics store state
 */
export interface AnalyticsState {
	data: AnalyticsData | null
	isLoading: boolean
	error: null | string
	selectedPeriod: TimePeriod
}

/**
 * Category distribution data
 */
export interface CategoryDistribution {
	categoryId: string
	categoryName: string
	count: number
	completed: number
}

/**
 * Priority distribution data
 */
export interface PriorityDistribution {
	low: number
	medium: number
	high: number
	urgent: number
}

/**
 * Productivity insights
 */
export interface ProductivityInsights {
	totalTasksCompleted: number
	averageCompletionTime: number // in hours
	mostProductiveDay: string
	streakDays: number
	tasksPerDay: number
	onTimeCompletionRate: number
}

/**
 * Task statistics
 */
export interface TaskStatistics {
	total: number
	completed: number
	inProgress: number
	todo: number
	overdue: number
	cancelled: number
	completionRate: number
}

/**
 * Time-based analytics data
 */
export interface TimeBasedData {
	date: string
	completed: number
	created: number
	total: number
}

/**
 * Time period for analytics
 */
export type TimePeriod = 'daily' | 'monthly' | 'weekly'

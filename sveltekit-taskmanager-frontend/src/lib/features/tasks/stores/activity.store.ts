/**
 * Activity log store using Svelte 5 runes
 * Manages task activity history and filtering
 */

import { withErrorHandling } from '$lib/shared/utils/error.utils'

import type { Activity, ActivityFilters, ActivityState } from '../types/task.types'

/**
 * Create activity store with Svelte 5 runes
 */
function createActivityStore() {
	const state = $state<ActivityState>({
		activities: [],
		error: null,
		isLoading: false,
	})

	/**
	 * Fetch activities with filters
	 */
	async function fetchActivities(_filters?: ActivityFilters): Promise<void> {
		state.isLoading = true
		state.error = null

		try {
			await withErrorHandling(async () => {
				// This would be replaced with actual API call
				// const queryParams = new URLSearchParams();
				// if (filters?.taskId) queryParams.append('taskId', filters.taskId);
				// if (filters?.type) queryParams.append('type', filters.type);
				// if (filters?.userId) queryParams.append('userId', filters.userId);
				// if (filters?.dateFrom) queryParams.append('dateFrom', filters.dateFrom);
				// if (filters?.dateTo) queryParams.append('dateTo', filters.dateTo);
				// const response = await httpClient.get<Activity[]>(`/activities?${queryParams}`);

				// Mock response for development
				const mockActivities: Activity[] = []

				state.activities = mockActivities
			}, 'Failed to fetch activities')
		} catch (error) {
			state.error = error instanceof Error ? error.message : 'Failed to fetch activities'
			throw error
		} finally {
			state.isLoading = false
		}
	}

	/**
	 * Add new activity
	 */
	async function addActivity(activity: Omit<Activity, 'createdAt' | 'id'>): Promise<Activity> {
		try {
			const newActivity = await withErrorHandling(async () => {
				// This would be replaced with actual API call
				// const response = await httpClient.post<Activity>('/activities', activity);

				// Mock response for development
				const mockActivity: Activity = {
					...activity,
					createdAt: new Date().toISOString(),
					id: Date.now().toString(),
				}

				state.activities = [mockActivity, ...state.activities]

				return mockActivity
			}, 'Failed to add activity')

			return newActivity
		} catch (error) {
			state.error = error instanceof Error ? error.message : 'Failed to add activity'
			throw error
		}
	}

	/**
	 * Log task activity (convenience method)
	 */
	async function logTaskActivity(
		taskId: string,
		type: Activity['type'],
		description: string,
		changes?: Record<string, { new: unknown; old: unknown }>
	): Promise<void> {
		await addActivity({
			changes,
			description,
			taskId,
			type,
			userId: '1', // This would come from auth context
			userName: 'Current User', // This would come from auth context
		})
	}

	/**
	 * Clear error state
	 */
	function clearError(): void {
		state.error = null
	}

	/**
	 * Reset store state
	 */
	function reset(): void {
		state.activities = []
		state.isLoading = false
		state.error = null
	}

	return {
		addActivity,
		clearError,
		fetchActivities,
		logTaskActivity,
		reset,
		get state() {
			return state
		},
	}
}

/**
 * Export activity store instance
 * Only create store instance on client side to avoid SSR issues
 */
let activityStoreInstance: null | ReturnType<typeof createActivityStore> = null

export const activityStore = new Proxy({} as ReturnType<typeof createActivityStore>, {
	get(_target, prop) {
		if (!activityStoreInstance) {
			if (typeof window === 'undefined') {
				throw new Error('activityStore can only be accessed on the client side')
			}
			activityStoreInstance = createActivityStore()
		}
		return activityStoreInstance[prop as keyof ReturnType<typeof createActivityStore>]
	},
})

/**
 * Export store creator for testing
 */
export { createActivityStore }

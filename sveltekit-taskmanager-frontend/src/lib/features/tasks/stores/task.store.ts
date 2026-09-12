/**
 * Task management store using Svelte 5 runes
 * Manages tasks list, filters, and pagination
 */

import type { PaginatedResponse, PaginationParams } from '$lib/shared/types/api.types'

import { ValidationError, withErrorHandling } from '$lib/shared/utils/error.utils'
import { createTaskSchema, updateTaskSchema } from '../schemas/task.schemas'
import type { Subtask, Task, TaskFilters, TaskSort, TaskState } from '../types/task.types'

/**
 * Create task store with Svelte 5 runes
 */
function createTaskStore() {
	const state = $state<TaskState>({
		currentTask: null,
		error: null,
		filters: {},
		isLoading: false,
		pagination: {
			limit: 10,
			page: 1,
			total: 0,
			totalPages: 0,
		},
		sort: { field: 'createdAt', order: 'desc' },
		tasks: [],
	})

	/**
	 * Fetch tasks with filters and pagination
	 */
	async function fetchTasks(params?: PaginationParams): Promise<void> {
		state.isLoading = true
		state.error = null

		try {
			await withErrorHandling(async () => {
				const queryParams = new URLSearchParams()

				// Add pagination
				queryParams.append('page', (params?.page || state.pagination.page).toString())
				queryParams.append('limit', (params?.limit || state.pagination.limit).toString())

				// Add sort
				queryParams.append('sort', state.sort.field)
				queryParams.append('order', state.sort.order)

				// Add filters
				if (state.filters.status) queryParams.append('status', state.filters.status)
				if (state.filters.priority) queryParams.append('priority', state.filters.priority)
				if (state.filters.categoryId) queryParams.append('categoryId', state.filters.categoryId)
				if (state.filters.search) queryParams.append('search', state.filters.search)
				if (state.filters.dueDateFrom) queryParams.append('dueDateFrom', state.filters.dueDateFrom)
				if (state.filters.dueDateTo) queryParams.append('dueDateTo', state.filters.dueDateTo)
				if (state.filters.tags?.length) {
					state.filters.tags.forEach((tag) => {
						queryParams.append('tags', tag)
					})
				}

				// This would be replaced with actual API call
				// const response = await httpClient.get<PaginatedResponse<Task>>(`/tasks?${queryParams}`);

				// Mock response for development
				const mockResponse: PaginatedResponse<Task> = {
					data: [],
					limit: state.pagination.limit,
					page: state.pagination.page,
					total: 0,
					totalPages: 0,
				}

				state.tasks = mockResponse.data
				state.pagination = {
					limit: mockResponse.limit,
					page: mockResponse.page,
					total: mockResponse.total,
					totalPages: mockResponse.totalPages,
				}
			}, 'Failed to fetch tasks')
		} catch (error) {
			state.error = error instanceof Error ? error.message : 'Failed to fetch tasks'
			throw error
		} finally {
			state.isLoading = false
		}
	}

	/**
	 * Fetch single task by ID
	 */
	async function fetchTaskById(id: string): Promise<void> {
		state.isLoading = true
		state.error = null

		try {
			await withErrorHandling(async () => {
				// This would be replaced with actual API call
				// const response = await httpClient.get<Task>(`/tasks/${id}`);

				// Mock response for development
				const mockTask: Task = {
					createdAt: new Date().toISOString(),
					id,
					priority: 'medium',
					status: 'todo',
					title: 'Mock Task',
					updatedAt: new Date().toISOString(),
					userId: '1',
				}

				state.currentTask = mockTask
			}, 'Failed to fetch task')
		} catch (error) {
			state.error = error instanceof Error ? error.message : 'Failed to fetch task'
			throw error
		} finally {
			state.isLoading = false
		}
	}

	/**
	 * Create new task with optimistic UI update
	 */
	async function createTask(payload: unknown): Promise<Task> {
		state.error = null

		try {
			// Validate input with Zod
			const validatedPayload = createTaskSchema.parse(payload)

			// Optimistic update: create temporary task with generated ID
			const tempId = `temp-${Date.now()}`
			const optimisticTask: Task = {
				categoryId: validatedPayload.categoryId,
				createdAt: new Date().toISOString(),
				description: validatedPayload.description,
				dueDate: validatedPayload.dueDate,
				id: tempId,
				priority: validatedPayload.priority || 'medium',
				status: 'todo',
				tags: validatedPayload.tags,
				title: validatedPayload.title,
				updatedAt: new Date().toISOString(),
				userId: '1',
			}

			// Add to state immediately
			state.tasks = [optimisticTask, ...state.tasks]
			state.pagination.total += 1

			const newTask = await withErrorHandling(async () => {
				// This would be replaced with actual API call
				// const response = await httpClient.post<Task>('/tasks', validatedPayload);

				// Mock response for development
				const mockTask: Task = {
					...optimisticTask,
					id: Date.now().toString(), // Replace temp ID with real ID
				}

				// Replace optimistic task with real task
				state.tasks = state.tasks.map((task) => (task.id === tempId ? mockTask : task))

				return mockTask
			}, 'Failed to create task')

			return newTask
		} catch (error) {
			// Rollback optimistic update on error
			state.tasks = state.tasks.filter((task) => !task.id.startsWith('temp-'))
			state.pagination.total -= 1

			if (error instanceof Error && error.name === 'ZodError') {
				state.error = `Invalid input: ${error.message}`
				throw new ValidationError('task', error.message)
			}
			state.error = error instanceof Error ? error.message : 'Failed to create task'
			throw error
		}
	}

	/**
	 * Update existing task with optimistic UI update
	 */
	async function updateTask(id: string, payload: unknown): Promise<Task> {
		state.error = null

		// Store previous state for rollback
		const previousTask = state.tasks.find((t) => t.id === id)
		const previousCurrentTask = state.currentTask?.id === id ? { ...state.currentTask } : null

		try {
			// Validate input with Zod
			const validatedPayload = updateTaskSchema.parse(payload)

			// Optimistic update
			const optimisticTask: Task = {
				...previousTask!,
				...validatedPayload,
				updatedAt: new Date().toISOString(),
			}

			state.tasks = state.tasks.map((task) => (task.id === id ? optimisticTask : task))
			if (state.currentTask?.id === id) {
				state.currentTask = optimisticTask
			}

			const updatedTask = await withErrorHandling(async () => {
				// This would be replaced with actual API call
				// const response = await httpClient.patch<Task>(`/tasks/${id}`, validatedPayload);

				// Mock response for development
				const mockTask: Task = {
					...optimisticTask,
					updatedAt: new Date().toISOString(),
				}

				state.tasks = state.tasks.map((task) => (task.id === id ? mockTask : task))
				if (state.currentTask?.id === id) {
					state.currentTask = mockTask
				}

				return mockTask
			}, 'Failed to update task')

			return updatedTask
		} catch (error) {
			// Rollback optimistic update on error
			if (previousTask) {
				state.tasks = state.tasks.map((task) => (task.id === id ? { ...previousTask } : task))
			}
			if (previousCurrentTask) {
				state.currentTask = previousCurrentTask
			}

			if (error instanceof Error && error.name === 'ZodError') {
				state.error = `Invalid input: ${error.message}`
				throw new ValidationError('task', error.message)
			}
			state.error = error instanceof Error ? error.message : 'Failed to update task'
			throw error
		}
	}

	/**
	 * Delete task (soft delete) with optimistic UI update
	 */
	async function deleteTask(id: string): Promise<void> {
		state.error = null

		// Store previous state for rollback
		const previousTask = state.tasks.find((t) => t.id === id)
		const previousCurrentTask = state.currentTask?.id === id ? { ...state.currentTask } : null

		try {
			// Optimistic update: soft delete by updating status
			state.tasks = state.tasks.map((task) =>
				task.id === id
					? { ...task, status: 'deleted' as const, updatedAt: new Date().toISOString() }
					: task
			)
			state.pagination.total -= 1

			if (state.currentTask?.id === id) {
				state.currentTask = null
			}

			await withErrorHandling(async () => {
				// This would be replaced with actual API call for soft delete
				// await httpClient.patch(`/tasks/${id}`, { status: 'deleted' });
			}, 'Failed to delete task')
		} catch (error) {
			// Rollback optimistic update on error
			if (previousTask) {
				state.tasks = state.tasks.map((task) => (task.id === id ? { ...previousTask } : task))
				state.pagination.total += 1
			}
			if (previousCurrentTask) {
				state.currentTask = previousCurrentTask
			}

			state.error = error instanceof Error ? error.message : 'Failed to delete task'
			throw error
		}
	}

	/**
	 * Restore deleted task
	 */
	async function restoreTask(id: string): Promise<void> {
		state.isLoading = true
		state.error = null

		try {
			await withErrorHandling(async () => {
				// This would be replaced with actual API call
				// await httpClient.patch(`/tasks/${id}`, { status: 'todo' });

				state.tasks = state.tasks.map((task) =>
					task.id === id
						? { ...task, status: 'todo' as const, updatedAt: new Date().toISOString() }
						: task
				)
				state.pagination.total += 1
			}, 'Failed to restore task')
		} catch (error) {
			state.error = error instanceof Error ? error.message : 'Failed to restore task'
			throw error
		} finally {
			state.isLoading = false
		}
	}

	/**
	 * Permanently delete task
	 */
	async function permanentDeleteTask(id: string): Promise<void> {
		state.isLoading = true
		state.error = null

		try {
			await withErrorHandling(async () => {
				// This would be replaced with actual API call
				// await httpClient.delete(`/tasks/${id}/permanent`);

				state.tasks = state.tasks.filter((task) => task.id !== id)
				if (state.currentTask?.id === id) {
					state.currentTask = null
				}
			}, 'Failed to permanently delete task')
		} catch (error) {
			state.error = error instanceof Error ? error.message : 'Failed to permanently delete task'
			throw error
		} finally {
			state.isLoading = false
		}
	}

	/**
	 * Update filters
	 */
	function setFilters(filters: Partial<TaskFilters>): void {
		state.filters = { ...state.filters, ...filters }
		state.pagination.page = 1 // Reset to first page when filters change
	}

	/**
	 * Clear all filters
	 */
	function clearFilters(): void {
		state.filters = {}
		state.pagination.page = 1
	}

	/**
	 * Update sort
	 */
	function setSort(sort: TaskSort): void {
		state.sort = sort
	}

	/**
	 * Change pagination page
	 */
	function setPage(page: number): void {
		state.pagination.page = page
	}

	/**
	 * Change pagination limit
	 */
	function setLimit(limit: number): void {
		state.pagination.limit = limit
		state.pagination.page = 1 // Reset to first page when limit changes
	}

	/**
	 * Clear error state
	 */
	function clearError(): void {
		state.error = null
	}

	/**
	 * Add subtask to task
	 */
	async function addSubtask(taskId: string, title: string): Promise<Subtask> {
		try {
			const subtask = await withErrorHandling(async () => {
				// This would be replaced with actual API call
				// const response = await httpClient.post<Subtask>(`/tasks/${taskId}/subtasks`, { title });

				// Mock response for development
				const mockSubtask: Subtask = {
					createdAt: new Date().toISOString(),
					id: Date.now().toString(),
					isCompleted: false,
					taskId,
					title: title.trim(),
					updatedAt: new Date().toISOString(),
				}

				// Update task's subtasks and progress
				state.tasks = state.tasks.map((task) => {
					if (task.id === taskId) {
						const updatedSubtasks = [...(task.subtasks || []), mockSubtask]
						const progress = calculateProgress(updatedSubtasks)
						return { ...task, progress, subtasks: updatedSubtasks }
					}
					return task
				})

				if (state.currentTask?.id === taskId) {
					const updatedSubtasks = [...(state.currentTask.subtasks || []), mockSubtask]
					const progress = calculateProgress(updatedSubtasks)
					state.currentTask = { ...state.currentTask, progress, subtasks: updatedSubtasks }
				}

				return mockSubtask
			}, 'Failed to add subtask')

			return subtask
		} catch (error) {
			state.error = error instanceof Error ? error.message : 'Failed to add subtask'
			throw error
		}
	}

	/**
	 * Toggle subtask completion with optimistic UI update
	 */
	async function toggleSubtask(taskId: string, subtaskId: string): Promise<void> {
		// Store previous state for rollback
		const previousTask = state.tasks.find((t) => t.id === taskId)
		const previousCurrentTask = state.currentTask?.id === taskId ? { ...state.currentTask } : null

		try {
			// Optimistic update
			state.tasks = state.tasks.map((task) => {
				if (task.id === taskId) {
					const updatedSubtasks =
						task.subtasks?.map((subtask) =>
							subtask.id === subtaskId
								? {
										...subtask,
										isCompleted: !subtask.isCompleted,
										updatedAt: new Date().toISOString(),
									}
								: subtask
						) || []
					const progress = calculateProgress(updatedSubtasks)
					return { ...task, progress, subtasks: updatedSubtasks }
				}
				return task
			})

			if (state.currentTask?.id === taskId) {
				const updatedSubtasks =
					state.currentTask.subtasks?.map((subtask) =>
						subtask.id === subtaskId
							? {
									...subtask,
									isCompleted: !subtask.isCompleted,
									updatedAt: new Date().toISOString(),
								}
							: subtask
					) || []
				const progress = calculateProgress(updatedSubtasks)
				state.currentTask = { ...state.currentTask, progress, subtasks: updatedSubtasks }
			}

			await withErrorHandling(async () => {
				// This would be replaced with actual API call
				// await httpClient.patch(`/tasks/${taskId}/subtasks/${subtaskId}`, { isCompleted: !isCompleted });
			}, 'Failed to toggle subtask')
		} catch (error) {
			// Rollback optimistic update on error
			if (previousTask) {
				state.tasks = state.tasks.map((task) => (task.id === taskId ? { ...previousTask } : task))
			}
			if (previousCurrentTask) {
				state.currentTask = previousCurrentTask
			}

			state.error = error instanceof Error ? error.message : 'Failed to toggle subtask'
			throw error
		}
	}

	/**
	 * Delete subtask
	 */
	async function deleteSubtask(taskId: string, subtaskId: string): Promise<void> {
		try {
			await withErrorHandling(async () => {
				// This would be replaced with actual API call
				// await httpClient.delete(`/tasks/${taskId}/subtasks/${subtaskId}`);

				// Update task's subtasks and progress
				state.tasks = state.tasks.map((task) => {
					if (task.id === taskId) {
						const updatedSubtasks =
							task.subtasks?.filter((subtask) => subtask.id !== subtaskId) || []
						const progress = calculateProgress(updatedSubtasks)
						return { ...task, progress, subtasks: updatedSubtasks }
					}
					return task
				})

				if (state.currentTask?.id === taskId) {
					const updatedSubtasks =
						state.currentTask.subtasks?.filter((subtask) => subtask.id !== subtaskId) || []
					const progress = calculateProgress(updatedSubtasks)
					state.currentTask = { ...state.currentTask, progress, subtasks: updatedSubtasks }
				}
			}, 'Failed to delete subtask')
		} catch (error) {
			state.error = error instanceof Error ? error.message : 'Failed to delete subtask'
			throw error
		}
	}

	/**
	 * Bulk complete subtasks
	 */
	async function bulkCompleteSubtasks(taskId: string, subtaskIds: string[]): Promise<void> {
		try {
			await withErrorHandling(async () => {
				// This would be replaced with actual API call
				// await httpClient.post(`/tasks/${taskId}/subtasks/bulk-complete`, { subtaskIds });

				// Update task's subtasks and progress
				state.tasks = state.tasks.map((task) => {
					if (task.id === taskId) {
						const updatedSubtasks =
							task.subtasks?.map((subtask) =>
								subtaskIds.includes(subtask.id)
									? { ...subtask, isCompleted: true, updatedAt: new Date().toISOString() }
									: subtask
							) || []
						const progress = calculateProgress(updatedSubtasks)
						return { ...task, progress, subtasks: updatedSubtasks }
					}
					return task
				})

				if (state.currentTask?.id === taskId) {
					const updatedSubtasks =
						state.currentTask.subtasks?.map((subtask) =>
							subtaskIds.includes(subtask.id)
								? { ...subtask, isCompleted: true, updatedAt: new Date().toISOString() }
								: subtask
						) || []
					const progress = calculateProgress(updatedSubtasks)
					state.currentTask = { ...state.currentTask, progress, subtasks: updatedSubtasks }
				}
			}, 'Failed to bulk complete subtasks')
		} catch (error) {
			state.error = error instanceof Error ? error.message : 'Failed to bulk complete subtasks'
			throw error
		}
	}

	/**
	 * Bulk delete subtasks
	 */
	async function bulkDeleteSubtasks(taskId: string, subtaskIds: string[]): Promise<void> {
		try {
			await withErrorHandling(async () => {
				// This would be replaced with actual API call
				// await httpClient.post(`/tasks/${taskId}/subtasks/bulk-delete`, { subtaskIds });

				// Update task's subtasks and progress
				state.tasks = state.tasks.map((task) => {
					if (task.id === taskId) {
						const updatedSubtasks =
							task.subtasks?.filter((subtask) => !subtaskIds.includes(subtask.id)) || []
						const progress = calculateProgress(updatedSubtasks)
						return { ...task, progress, subtasks: updatedSubtasks }
					}
					return task
				})

				if (state.currentTask?.id === taskId) {
					const updatedSubtasks =
						state.currentTask.subtasks?.filter((subtask) => !subtaskIds.includes(subtask.id)) || []
					const progress = calculateProgress(updatedSubtasks)
					state.currentTask = { ...state.currentTask, progress, subtasks: updatedSubtasks }
				}
			}, 'Failed to bulk delete subtasks')
		} catch (error) {
			state.error = error instanceof Error ? error.message : 'Failed to bulk delete subtasks'
			throw error
		}
	}

	/**
	 * Calculate progress percentage based on completed subtasks
	 */
	function calculateProgress(subtasks: Subtask[]): number {
		if (subtasks.length === 0) return 0
		const completed = subtasks.filter((s) => s.isCompleted).length
		return Math.round((completed / subtasks.length) * 100)
	}

	/**
	 * Reset store state
	 */
	function reset(): void {
		state.tasks = []
		state.currentTask = null
		state.filters = {}
		state.sort = { field: 'createdAt', order: 'desc' }
		state.pagination = {
			limit: 10,
			page: 1,
			total: 0,
			totalPages: 0,
		}
		state.isLoading = false
		state.error = null
	}

	return {
		addSubtask,
		bulkCompleteSubtasks,
		bulkDeleteSubtasks,
		clearError,
		clearFilters,
		createTask,
		deleteSubtask,
		deleteTask,
		fetchTaskById,
		fetchTasks,
		permanentDeleteTask,
		reset,
		restoreTask,
		setFilters,
		setLimit,
		setPage,
		setSort,
		get state() {
			return state
		},
		toggleSubtask,
		updateTask,
	}
}

/**
 * Export task store instance
 * Only create store instance on client side to avoid SSR issues
 */
let taskStoreInstance: null | ReturnType<typeof createTaskStore> = null

export const taskStore = new Proxy({} as ReturnType<typeof createTaskStore>, {
	get(_target, prop) {
		if (!taskStoreInstance) {
			if (typeof window === 'undefined') {
				throw new Error('taskStore can only be accessed on the client side')
			}
			taskStoreInstance = createTaskStore()
		}
		return taskStoreInstance[prop as keyof ReturnType<typeof createTaskStore>]
	},
})

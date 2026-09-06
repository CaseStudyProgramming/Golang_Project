/**
 * Task management store using Svelte 5 runes
 * Manages tasks list, filters, and pagination
 */

import { httpClient } from '$lib/shared/utils/api.utils';
import { withErrorHandling, ValidationError } from '$lib/shared/utils/error.utils';
import { createTaskSchema, updateTaskSchema } from '../schemas/task.schemas';
import type { Task, TaskFilters, TaskSort, TaskState } from '../types/task.types';
import type { PaginatedResponse, PaginationParams } from '$lib/shared/types/api.types';

/**
 * Create task store with Svelte 5 runes
 */
function createTaskStore() {
	const state = $state<TaskState>({
		tasks: [],
		currentTask: null,
		filters: {},
		sort: { field: 'createdAt', order: 'desc' },
		pagination: {
			page: 1,
			limit: 10,
			total: 0,
			totalPages: 0
		},
		isLoading: false,
		error: null
	});

	/**
	 * Fetch tasks with filters and pagination
	 */
	async function fetchTasks(params?: PaginationParams): Promise<void> {
		state.isLoading = true;
		state.error = null;

		try {
			await withErrorHandling(async () => {
				const queryParams = new URLSearchParams();
				
				// Add pagination
				queryParams.append('page', (params?.page || state.pagination.page).toString());
				queryParams.append('limit', (params?.limit || state.pagination.limit).toString());
				
				// Add sort
				queryParams.append('sort', state.sort.field);
				queryParams.append('order', state.sort.order);
				
				// Add filters
				if (state.filters.status) queryParams.append('status', state.filters.status);
				if (state.filters.priority) queryParams.append('priority', state.filters.priority);
				if (state.filters.categoryId) queryParams.append('categoryId', state.filters.categoryId);
				if (state.filters.search) queryParams.append('search', state.filters.search);
				if (state.filters.dueDateFrom) queryParams.append('dueDateFrom', state.filters.dueDateFrom);
				if (state.filters.dueDateTo) queryParams.append('dueDateTo', state.filters.dueDateTo);
				if (state.filters.tags?.length) {
					state.filters.tags.forEach((tag) => queryParams.append('tags', tag));
				}

				// This would be replaced with actual API call
				// const response = await httpClient.get<PaginatedResponse<Task>>(`/tasks?${queryParams}`);
				
				// Mock response for development
				const mockResponse: PaginatedResponse<Task> = {
					data: [],
					total: 0,
					page: state.pagination.page,
					limit: state.pagination.limit,
					totalPages: 0
				};

				state.tasks = mockResponse.data;
				state.pagination = {
					page: mockResponse.page,
					limit: mockResponse.limit,
					total: mockResponse.total,
					totalPages: mockResponse.totalPages
				};
			}, 'Failed to fetch tasks');
		} catch (error) {
			state.error = error instanceof Error ? error.message : 'Failed to fetch tasks';
			throw error;
		} finally {
			state.isLoading = false;
		}
	}

	/**
	 * Fetch single task by ID
	 */
	async function fetchTaskById(id: string): Promise<void> {
		state.isLoading = true;
		state.error = null;

		try {
			await withErrorHandling(async () => {
				// This would be replaced with actual API call
				// const response = await httpClient.get<Task>(`/tasks/${id}`);
				
				// Mock response for development
				const mockTask: Task = {
					id,
					title: 'Mock Task',
					status: 'todo',
					priority: 'medium',
					createdAt: new Date().toISOString(),
					updatedAt: new Date().toISOString(),
					userId: '1'
				};

				state.currentTask = mockTask;
			}, 'Failed to fetch task');
		} catch (error) {
			state.error = error instanceof Error ? error.message : 'Failed to fetch task';
			throw error;
		} finally {
			state.isLoading = false;
		}
	}

	/**
	 * Create new task
	 */
	async function createTask(payload: unknown): Promise<Task> {
		state.isLoading = true;
		state.error = null;

		try {
			// Validate input with Zod
			const validatedPayload = createTaskSchema.parse(payload);

			const newTask = await withErrorHandling(async () => {
				// This would be replaced with actual API call
				// const response = await httpClient.post<Task>('/tasks', validatedPayload);
				
				// Mock response for development
				const mockTask: Task = {
					id: Date.now().toString(),
					title: validatedPayload.title,
					description: validatedPayload.description,
					status: 'todo',
					priority: validatedPayload.priority || 'medium',
					dueDate: validatedPayload.dueDate,
					categoryId: validatedPayload.categoryId,
					tags: validatedPayload.tags,
					createdAt: new Date().toISOString(),
					updatedAt: new Date().toISOString(),
					userId: '1'
				};

				state.tasks = [mockTask, ...state.tasks];
				state.pagination.total += 1;

				return mockTask;
			}, 'Failed to create task');

			return newTask;
		} catch (error) {
			if (error instanceof Error && error.name === 'ZodError') {
				state.error = 'Invalid input: ' + error.message;
				throw new ValidationError('task', error.message);
			}
			state.error = error instanceof Error ? error.message : 'Failed to create task';
			throw error;
		} finally {
			state.isLoading = false;
		}
	}

	/**
	 * Update existing task
	 */
	async function updateTask(id: string, payload: unknown): Promise<Task> {
		state.isLoading = true;
		state.error = null;

		try {
			// Validate input with Zod
			const validatedPayload = updateTaskSchema.parse(payload);

			const updatedTask = await withErrorHandling(async () => {
				// This would be replaced with actual API call
				// const response = await httpClient.patch<Task>(`/tasks/${id}`, validatedPayload);
				
				// Mock response for development
				const mockTask: Task = {
					...state.tasks.find((t) => t.id === id)!,
					...validatedPayload,
					updatedAt: new Date().toISOString()
				};

				state.tasks = state.tasks.map((task) => (task.id === id ? mockTask : task));
				if (state.currentTask?.id === id) {
					state.currentTask = mockTask;
				}

				return mockTask;
			}, 'Failed to update task');

			return updatedTask;
		} catch (error) {
			if (error instanceof Error && error.name === 'ZodError') {
				state.error = 'Invalid input: ' + error.message;
				throw new ValidationError('task', error.message);
			}
			state.error = error instanceof Error ? error.message : 'Failed to update task';
			throw error;
		} finally {
			state.isLoading = false;
		}
	}

	/**
	 * Delete task
	 */
	async function deleteTask(id: string): Promise<void> {
		state.isLoading = true;
		state.error = null;

		try {
			await withErrorHandling(async () => {
				// This would be replaced with actual API call
				// await httpClient.delete(`/tasks/${id}`);
				
				state.tasks = state.tasks.filter((task) => task.id !== id);
				state.pagination.total -= 1;

				if (state.currentTask?.id === id) {
					state.currentTask = null;
				}
			}, 'Failed to delete task');
		} catch (error) {
			state.error = error instanceof Error ? error.message : 'Failed to delete task';
			throw error;
		} finally {
			state.isLoading = false;
		}
	}

	/**
	 * Update filters
	 */
	function setFilters(filters: Partial<TaskFilters>): void {
		state.filters = { ...state.filters, ...filters };
		state.pagination.page = 1; // Reset to first page when filters change
	}

	/**
	 * Clear all filters
	 */
	function clearFilters(): void {
		state.filters = {};
		state.pagination.page = 1;
	}

	/**
	 * Update sort
	 */
	function setSort(sort: TaskSort): void {
		state.sort = sort;
	}

	/**
	 * Change pagination page
	 */
	function setPage(page: number): void {
		state.pagination.page = page;
	}

	/**
	 * Change pagination limit
	 */
	function setLimit(limit: number): void {
		state.pagination.limit = limit;
		state.pagination.page = 1; // Reset to first page when limit changes
	}

	/**
	 * Clear error state
	 */
	function clearError(): void {
		state.error = null;
	}

	/**
	 * Reset store state
	 */
	function reset(): void {
		state.tasks = [];
		state.currentTask = null;
		state.filters = {};
		state.sort = { field: 'createdAt', order: 'desc' };
		state.pagination = {
			page: 1,
			limit: 10,
			total: 0,
			totalPages: 0
		};
		state.isLoading = false;
		state.error = null;
	}

	return {
		get state() {
			return state;
		},
		fetchTasks,
		fetchTaskById,
		createTask,
		updateTask,
		deleteTask,
		setFilters,
		clearFilters,
		setSort,
		setPage,
		setLimit,
		clearError,
		reset
	};
}

/**
 * Export task store instance
 * Only create store instance on client side to avoid SSR issues
 */
let taskStoreInstance: ReturnType<typeof createTaskStore> | null = null;

export const taskStore = new Proxy({} as ReturnType<typeof createTaskStore>, {
	get(_target, prop) {
		if (!taskStoreInstance) {
			if (typeof window === 'undefined') {
				throw new Error('taskStore can only be accessed on the client side');
			}
			taskStoreInstance = createTaskStore();
		}
		return taskStoreInstance[prop as keyof ReturnType<typeof createTaskStore>];
	}
});

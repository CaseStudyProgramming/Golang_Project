/**
 * Task management store using Svelte 5 runes
 * Manages tasks list, filters, and pagination
 */

import { httpClient } from '../utils/api.utils';
import type { Task, TaskFilters, TaskSort, CreateTaskPayload, UpdateTaskPayload } from '../types/task.types';
import type { PaginatedResponse, PaginationParams } from '../types/api.types';

/**
 * Task store state interface
 */
interface TaskState {
	tasks: Task[];
	currentTask: Task | null;
	filters: TaskFilters;
	sort: TaskSort;
	pagination: {
		page: number;
		limit: number;
		total: number;
		totalPages: number;
	};
	isLoading: boolean;
	error: string | null;
}

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
	async function createTask(payload: CreateTaskPayload): Promise<Task> {
		state.isLoading = true;
		state.error = null;

		try {
			// This would be replaced with actual API call
			// const response = await httpClient.post<Task>('/tasks', payload);
			
			// Mock response for development
			const newTask: Task = {
				id: Date.now().toString(),
				title: payload.title,
				description: payload.description,
				status: 'todo',
				priority: payload.priority || 'medium',
				dueDate: payload.dueDate,
				categoryId: payload.categoryId,
				tags: payload.tags,
				createdAt: new Date().toISOString(),
				updatedAt: new Date().toISOString(),
				userId: '1'
			};

			state.tasks = [newTask, ...state.tasks];
			state.pagination.total += 1;

			return newTask;
		} catch (error) {
			state.error = error instanceof Error ? error.message : 'Failed to create task';
			throw error;
		} finally {
			state.isLoading = false;
		}
	}

	/**
	 * Update existing task
	 */
	async function updateTask(id: string, payload: UpdateTaskPayload): Promise<Task> {
		state.isLoading = true;
		state.error = null;

		try {
			// This would be replaced with actual API call
			// const response = await httpClient.patch<Task>(`/tasks/${id}`, payload);
			
			// Mock response for development
			const updatedTask: Task = {
				...state.tasks.find((t) => t.id === id)!,
				...payload,
				updatedAt: new Date().toISOString()
			};

			state.tasks = state.tasks.map((task) => (task.id === id ? updatedTask : task));
			if (state.currentTask?.id === id) {
				state.currentTask = updatedTask;
			}

			return updatedTask;
		} catch (error) {
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
			// This would be replaced with actual API call
			// await httpClient.delete(`/tasks/${id}`);
			
			state.tasks = state.tasks.filter((task) => task.id !== id);
			state.pagination.total -= 1;

			if (state.currentTask?.id === id) {
				state.currentTask = null;
			}
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
 */
export const taskStore = createTaskStore();

/**
 * Task feature types
 */

/**
 * Task priority levels
 */
export type TaskPriority = 'low' | 'medium' | 'high' | 'urgent';

/**
 * Task status
 */
export type TaskStatus = 'todo' | 'in_progress' | 'completed' | 'cancelled';

/**
 * Task interface
 */
export interface Task {
	id: string;
	title: string;
	description?: string;
	status: TaskStatus;
	priority: TaskPriority;
	dueDate?: string;
	completedAt?: string;
	createdAt: string;
	updatedAt: string;
	userId: string;
	categoryId?: string;
	tags?: string[];
}

/**
 * Task filter options
 */
export interface TaskFilters {
	status?: TaskStatus;
	priority?: TaskPriority;
	categoryId?: string;
	search?: string;
	dueDateFrom?: string;
	dueDateTo?: string;
	tags?: string[];
}

/**
 * Task sort options
 */
export interface TaskSort {
	field: 'title' | 'dueDate' | 'priority' | 'createdAt' | 'updatedAt';
	order: 'asc' | 'desc';
}

/**
 * Task creation payload
 */
export interface CreateTaskPayload {
	title: string;
	description?: string;
	priority?: TaskPriority;
	dueDate?: string;
	categoryId?: string;
	tags?: string[];
}

/**
 * Task update payload
 */
export interface UpdateTaskPayload {
	title?: string;
	description?: string;
	status?: TaskStatus;
	priority?: TaskPriority;
	dueDate?: string;
	categoryId?: string;
	tags?: string[];
}

/**
 * Task store state interface
 */
export interface TaskState {
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

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
export type TaskStatus = 'todo' | 'in_progress' | 'completed' | 'cancelled' | 'deleted';

/**
 * Subtask interface
 */
export interface Subtask {
	id: string;
	taskId: string;
	title: string;
	isCompleted: boolean;
	createdAt: string;
	updatedAt: string;
}

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
	subtasks?: Subtask[];
	progress?: number;
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
 * Activity types
 */
export type ActivityType =
	| 'task_created'
	| 'task_updated'
	| 'task_deleted'
	| 'task_completed'
	| 'subtask_added'
	| 'subtask_completed'
	| 'subtask_deleted'
	| 'category_assigned'
	| 'tag_added'
	| 'tag_removed'
	| 'status_changed'
	| 'priority_changed';

/**
 * Activity log interface
 */
export interface Activity {
	id: string;
	taskId: string;
	type: ActivityType;
	description: string;
	userId: string;
	userName?: string;
	changes?: Record<string, { old: unknown; new: unknown }>;
	createdAt: string;
}

/**
 * Activity filter options
 */
export interface ActivityFilters {
	taskId?: string;
	type?: ActivityType;
	userId?: string;
	dateFrom?: string;
	dateTo?: string;
}

/**
 * Activity store state interface
 */
export interface ActivityState {
	activities: Activity[];
	isLoading: boolean;
	error: string | null;
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

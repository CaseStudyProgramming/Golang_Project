/**
 * Task feature types
 */

/**
 * Activity log interface
 */
export interface Activity {
	id: string
	taskId: string
	type: ActivityType
	description: string
	userId: string
	userName?: string
	changes?: Record<string, { new: unknown; old: unknown }>
	createdAt: string
}

/**
 * Activity filter options
 */
export interface ActivityFilters {
	taskId?: string
	type?: ActivityType
	userId?: string
	dateFrom?: string
	dateTo?: string
}

/**
 * Activity store state interface
 */
export interface ActivityState {
	activities: Activity[]
	isLoading: boolean
	error: null | string
}

/**
 * Activity types
 */
export type ActivityType =
	| 'category_assigned'
	| 'priority_changed'
	| 'status_changed'
	| 'subtask_added'
	| 'subtask_completed'
	| 'subtask_deleted'
	| 'tag_added'
	| 'tag_removed'
	| 'task_completed'
	| 'task_created'
	| 'task_deleted'
	| 'task_updated'

/**
 * Subtask interface
 */
export interface Subtask {
	id: string
	taskId: string
	title: string
	isCompleted: boolean
	createdAt: string
	updatedAt: string
}

/**
 * Task interface
 */
export interface Task {
	id: string
	title: string
	description?: string
	status: TaskStatus
	priority: TaskPriority
	dueDate?: string
	completedAt?: string
	createdAt: string
	updatedAt: string
	userId: string
	categoryId?: string
	tags?: string[]
	subtasks?: Subtask[]
	progress?: number
}

/**
 * Task filter options
 */
export interface TaskFilters {
	status?: TaskStatus
	priority?: TaskPriority
	categoryId?: string
	search?: string
	dueDateFrom?: string
	dueDateTo?: string
	tags?: string[]
}

/**
 * Task priority levels
 */
export type TaskPriority = 'high' | 'low' | 'medium' | 'urgent'

/**
 * Task sort options
 */
export interface TaskSort {
	field: 'createdAt' | 'dueDate' | 'priority' | 'title' | 'updatedAt'
	order: 'asc' | 'desc'
}

/**
 * Task store state interface
 */
export interface TaskState {
	tasks: Task[]
	currentTask: null | Task
	filters: TaskFilters
	sort: TaskSort
	pagination: {
		limit: number
		page: number
		total: number
		totalPages: number
	}
	isLoading: boolean
	error: null | string
}

/**
 * Task status
 */
export type TaskStatus = 'cancelled' | 'completed' | 'deleted' | 'in_progress' | 'todo'

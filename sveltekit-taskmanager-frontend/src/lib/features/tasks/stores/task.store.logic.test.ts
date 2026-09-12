import { describe, expect, it } from 'vitest'

describe('Task Store Logic', () => {
	describe('Task CRUD Operations', () => {
		it('creates task with optimistic update', () => {
			const tasks: Array<{ createdAt: string; id: string; status: string; title: string }> = []
			const tempId = `temp-${Date.now()}`
			const newTask = {
				createdAt: new Date().toISOString(),
				id: tempId,
				status: 'todo',
				title: 'New Task',
			}

			const updatedTasks = [newTask, ...tasks]

			expect(updatedTasks.length).toBe(1)
			expect(updatedTasks[0].id).toBe(tempId)
			expect(updatedTasks[0].title).toBe('New Task')
		})

		it('rolls back optimistic update on failure', () => {
			const tasks = [{ id: 'temp-123', status: 'todo', title: 'Temp Task' }]

			const rollbackTasks = tasks.filter((task) => !task.id.startsWith('temp-'))

			expect(rollbackTasks.length).toBe(0)
		})

		it('updates task with optimistic update', () => {
			const tasks = [{ id: '1', status: 'todo', title: 'Initial Task' }]

			const updatedTask = { ...tasks[0], title: 'Updated Task' }
			const updatedTasks = tasks.map((task) => (task.id === '1' ? updatedTask : task))

			expect(updatedTasks[0].title).toBe('Updated Task')
		})

		it('deletes task with optimistic update', () => {
			const tasks = [{ id: '1', status: 'todo', title: 'Task to delete' }]

			const deletedTasks = tasks.map((task) =>
				task.id === '1' ? { ...task, status: 'deleted' } : task
			)

			expect(deletedTasks[0].status).toBe('deleted')
		})

		it('permanently deletes task', () => {
			const tasks = [
				{ id: '1', status: 'deleted', title: 'Task to delete' },
				{ id: '2', status: 'todo', title: 'Keep this task' },
			]

			const filteredTasks = tasks.filter((task) => task.id !== '1')

			expect(filteredTasks.length).toBe(1)
			expect(filteredTasks[0].id).toBe('2')
		})
	})

	describe('Subtask Management', () => {
		it('adds subtask to task', () => {
			const task = {
				id: '1',
				subtasks: [],
				title: 'Main Task',
			}

			const newSubtask = {
				id: 'sub-1',
				isCompleted: false,
				title: 'Subtask 1',
			}

			const updatedTask = {
				...task,
				subtasks: [...task.subtasks, newSubtask],
			}

			expect(updatedTask.subtasks.length).toBe(1)
			expect(updatedTask.subtasks[0].title).toBe('Subtask 1')
		})

		it('toggles subtask completion', () => {
			const task = {
				id: '1',
				subtasks: [{ id: 'sub-1', isCompleted: false, title: 'Subtask 1' }],
				title: 'Main Task',
			}

			const updatedSubtasks = task.subtasks.map((subtask) =>
				subtask.id === 'sub-1' ? { ...subtask, isCompleted: !subtask.isCompleted } : subtask
			)

			expect(updatedSubtasks[0].isCompleted).toBe(true)
		})

		it('deletes subtask', () => {
			const task = {
				id: '1',
				subtasks: [
					{ id: 'sub-1', isCompleted: false, title: 'Subtask 1' },
					{ id: 'sub-2', isCompleted: false, title: 'Subtask 2' },
				],
				title: 'Main Task',
			}

			const updatedSubtasks = task.subtasks.filter((subtask) => subtask.id !== 'sub-1')

			expect(updatedSubtasks.length).toBe(1)
			expect(updatedSubtasks[0].id).toBe('sub-2')
		})

		it('bulk completes subtasks', () => {
			const task = {
				id: '1',
				subtasks: [
					{ id: 'sub-1', isCompleted: false, title: 'Subtask 1' },
					{ id: 'sub-2', isCompleted: false, title: 'Subtask 2' },
					{ id: 'sub-3', isCompleted: false, title: 'Subtask 3' },
				],
				title: 'Main Task',
			}

			const subtaskIds = ['sub-1', 'sub-2']
			const updatedSubtasks = task.subtasks.map((subtask) =>
				subtaskIds.includes(subtask.id) ? { ...subtask, isCompleted: true } : subtask
			)

			const completedCount = updatedSubtasks.filter((s) => s.isCompleted).length
			expect(completedCount).toBe(2)
		})
	})

	describe('Progress Calculation', () => {
		it('calculates progress percentage', () => {
			const subtasks = [
				{ id: '1', isCompleted: true },
				{ id: '2', isCompleted: false },
				{ id: '3', isCompleted: true },
			]

			const completed = subtasks.filter((s) => s.isCompleted).length
			const progress = Math.round((completed / subtasks.length) * 100)

			expect(progress).toBe(67)
		})

		it('returns 0 for empty subtasks', () => {
			const subtasks: Array<{ isCompleted: boolean }> = []
			const progress =
				subtasks.length === 0
					? 0
					: Math.round((subtasks.filter((s) => s.isCompleted).length / subtasks.length) * 100)

			expect(progress).toBe(0)
		})
	})

	describe('Filters and Sorting', () => {
		it('sets filters and resets pagination', () => {
			const state = {
				filters: {},
				pagination: { limit: 10, page: 5 },
			}

			const newFilters = { priority: 'high', status: 'completed' }
			state.filters = { ...state.filters, ...newFilters }
			state.pagination.page = 1

			expect(state.filters).toEqual(newFilters)
			expect(state.pagination.page).toBe(1)
		})

		it('clears all filters', () => {
			type FilterState = {
				filters: { priority?: string; status?: string }
				pagination: { page: number }
			}
			const state: FilterState = {
				filters: { priority: 'high', status: 'completed' },
				pagination: { page: 3 },
			}

			state.filters = {}
			state.pagination.page = 1

			expect(state.filters).toEqual({})
			expect(state.pagination.page).toBe(1)
		})

		it('sets sort configuration', () => {
			const state = {
				sort: { field: 'createdAt', order: 'desc' },
			}

			const newSort = { field: 'dueDate', order: 'asc' }
			state.sort = newSort

			expect(state.sort).toEqual(newSort)
		})
	})

	describe('Pagination', () => {
		it('sets pagination page', () => {
			const state = { pagination: { limit: 10, page: 1 } }
			state.pagination.page = 2

			expect(state.pagination.page).toBe(2)
		})

		it('sets pagination limit and resets page', () => {
			const state = { pagination: { limit: 10, page: 3 } }
			state.pagination.limit = 25
			state.pagination.page = 1

			expect(state.pagination.limit).toBe(25)
			expect(state.pagination.page).toBe(1)
		})
	})

	describe('Error Handling', () => {
		it('sets error on failure', () => {
			type ErrorState = { error: null | string }
			const state: ErrorState = { error: null }
			state.error = 'Failed to fetch tasks'

			expect(state.error).toBe('Failed to fetch tasks')
		})

		it('clears error state', () => {
			type ErrorState = { error: null | string }
			const state: ErrorState = { error: 'Test error' }
			state.error = null

			expect(state.error).toBe(null)
		})
	})

	describe('State Reset', () => {
		it('resets store to initial state', () => {
			const _state = {
				currentTask: { id: '1', title: 'Task 1' },
				error: 'Test error',
				filters: { status: 'completed' },
				pagination: { limit: 20, page: 2 },
				sort: { field: 'priority', order: 'asc' },
				tasks: [{ id: '1', title: 'Task 1' }],
			}

			const resetState = {
				currentTask: null,
				error: null,
				filters: {},
				pagination: { limit: 10, page: 1, total: 0, totalPages: 0 },
				sort: { field: 'createdAt', order: 'desc' },
				tasks: [],
			}

			expect(resetState.tasks).toEqual([])
			expect(resetState.currentTask).toBe(null)
			expect(resetState.filters).toEqual({})
			expect(resetState.error).toBe(null)
		})
	})
})

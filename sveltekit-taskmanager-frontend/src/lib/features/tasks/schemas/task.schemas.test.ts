import { describe, it, expect } from 'vitest'
import {
	createTaskSchema,
	updateTaskSchema,
	type CreateTaskPayload,
	type UpdateTaskPayload,
} from './task.schemas'

describe('Task Schemas', () => {
	describe('createTaskSchema', () => {
		it('should validate valid task creation data', () => {
			const validData = {
				title: 'Test Task',
				description: 'Task description',
				priority: 'high',
				dueDate: '2024-12-31',
				tags: ['urgent', 'work'],
				categoryId: 'category-123',
			}
			const result = createTaskSchema.safeParse(validData)
			expect(result.success).toBe(true)
		})

		it('should validate task with only required fields', () => {
			const validData = {
				title: 'Test Task',
			}
			const result = createTaskSchema.safeParse(validData)
			expect(result.success).toBe(true)
		})

		it('should reject empty title', () => {
			const invalidData = {
				title: '',
			}
			const result = createTaskSchema.safeParse(invalidData)
			expect(result.success).toBe(false)
		})

		it('should reject title longer than 200 characters', () => {
			const invalidData = {
				title: 'A'.repeat(201),
			}
			const result = createTaskSchema.safeParse(invalidData)
			expect(result.success).toBe(false)
		})

		it('should reject description longer than 2000 characters', () => {
			const invalidData = {
				title: 'Test Task',
				description: 'A'.repeat(2001),
			}
			const result = createTaskSchema.safeParse(invalidData)
			expect(result.success).toBe(false)
		})

		it('should reject invalid date format', () => {
			const invalidData = {
				title: 'Test Task',
				dueDate: '31-12-2024',
			}
			const result = createTaskSchema.safeParse(invalidData)
			expect(result.success).toBe(false)
		})

		it('should reject invalid priority', () => {
			const invalidData = {
				title: 'Test Task',
				priority: 'invalid',
			}
			const result = createTaskSchema.safeParse(invalidData)
			expect(result.success).toBe(false)
		})

		it('should accept valid priority values', () => {
			const priorities = ['low', 'medium', 'high', 'urgent']
			priorities.forEach((priority) => {
				const data = {
					title: 'Test Task',
					priority,
				}
				const result = createTaskSchema.safeParse(data)
				expect(result.success).toBe(true)
			})
		})

		it('should trim title', () => {
			const data = {
				title: '  Test Task  ',
			}
			const result = createTaskSchema.safeParse(data)
			expect(result.success).toBe(true)
			if (result.success) {
				expect(result.data.title).toBe('Test Task')
			}
		})

		it('should trim description', () => {
			const data = {
				title: 'Test Task',
				description: '  Task description  ',
			}
			const result = createTaskSchema.safeParse(data)
			expect(result.success).toBe(true)
			if (result.success) {
				expect(result.data.description).toBe('Task description')
			}
		})

		it('should accept empty tags array', () => {
			const data = {
				title: 'Test Task',
				tags: [],
			}
			const result = createTaskSchema.safeParse(data)
			expect(result.success).toBe(true)
		})

		it('should accept subtasks', () => {
			const data = {
				title: 'Test Task',
				subtasks: [
					{
						id: 'subtask-1',
						taskId: 'task-1',
						title: 'Subtask 1',
						isCompleted: false,
						createdAt: '2024-01-01',
						updatedAt: '2024-01-01',
					},
				],
			}
			const result = createTaskSchema.safeParse(data)
			expect(result.success).toBe(true)
		})

		it('should infer correct type', () => {
			const data: CreateTaskPayload = {
				title: 'Test Task',
				description: 'Task description',
				priority: 'high',
			}
			expect(data.title).toBeDefined()
		})
	})

	describe('updateTaskSchema', () => {
		it('should validate valid task update data', () => {
			const validData = {
				title: 'Updated Task',
				description: 'Updated description',
				priority: 'medium',
				status: 'in_progress',
				dueDate: '2024-12-31',
				tags: ['updated'],
				categoryId: 'category-456',
			}
			const result = updateTaskSchema.safeParse(validData)
			expect(result.success).toBe(true)
		})

		it('should validate partial update with only title', () => {
			const validData = {
				title: 'Updated Task',
			}
			const result = updateTaskSchema.safeParse(validData)
			expect(result.success).toBe(true)
		})

		it('should validate empty update (all fields optional)', () => {
			const validData = {}
			const result = updateTaskSchema.safeParse(validData)
			expect(result.success).toBe(true)
		})

		it('should reject invalid status', () => {
			const invalidData = {
				status: 'invalid',
			}
			const result = updateTaskSchema.safeParse(invalidData)
			expect(result.success).toBe(false)
		})

		it('should accept valid status values', () => {
			const statuses = ['todo', 'in_progress', 'completed', 'cancelled', 'deleted']
			statuses.forEach((status) => {
				const data = {
					status,
				}
				const result = updateTaskSchema.safeParse(data)
				expect(result.success).toBe(true)
			})
		})

		it('should reject invalid priority', () => {
			const invalidData = {
				priority: 'invalid',
			}
			const result = updateTaskSchema.safeParse(invalidData)
			expect(result.success).toBe(false)
		})

		it('should reject invalid date format', () => {
			const invalidData = {
				dueDate: '31-12-2024',
			}
			const result = updateTaskSchema.safeParse(invalidData)
			expect(result.success).toBe(false)
		})

		it('should reject empty category ID when provided', () => {
			const invalidData = {
				categoryId: '',
			}
			const result = updateTaskSchema.safeParse(invalidData)
			expect(result.success).toBe(false)
		})

		it('should infer correct type', () => {
			const data: UpdateTaskPayload = {
				title: 'Updated Task',
				status: 'completed',
			}
			expect(data.title).toBeDefined()
			expect(data.status).toBeDefined()
		})
	})
})

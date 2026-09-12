/**
 * Task schemas with Zod validation for OWASP compliance
 */

import { z } from 'zod'

/**
 * Title validation schema
 */
const titleSchema = z
	.string()
	.min(1, 'Title is required')
	.max(200, 'Title is too long')
	.transform((val) => val.trim())

/**
 * Description validation schema
 */
const descriptionSchema = z
	.string()
	.max(2000, 'Description is too long')
	.transform((val) => val.trim())
	.optional()

/**
 * Date validation schema
 */
const dateSchema = z
	.string()
	.regex(/^\d{4}-\d{2}-\d{2}$/, 'Invalid date format. Use YYYY-MM-DD')
	.optional()

/**
 * ID validation schema
 */
const idSchema = z.string().min(1, 'ID is required')

/**
 * Task priority validation schema
 */
const prioritySchema = z.enum(['low', 'medium', 'high', 'urgent'], {
	message: 'Invalid priority. Must be low, medium, high, or urgent',
})

/**
 * Task status validation schema
 */
const statusSchema = z.enum(['todo', 'in_progress', 'completed', 'cancelled', 'deleted'], {
	message: 'Invalid status. Must be todo, in_progress, completed, cancelled, or deleted',
})

/**
 * Subtask validation schema
 */
const subtaskSchema = z.object({
	createdAt: z.string(),
	id: idSchema,
	isCompleted: z.boolean(),
	taskId: idSchema,
	title: titleSchema,
	updatedAt: z.string(),
})

/**
 * Create task validation schema
 */
export const createTaskSchema = z.object({
	categoryId: idSchema.optional(),
	description: descriptionSchema,
	dueDate: dateSchema,
	priority: prioritySchema.optional(),
	subtasks: z.array(subtaskSchema).optional(),
	tags: z.array(z.string()).optional(),
	title: titleSchema,
})

/**
 * Update task validation schema
 */
export const updateTaskSchema = z.object({
	categoryId: idSchema.optional(),
	description: descriptionSchema,
	dueDate: dateSchema.optional(),
	priority: prioritySchema.optional(),
	status: statusSchema.optional(),
	tags: z.array(z.string()).optional(),
	title: titleSchema.optional(),
})

/**
 * Type inference for create task payload
 */
export type CreateTaskPayload = z.infer<typeof createTaskSchema>

/**
 * Type inference for update task payload
 */
export type UpdateTaskPayload = z.infer<typeof updateTaskSchema>

/**
 * Task schemas with Zod validation for OWASP compliance
 */

import { z } from 'zod';

/**
 * Title validation schema
 */
const titleSchema = z
	.string()
	.min(1, 'Title is required')
	.max(200, 'Title is too long')
	.transform((val) => val.trim());

/**
 * Description validation schema
 */
const descriptionSchema = z
	.string()
	.max(2000, 'Description is too long')
	.transform((val) => val.trim())
	.optional();

/**
 * Date validation schema
 */
const dateSchema = z
	.string()
	.regex(/^\d{4}-\d{2}-\d{2}$/, 'Invalid date format. Use YYYY-MM-DD')
	.optional();

/**
 * ID validation schema
 */
const idSchema = z.string().min(1, 'ID is required');

/**
 * Task priority validation schema
 */
const prioritySchema = z.enum(['low', 'medium', 'high', 'urgent'], {
	message: 'Invalid priority. Must be low, medium, high, or urgent'
});

/**
 * Task status validation schema
 */
const statusSchema = z.enum(['todo', 'in_progress', 'completed', 'cancelled'], {
	message: 'Invalid status. Must be todo, in_progress, completed, or cancelled'
});

/**
 * Create task validation schema
 */
export const createTaskSchema = z.object({
	title: titleSchema,
	description: descriptionSchema,
	priority: prioritySchema.optional(),
	dueDate: dateSchema,
	categoryId: idSchema.optional(),
	tags: z.array(z.string()).optional()
});

/**
 * Update task validation schema
 */
export const updateTaskSchema = z.object({
	title: titleSchema.optional(),
	description: descriptionSchema,
	status: statusSchema.optional(),
	priority: prioritySchema.optional(),
	dueDate: dateSchema.optional(),
	categoryId: idSchema.optional(),
	tags: z.array(z.string()).optional()
});

/**
 * Type inference for create task payload
 */
export type CreateTaskPayload = z.infer<typeof createTaskSchema>;

/**
 * Type inference for update task payload
 */
export type UpdateTaskPayload = z.infer<typeof updateTaskSchema>;

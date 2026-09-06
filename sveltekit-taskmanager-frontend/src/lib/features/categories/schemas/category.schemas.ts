/**
 * Category schemas with Zod validation for OWASP compliance
 */

import { z } from 'zod';

/**
 * Name validation schema
 */
const nameSchema = z
	.string()
	.min(1, 'Name is required')
	.max(100, 'Name is too long')
	.transform((val) => val.trim());

/**
 * Description validation schema
 */
const descriptionSchema = z
	.string()
	.max(500, 'Description is too long')
	.transform((val) => val.trim())
	.optional();

/**
 * Color validation schema (hex color)
 */
const colorSchema = z
	.string()
	.regex(/^#[0-9A-Fa-f]{6}$/, 'Invalid color format. Use hex format like #FF0000')
	.optional();

/**
 * Create category validation schema
 */
export const createCategorySchema = z.object({
	name: nameSchema,
	description: descriptionSchema,
	color: colorSchema,
	icon: z.string().max(50, 'Icon is too long').optional()
});

/**
 * Update category validation schema
 */
export const updateCategorySchema = z.object({
	name: nameSchema.optional(),
	description: descriptionSchema,
	color: colorSchema,
	icon: z.string().max(50, 'Icon is too long').optional()
});

/**
 * Type inference for create category payload
 */
export type CreateCategoryPayload = z.infer<typeof createCategorySchema>;

/**
 * Type inference for update category payload
 */
export type UpdateCategoryPayload = z.infer<typeof updateCategorySchema>;

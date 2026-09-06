/**
 * Tag schemas with Zod validation for OWASP compliance
 */

import { z } from 'zod';

/**
 * Name validation schema
 */
const nameSchema = z
	.string()
	.min(1, 'Name is required')
	.max(50, 'Name is too long')
	.transform((val) => val.trim());

/**
 * Color validation schema (hex color)
 */
const colorSchema = z
	.string()
	.regex(/^#[0-9A-Fa-f]{6}$/, 'Invalid color format. Use hex format like #FF0000')
	.optional();

/**
 * Create tag validation schema
 */
export const createTagSchema = z.object({
	name: nameSchema,
	color: colorSchema
});

/**
 * Update tag validation schema
 */
export const updateTagSchema = z.object({
	name: nameSchema.optional(),
	color: colorSchema
});

/**
 * Type inference for create tag payload
 */
export type CreateTagPayload = z.infer<typeof createTagSchema>;

/**
 * Type inference for update tag payload
 */
export type UpdateTagPayload = z.infer<typeof updateTagSchema>;

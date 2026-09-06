/**
 * Category and tag related types and interfaces
 */

/**
 * Category interface
 */
export interface Category {
	id: string;
	name: string;
	description?: string;
	color?: string;
	icon?: string;
	userId: string;
	createdAt: string;
	updatedAt: string;
}

/**
 * Tag interface
 */
export interface Tag {
	id: string;
	name: string;
	color?: string;
	userId: string;
	createdAt: string;
	updatedAt: string;
}

/**
 * Category creation payload
 */
export interface CreateCategoryPayload {
	name: string;
	description?: string;
	color?: string;
	icon?: string;
}

/**
 * Category update payload
 */
export interface UpdateCategoryPayload {
	name?: string;
	description?: string;
	color?: string;
	icon?: string;
}

/**
 * Tag creation payload
 */
export interface CreateTagPayload {
	name: string;
	color?: string;
}

/**
 * Tag update payload
 */
export interface UpdateTagPayload {
	name?: string;
	color?: string;
}

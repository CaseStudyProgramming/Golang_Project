/**
 * Category feature types
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
 * Category store state interface
 */
export interface CategoryState {
	categories: Category[];
	currentCategory: Category | null;
	isLoading: boolean;
	error: string | null;
}

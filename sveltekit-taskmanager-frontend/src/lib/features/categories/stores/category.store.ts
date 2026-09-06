/**
 * Category management store using Svelte 5 runes
 */

import { httpClient } from '$lib/shared/utils/api.utils';
import { withErrorHandling, ValidationError } from '$lib/shared/utils/error.utils';
import { createCategorySchema, updateCategorySchema } from '../schemas/category.schemas';
import type { Category, CategoryState } from '../types/category.types';

/**
 * Create category store with Svelte 5 runes
 */
function createCategoryStore() {
	const state = $state<CategoryState>({
		categories: [],
		currentCategory: null,
		isLoading: false,
		error: null
	});

	/**
	 * Fetch all categories
	 */
	async function fetchCategories(): Promise<void> {
		state.isLoading = true;
		state.error = null;

		try {
			await withErrorHandling(async () => {
				// This would be replaced with actual API call
				// const response = await httpClient.get<Category[]>('/categories');
				
				// Mock response for development
				const mockCategories: Category[] = [];

				state.categories = mockCategories;
			}, 'Failed to fetch categories');
		} catch (error) {
			state.error = error instanceof Error ? error.message : 'Failed to fetch categories';
			throw error;
		} finally {
			state.isLoading = false;
		}
	}

	/**
	 * Fetch single category by ID
	 */
	async function fetchCategoryById(id: string): Promise<void> {
		state.isLoading = true;
		state.error = null;

		try {
			await withErrorHandling(async () => {
				// This would be replaced with actual API call
				// const response = await httpClient.get<Category>(`/categories/${id}`);
				
				// Mock response for development
				const mockCategory: Category = {
					id,
					name: 'Mock Category',
					createdAt: new Date().toISOString(),
					updatedAt: new Date().toISOString(),
					userId: '1'
				};

				state.currentCategory = mockCategory;
			}, 'Failed to fetch category');
		} catch (error) {
			state.error = error instanceof Error ? error.message : 'Failed to fetch category';
			throw error;
		} finally {
			state.isLoading = false;
		}
	}

	/**
	 * Create new category
	 */
	async function createCategory(payload: unknown): Promise<Category> {
		state.isLoading = true;
		state.error = null;

		try {
			// Validate input with Zod
			const validatedPayload = createCategorySchema.parse(payload);

			const newCategory = await withErrorHandling(async () => {
				// This would be replaced with actual API call
				// const response = await httpClient.post<Category>('/categories', validatedPayload);
				
				// Mock response for development
				const mockCategory: Category = {
					id: Date.now().toString(),
					name: validatedPayload.name,
					description: validatedPayload.description,
					color: validatedPayload.color,
					icon: validatedPayload.icon,
					createdAt: new Date().toISOString(),
					updatedAt: new Date().toISOString(),
					userId: '1'
				};

				state.categories = [...state.categories, mockCategory];

				return mockCategory;
			}, 'Failed to create category');

			return newCategory;
		} catch (error) {
			if (error instanceof Error && error.name === 'ZodError') {
				state.error = 'Invalid input: ' + error.message;
				throw new ValidationError('category', error.message);
			}
			state.error = error instanceof Error ? error.message : 'Failed to create category';
			throw error;
		} finally {
			state.isLoading = false;
		}
	}

	/**
	 * Update existing category
	 */
	async function updateCategory(id: string, payload: unknown): Promise<Category> {
		state.isLoading = true;
		state.error = null;

		try {
			// Validate input with Zod
			const validatedPayload = updateCategorySchema.parse(payload);

			const updatedCategory = await withErrorHandling(async () => {
				// This would be replaced with actual API call
				// const response = await httpClient.patch<Category>(`/categories/${id}`, validatedPayload);
				
				// Mock response for development
				const mockCategory: Category = {
					...state.categories.find((c) => c.id === id)!,
					...validatedPayload,
					updatedAt: new Date().toISOString()
				};

				state.categories = state.categories.map((category) =>
					category.id === id ? mockCategory : category
				);

				if (state.currentCategory?.id === id) {
					state.currentCategory = mockCategory;
				}

				return mockCategory;
			}, 'Failed to update category');

			return updatedCategory;
		} catch (error) {
			if (error instanceof Error && error.name === 'ZodError') {
				state.error = 'Invalid input: ' + error.message;
				throw new ValidationError('category', error.message);
			}
			state.error = error instanceof Error ? error.message : 'Failed to update category';
			throw error;
		} finally {
			state.isLoading = false;
		}
	}

	/**
	 * Delete category
	 */
	async function deleteCategory(id: string): Promise<void> {
		state.isLoading = true;
		state.error = null;

		try {
			await withErrorHandling(async () => {
				// This would be replaced with actual API call
				// await httpClient.delete(`/categories/${id}`);
				
				state.categories = state.categories.filter((category) => category.id !== id);

				if (state.currentCategory?.id === id) {
					state.currentCategory = null;
				}
			}, 'Failed to delete category');
		} catch (error) {
			state.error = error instanceof Error ? error.message : 'Failed to delete category';
			throw error;
		} finally {
			state.isLoading = false;
		}
	}

	/**
	 * Clear error state
	 */
	function clearError(): void {
		state.error = null;
	}

	/**
	 * Reset store state
	 */
	function reset(): void {
		state.categories = [];
		state.currentCategory = null;
		state.isLoading = false;
		state.error = null;
	}

	return {
		get state() {
			return state;
		},
		fetchCategories,
		fetchCategoryById,
		createCategory,
		updateCategory,
		deleteCategory,
		clearError,
		reset
	};
}

/**
 * Export category store instance
 */
export const categoryStore = createCategoryStore();

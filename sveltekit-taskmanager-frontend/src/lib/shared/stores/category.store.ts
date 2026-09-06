/**
 * Category management store using Svelte 5 runes
 */

import { httpClient } from '../utils/api.utils';
import type { Category, CreateCategoryPayload, UpdateCategoryPayload } from '../types/category.types';

/**
 * Category store state interface
 */
interface CategoryState {
	categories: Category[];
	currentCategory: Category | null;
	isLoading: boolean;
	error: string | null;
}

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
			// This would be replaced with actual API call
			// const response = await httpClient.get<Category[]>('/categories');
			
			// Mock response for development
			const mockCategories: Category[] = [];

			state.categories = mockCategories;
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
	async function createCategory(payload: CreateCategoryPayload): Promise<Category> {
		state.isLoading = true;
		state.error = null;

		try {
			// This would be replaced with actual API call
			// const response = await httpClient.post<Category>('/categories', payload);
			
			// Mock response for development
			const newCategory: Category = {
				id: Date.now().toString(),
				name: payload.name,
				description: payload.description,
				color: payload.color,
				icon: payload.icon,
				createdAt: new Date().toISOString(),
				updatedAt: new Date().toISOString(),
				userId: '1'
			};

			state.categories = [...state.categories, newCategory];

			return newCategory;
		} catch (error) {
			state.error = error instanceof Error ? error.message : 'Failed to create category';
			throw error;
		} finally {
			state.isLoading = false;
		}
	}

	/**
	 * Update existing category
	 */
	async function updateCategory(id: string, payload: UpdateCategoryPayload): Promise<Category> {
		state.isLoading = true;
		state.error = null;

		try {
			// This would be replaced with actual API call
			// const response = await httpClient.patch<Category>(`/categories/${id}`, payload);
			
			// Mock response for development
			const updatedCategory: Category = {
				...state.categories.find((c) => c.id === id)!,
				...payload,
				updatedAt: new Date().toISOString()
			};

			state.categories = state.categories.map((category) =>
				category.id === id ? updatedCategory : category
			);

			if (state.currentCategory?.id === id) {
				state.currentCategory = updatedCategory;
			}

			return updatedCategory;
		} catch (error) {
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
			// This would be replaced with actual API call
			// await httpClient.delete(`/categories/${id}`);
			
			state.categories = state.categories.filter((category) => category.id !== id);

			if (state.currentCategory?.id === id) {
				state.currentCategory = null;
			}
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

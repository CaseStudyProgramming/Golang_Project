/**
 * Tag management store using Svelte 5 runes
 */

import { httpClient } from '../utils/api.utils';
import type { Tag, CreateTagPayload, UpdateTagPayload } from '../types/category.types';

/**
 * Tag store state interface
 */
interface TagState {
	tags: Tag[];
	currentTag: Tag | null;
	isLoading: boolean;
	error: string | null;
}

/**
 * Create tag store with Svelte 5 runes
 */
function createTagStore() {
	const state = $state<TagState>({
		tags: [],
		currentTag: null,
		isLoading: false,
		error: null
	});

	/**
	 * Fetch all tags
	 */
	async function fetchTags(): Promise<void> {
		state.isLoading = true;
		state.error = null;

		try {
			// This would be replaced with actual API call
			// const response = await httpClient.get<Tag[]>('/tags');
			
			// Mock response for development
			const mockTags: Tag[] = [];

			state.tags = mockTags;
		} catch (error) {
			state.error = error instanceof Error ? error.message : 'Failed to fetch tags';
			throw error;
		} finally {
			state.isLoading = false;
		}
	}

	/**
	 * Fetch single tag by ID
	 */
	async function fetchTagById(id: string): Promise<void> {
		state.isLoading = true;
		state.error = null;

		try {
			// This would be replaced with actual API call
			// const response = await httpClient.get<Tag>(`/tags/${id}`);
			
			// Mock response for development
			const mockTag: Tag = {
				id,
				name: 'Mock Tag',
				createdAt: new Date().toISOString(),
				updatedAt: new Date().toISOString(),
				userId: '1'
			};

			state.currentTag = mockTag;
		} catch (error) {
			state.error = error instanceof Error ? error.message : 'Failed to fetch tag';
			throw error;
		} finally {
			state.isLoading = false;
		}
	}

	/**
	 * Create new tag
	 */
	async function createTag(payload: CreateTagPayload): Promise<Tag> {
		state.isLoading = true;
		state.error = null;

		try {
			// This would be replaced with actual API call
			// const response = await httpClient.post<Tag>('/tags', payload);
			
			// Mock response for development
			const newTag: Tag = {
				id: Date.now().toString(),
				name: payload.name,
				color: payload.color,
				createdAt: new Date().toISOString(),
				updatedAt: new Date().toISOString(),
				userId: '1'
			};

			state.tags = [...state.tags, newTag];

			return newTag;
		} catch (error) {
			state.error = error instanceof Error ? error.message : 'Failed to create tag';
			throw error;
		} finally {
			state.isLoading = false;
		}
	}

	/**
	 * Update existing tag
	 */
	async function updateTag(id: string, payload: UpdateTagPayload): Promise<Tag> {
		state.isLoading = true;
		state.error = null;

		try {
			// This would be replaced with actual API call
			// const response = await httpClient.patch<Tag>(`/tags/${id}`, payload);
			
			// Mock response for development
			const updatedTag: Tag = {
				...state.tags.find((t) => t.id === id)!,
				...payload,
				updatedAt: new Date().toISOString()
			};

			state.tags = state.tags.map((tag) => (tag.id === id ? updatedTag : tag));

			if (state.currentTag?.id === id) {
				state.currentTag = updatedTag;
			}

			return updatedTag;
		} catch (error) {
			state.error = error instanceof Error ? error.message : 'Failed to update tag';
			throw error;
		} finally {
			state.isLoading = false;
		}
	}

	/**
	 * Delete tag
	 */
	async function deleteTag(id: string): Promise<void> {
		state.isLoading = true;
		state.error = null;

		try {
			// This would be replaced with actual API call
			// await httpClient.delete(`/tags/${id}`);
			
			state.tags = state.tags.filter((tag) => tag.id !== id);

			if (state.currentTag?.id === id) {
				state.currentTag = null;
			}
		} catch (error) {
			state.error = error instanceof Error ? error.message : 'Failed to delete tag';
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
		state.tags = [];
		state.currentTag = null;
		state.isLoading = false;
		state.error = null;
	}

	return {
		get state() {
			return state;
		},
		fetchTags,
		fetchTagById,
		createTag,
		updateTag,
		deleteTag,
		clearError,
		reset
	};
}

/**
 * Export tag store instance
 */
export const tagStore = createTagStore();

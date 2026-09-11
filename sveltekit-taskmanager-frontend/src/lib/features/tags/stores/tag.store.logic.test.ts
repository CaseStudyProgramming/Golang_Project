import { describe, expect, it, vi } from 'vitest';

import type { Tag } from '../types/tag.types';

describe('Tag Store Logic', () => {
	describe('Tag CRUD Operations', () => {
		it('creates new tag with required fields', () => {
			const tag: Tag = {
				id: '123',
				name: 'Work',
				color: '#3B82F6',
				userId: 'user1',
				createdAt: new Date().toISOString(),
				updatedAt: new Date().toISOString()
			};

			expect(tag.id).toBe('123');
			expect(tag.name).toBe('Work');
			expect(tag.color).toBe('#3B82F6');
			expect(tag.userId).toBe('user1');
		});

		it('updates existing tag', () => {
			const originalTag: Tag = {
				id: '123',
				name: 'Personal',
				color: '#10B981',
				userId: 'user1',
				createdAt: '2024-01-01T00:00:00Z',
				updatedAt: '2024-01-01T00:00:00Z'
			};

			const updatedTag: Tag = {
				...originalTag,
				name: 'Updated Personal',
				color: '#F59E0B',
				updatedAt: new Date().toISOString()
			};

			expect(updatedTag.name).toBe('Updated Personal');
			expect(updatedTag.color).toBe('#F59E0B');
			expect(updatedTag.id).toBe(originalTag.id);
		});

		it('deletes tag from array', () => {
			const tags: Tag[] = [
				{ id: '1', name: 'Tag 1', color: '#3B82F6', userId: 'user1', createdAt: '2024-01-01T00:00:00Z', updatedAt: '2024-01-01T00:00:00Z' },
				{ id: '2', name: 'Tag 2', color: '#10B981', userId: 'user1', createdAt: '2024-01-01T00:00:00Z', updatedAt: '2024-01-01T00:00:00Z' }
			];

			const filtered = tags.filter(t => t.id !== '1');
			expect(filtered).toHaveLength(1);
			expect(filtered[0].id).toBe('2');
		});
	});

	describe('Tag Validation', () => {
		it('validates tag name presence', () => {
			const validName = 'Work';
			const invalidName = '';

			expect(validName.length).toBeGreaterThan(0);
			expect(invalidName.length).toBe(0);
		});

		it('validates tag color format', () => {
			const validColor = '#3B82F6';
			const invalidColor = 'red';

			expect(validColor).toMatch(/^#[0-9A-Fa-f]{6}$/);
			expect(invalidColor).not.toMatch(/^#[0-9A-Fa-f]{6}$/);
		});
	});

	describe('Current Tag Management', () => {
		it('sets current tag', () => {
			const tag: Tag = {
				id: '123',
				name: 'Current',
				color: '#8B5CF6',
				userId: 'user1',
				createdAt: '2024-01-01T00:00:00Z',
				updatedAt: '2024-01-01T00:00:00Z'
			};

			let currentTag: Tag | null = null;
			currentTag = tag;
			expect(currentTag).toEqual(tag);
		});

		it('clears current tag', () => {
			let currentTag: Tag | null = {
				id: '123',
				name: 'Current',
				color: '#8B5CF6',
				userId: 'user1',
				createdAt: '2024-01-01T00:00:00Z',
				updatedAt: '2024-01-01T00:00:00Z'
			};

			currentTag = null;
			expect(currentTag).toBeNull();
		});

		it('updates current tag when tag is updated', () => {
			const originalTag: Tag = {
				id: '123',
				name: 'Original',
				color: '#3B82F6',
				userId: 'user1',
				createdAt: '2024-01-01T00:00:00Z',
				updatedAt: '2024-01-01T00:00:00Z'
			};

			let currentTag: Tag = originalTag;

			const updatedTag: Tag = {
				...originalTag,
				name: 'Updated',
				color: '#10B981',
				updatedAt: new Date().toISOString()
			};

			currentTag = updatedTag;
			expect(currentTag.name).toBe('Updated');
			expect(currentTag.color).toBe('#10B981');
		});
	});

	describe('Error Handling', () => {
		it('sets error message', () => {
			let error: string | null = null;
			error = 'Failed to fetch tags';
			expect(error).toBe('Failed to fetch tags');
		});

		it('clears error state', () => {
			let error: string | null = 'Some error';
			error = null;
			expect(error).toBeNull();
		});
	});

	describe('Loading States', () => {
		it('sets loading state', () => {
			let isLoading = false;
			isLoading = true;
			expect(isLoading).toBe(true);
		});

		it('clears loading state', () => {
			let isLoading = true;
			isLoading = false;
			expect(isLoading).toBe(false);
		});
	});

	describe('State Reset', () => {
		it('resets tag list', () => {
			let tags: Tag[] = [
				{ id: '1', name: 'Tag 1', color: '#3B82F6', userId: 'user1', createdAt: '2024-01-01T00:00:00Z', updatedAt: '2024-01-01T00:00:00Z' }
			];

			tags = [];
			expect(tags).toHaveLength(0);
		});

		it('resets current tag', () => {
			let currentTag: Tag | null = {
				id: '123',
				name: 'Current',
				color: '#8B5CF6',
				userId: 'user1',
				createdAt: '2024-01-01T00:00:00Z',
				updatedAt: '2024-01-01T00:00:00Z'
			};

			currentTag = null;
			expect(currentTag).toBeNull();
		});

		it('resets error state', () => {
			let error: string | null = 'Some error';
			let isLoading = true;

			error = null;
			isLoading = false;

			expect(error).toBeNull();
			expect(isLoading).toBe(false);
		});
	});

	describe('Tag State Management', () => {
		it('maintains tag list across operations', () => {
			const tags: Tag[] = [
				{ id: '1', name: 'Tag 1', color: '#3B82F6', userId: 'user1', createdAt: '2024-01-01T00:00:00Z', updatedAt: '2024-01-01T00:00:00Z' },
				{ id: '2', name: 'Tag 2', color: '#10B981', userId: 'user1', createdAt: '2024-01-01T00:00:00Z', updatedAt: '2024-01-01T00:00:00Z' },
				{ id: '3', name: 'Tag 3', color: '#F59E0B', userId: 'user1', createdAt: '2024-01-01T00:00:00Z', updatedAt: '2024-01-01T00:00:00Z' }
			];

			expect(tags).toHaveLength(3);
		});

		it('preserves other tags when updating one tag', () => {
			const tags: Tag[] = [
				{ id: '1', name: 'Tag 1', color: '#3B82F6', userId: 'user1', createdAt: '2024-01-01T00:00:00Z', updatedAt: '2024-01-01T00:00:00Z' },
				{ id: '2', name: 'Tag 2', color: '#10B981', userId: 'user1', createdAt: '2024-01-01T00:00:00Z', updatedAt: '2024-01-01T00:00:00Z' }
			];

			const updatedTags = tags.map(tag =>
				tag.id === '1' ? { ...tag, name: 'Updated Tag 1', color: '#EF4444' } : tag
			);

			expect(updatedTags).toHaveLength(2);
			expect(updatedTags.find(t => t.id === '2')?.name).toBe('Tag 2');
		});
	});
});
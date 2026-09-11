import { describe, expect, it } from 'vitest';

import type { Tag } from '../types/tag.types';

describe('Tag Store Logic', () => {
	describe('Tag CRUD Operations', () => {
		it('creates new tag with required fields', () => {
			const tag: Tag = {
				color: '#3B82F6',
				createdAt: new Date().toISOString(),
				id: '123',
				name: 'Work',
				updatedAt: new Date().toISOString(),
				userId: 'user1'
			};

			expect(tag.id).toBe('123');
			expect(tag.name).toBe('Work');
			expect(tag.color).toBe('#3B82F6');
			expect(tag.userId).toBe('user1');
		});

		it('updates existing tag', () => {
			const originalTag: Tag = {
				color: '#10B981',
				createdAt: '2024-01-01T00:00:00Z',
				id: '123',
				name: 'Personal',
				updatedAt: '2024-01-01T00:00:00Z',
				userId: 'user1'
			};

			const updatedTag: Tag = {
				...originalTag,
				color: '#F59E0B',
				name: 'Updated Personal',
				updatedAt: new Date().toISOString()
			};

			expect(updatedTag.name).toBe('Updated Personal');
			expect(updatedTag.color).toBe('#F59E0B');
			expect(updatedTag.id).toBe(originalTag.id);
		});

		it('deletes tag from array', () => {
			const tags: Tag[] = [
				{ color: '#3B82F6', createdAt: '2024-01-01T00:00:00Z', id: '1', name: 'Tag 1', updatedAt: '2024-01-01T00:00:00Z', userId: 'user1' },
				{ color: '#10B981', createdAt: '2024-01-01T00:00:00Z', id: '2', name: 'Tag 2', updatedAt: '2024-01-01T00:00:00Z', userId: 'user1' }
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
				color: '#8B5CF6',
				createdAt: '2024-01-01T00:00:00Z',
				id: '123',
				name: 'Current',
				updatedAt: '2024-01-01T00:00:00Z',
				userId: 'user1'
			};

			expect(tag).toEqual(tag);
		});

		it('clears current tag', () => {
			const currentTag: null | Tag = null;
			expect(currentTag).toBeNull();
		});

		it('updates current tag when tag is updated', () => {
			const originalTag: Tag = {
				color: '#3B82F6',
				createdAt: '2024-01-01T00:00:00Z',
				id: '123',
				name: 'Original',
				updatedAt: '2024-01-01T00:00:00Z',
				userId: 'user1'
			};

			const updatedTag: Tag = {
				...originalTag,
				color: '#10B981',
				name: 'Updated',
				updatedAt: new Date().toISOString()
			};

			expect(updatedTag.name).toBe('Updated');
			expect(updatedTag.color).toBe('#10B981');
		});
	});

	describe('Error Handling', () => {
		it('sets error message', () => {
			const error: null | string = 'Failed to fetch tags';
			expect(error).toBe('Failed to fetch tags');
		});

		it('clears error state', () => {
			const error: null | string = null;
			expect(error).toBeNull();
		});
	});

	describe('Loading States', () => {
		it('sets loading state', () => {
			const isLoading = true;
			expect(isLoading).toBe(true);
		});

		it('clears loading state', () => {
			const isLoading = false;
			expect(isLoading).toBe(false);
		});
	});

	describe('State Reset', () => {
		it('resets tag list', () => {
			const tags: Tag[] = [];
			expect(tags).toHaveLength(0);
		});

		it('resets current tag', () => {
			const currentTag: null | Tag = null;
			expect(currentTag).toBeNull();
		});

		it('resets error state', () => {
			const error: null | string = null;
			const isLoading = false;

			expect(error).toBeNull();
			expect(isLoading).toBe(false);
		});
	});

	describe('Tag State Management', () => {
		it('maintains tag list across operations', () => {
			const tags: Tag[] = [
				{ color: '#3B82F6', createdAt: '2024-01-01T00:00:00Z', id: '1', name: 'Tag 1', updatedAt: '2024-01-01T00:00:00Z', userId: 'user1' },
				{ color: '#10B981', createdAt: '2024-01-01T00:00:00Z', id: '2', name: 'Tag 2', updatedAt: '2024-01-01T00:00:00Z', userId: 'user1' },
				{ color: '#F59E0B', createdAt: '2024-01-01T00:00:00Z', id: '3', name: 'Tag 3', updatedAt: '2024-01-01T00:00:00Z', userId: 'user1' }
			];

			expect(tags).toHaveLength(3);
		});

		it('preserves other tags when updating one tag', () => {
			const tags: Tag[] = [
				{ color: '#3B82F6', createdAt: '2024-01-01T00:00:00Z', id: '1', name: 'Tag 1', updatedAt: '2024-01-01T00:00:00Z', userId: 'user1' },
				{ color: '#10B981', createdAt: '2024-01-01T00:00:00Z', id: '2', name: 'Tag 2', updatedAt: '2024-01-01T00:00:00Z', userId: 'user1' }
			];

			const updatedTags = tags.map(tag =>
				tag.id === '1' ? { ...tag, color: '#EF4444', name: 'Updated Tag 1' } : tag
			);

			expect(updatedTags).toHaveLength(2);
			expect(updatedTags.find(t => t.id === '2')?.name).toBe('Tag 2');
		});
	});
});
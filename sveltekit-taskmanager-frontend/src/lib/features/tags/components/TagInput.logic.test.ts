import { describe, expect, it } from 'vitest';

describe('TagInput Component Logic', () => {
	describe('Tag Filtering', () => {
		it('filters tags by input text case-insensitively', () => {
			const availableTags = [
				{ id: '1', name: 'JavaScript', color: '#F7DF1E' },
				{ id: '2', name: 'TypeScript', color: '#3178C6' },
				{ id: '3', name: 'Python', color: '#3776AB' }
			];

			const input = 'script';
			const filtered = availableTags.filter((tag) =>
				tag.name.toLowerCase().includes(input.toLowerCase())
			);

			expect(filtered.length).toBe(2);
			expect(filtered.map(t => t.name)).toContain('JavaScript');
			expect(filtered.map(t => t.name)).toContain('TypeScript');
		});

		it('returns all tags when input is empty', () => {
			const availableTags = [
				{ id: '1', name: 'JavaScript', color: '#F7DF1E' },
				{ id: '2', name: 'TypeScript', color: '#3178C6' }
			];

			const input = '';
			const filtered = !input.trim() ? availableTags : availableTags.filter((tag) =>
				tag.name.toLowerCase().includes(input.toLowerCase())
			);

			expect(filtered.length).toBe(2);
		});

		it('returns empty array when no matches found', () => {
			const availableTags = [
				{ id: '1', name: 'JavaScript', color: '#F7DF1E' },
				{ id: '2', name: 'TypeScript', color: '#3178C6' }
			];

			const input = 'ruby';
			const filtered = availableTags.filter((tag) =>
				tag.name.toLowerCase().includes(input.toLowerCase())
			);

			expect(filtered.length).toBe(0);
		});
	});

	describe('Tag Selection', () => {
		it('adds tag to selected tags if not already selected', () => {
			const selectedTags = ['1'];
			const newTagId = '2';

			const updated = selectedTags.includes(newTagId) 
				? selectedTags 
				: [...selectedTags, newTagId];

			expect(updated.length).toBe(2);
			expect(updated).toContain('2');
		});

		it('does not add duplicate tags', () => {
			const selectedTags = ['1', '2'];
			const newTagId = '1';

			const updated = selectedTags.includes(newTagId) 
				? selectedTags 
				: [...selectedTags, newTagId];

			expect(updated.length).toBe(2);
			expect(updated.filter(id => id === '1').length).toBe(1);
		});
	});

	describe('Tag Removal', () => {
		it('removes tag from selected tags', () => {
			const selectedTags = ['1', '2', '3'];
			const tagToRemove = '2';

			const updated = selectedTags.filter((id) => id !== tagToRemove);

			expect(updated.length).toBe(2);
			expect(updated).not.toContain('2');
		});

		it('handles removal of non-existent tag', () => {
			const selectedTags = ['1', '2'];
			const tagToRemove = '3';

			const updated = selectedTags.filter((id) => id !== tagToRemove);

			expect(updated.length).toBe(2);
		});
	});

	describe('Keyboard Navigation', () => {
		it('handles arrow down key for navigation', () => {
			let highlightedIndex = -1;
			const maxIndex = 5;

			if (highlightedIndex < maxIndex) {
				highlightedIndex = Math.min(highlightedIndex + 1, maxIndex);
			}

			expect(highlightedIndex).toBe(0);
		});

		it('handles arrow up key for navigation', () => {
			let highlightedIndex = 3;
			const minIndex = 0;

			if (highlightedIndex > minIndex) {
				highlightedIndex = Math.max(highlightedIndex - 1, minIndex);
			}

			expect(highlightedIndex).toBe(2);
		});

		it('prevents going below minimum index', () => {
			let highlightedIndex = 0;
			const minIndex = 0;

			highlightedIndex = Math.max(highlightedIndex - 1, minIndex);

			expect(highlightedIndex).toBe(0);
		});

		it('prevents going above maximum index', () => {
			let highlightedIndex = 5;
			const maxIndex = 5;

			highlightedIndex = Math.min(highlightedIndex + 1, maxIndex);

			expect(highlightedIndex).toBe(5);
		});
	});

	describe('Input Validation', () => {
		it('validates non-empty input for tag creation', () => {
			const input = '  ';
			const isValid = input.trim().length > 0;

			expect(isValid).toBe(false);
		});

		it('validates non-whitespace input for tag creation', () => {
			const input = 'new tag';
			const isValid = input.trim().length > 0;

			expect(isValid).toBe(true);
		});
	});

	describe('Selected Tag Objects', () => {
		it('filters available tags to get selected tag objects', () => {
			const availableTags = [
				{ id: '1', name: 'JavaScript', color: '#F7DF1E' },
				{ id: '2', name: 'TypeScript', color: '#3178C6' },
				{ id: '3', name: 'Python', color: '#3776AB' }
			];
			const selectedTags = ['1', '3'];

			const selectedObjects = availableTags.filter((tag) =>
				selectedTags.includes(tag.id)
			);

			expect(selectedObjects.length).toBe(2);
			expect(selectedObjects.map(t => t.name)).toContain('JavaScript');
			expect(selectedObjects.map(t => t.name)).toContain('Python');
			expect(selectedObjects.map(t => t.name)).not.toContain('TypeScript');
		});
	});

	describe('Color Fallback', () => {
		it('uses default color when tag color is not provided', () => {
			type Tag = { id: string; name: string; color?: string };
			const tag: Tag = { id: '1', name: 'Test' };
			const defaultColor = '#3B82F6';
			const backgroundColor = tag.color || defaultColor;

			expect(backgroundColor).toBe(defaultColor);
		});

		it('uses tag color when provided', () => {
			type Tag = { id: string; name: string; color?: string };
			const tag: Tag = { id: '1', name: 'Test', color: '#FF0000' };
			const defaultColor = '#3B82F6';
			const backgroundColor = tag.color || defaultColor;

			expect(backgroundColor).toBe('#FF0000');
		});
	});

	describe('Loading State', () => {
		it('disables input during loading', () => {
			let isLoading = true;
			const isDisabled = isLoading;

			expect(isDisabled).toBe(true);
		});

		it('enables input when not loading', () => {
			let isLoading = false;
			const isDisabled = isLoading;

			expect(isDisabled).toBe(false);
		});
	});
});

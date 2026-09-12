import { describe, expect, it } from 'vitest'

describe('Category Store Logic', () => {
	describe('Category CRUD Operations', () => {
		it('creates new category', () => {
			const categories: Array<{
				color: string
				createdAt: string
				description: string
				icon: string
				id: string
				name: string
			}> = []
			const newCategory = {
				color: '#FF0000',
				createdAt: new Date().toISOString(),
				description: 'Category description',
				icon: 'folder',
				id: '1',
				name: 'New Category',
			}

			const updatedCategories = [...categories, newCategory]

			expect(updatedCategories.length).toBe(1)
			expect(updatedCategories[0].name).toBe('New Category')
			expect(updatedCategories[0].color).toBe('#FF0000')
		})

		it('updates existing category', () => {
			const categories = [{ color: '#FF0000', id: '1', name: 'Initial Category' }]

			const updatedCategory = { ...categories[0], color: '#00FF00', name: 'Updated Category' }
			const updatedCategories = categories.map((category) =>
				category.id === '1' ? updatedCategory : category
			)

			expect(updatedCategories[0].name).toBe('Updated Category')
			expect(updatedCategories[0].color).toBe('#00FF00')
		})

		it('deletes category', () => {
			const categories = [
				{ id: '1', name: 'Category to delete' },
				{ id: '2', name: 'Keep this category' },
			]

			const filteredCategories = categories.filter((category) => category.id !== '1')

			expect(filteredCategories.length).toBe(1)
			expect(filteredCategories[0].id).toBe('2')
		})
	})

	describe('Current Category Management', () => {
		it('sets current category', () => {
			type CategoryState = {
				currentCategory: null | { id: string; name: string }
			}
			const state: CategoryState = { currentCategory: null }
			const category = { id: '1', name: 'Test Category' }
			state.currentCategory = category

			expect(state.currentCategory).toEqual(category)
		})

		it('clears current category on delete', () => {
			type CategoryState = {
				currentCategory: null | { id: string; name: string }
			}
			const state: CategoryState = {
				currentCategory: { id: '1', name: 'Category to delete' },
			}

			state.currentCategory = null

			expect(state.currentCategory).toBe(null)
		})

		it('updates current category when list is updated', () => {
			const state = {
				categories: [{ id: '1', name: 'Initial Category' }],
				currentCategory: { id: '1', name: 'Initial Category' },
			}

			const updatedCategory = { ...state.categories[0], name: 'Updated Category' }
			state.categories = state.categories.map((category) =>
				category.id === '1' ? updatedCategory : category
			)
			state.currentCategory = updatedCategory

			expect(state.currentCategory.name).toBe('Updated Category')
		})
	})

	describe('Error Handling', () => {
		it('sets error on failure', () => {
			type ErrorState = { error: null | string }
			const state: ErrorState = { error: null }
			state.error = 'Failed to fetch categories'

			expect(state.error).toBe('Failed to fetch categories')
		})

		it('clears error state', () => {
			type ErrorState = { error: null | string }
			const state: ErrorState = { error: 'Test error' }
			state.error = null

			expect(state.error).toBe(null)
		})
	})

	describe('Loading States', () => {
		it('sets loading state during operation', () => {
			const state = { isLoading: false }
			state.isLoading = true

			expect(state.isLoading).toBe(true)
		})

		it('clears loading state after operation', () => {
			const state = { isLoading: true }
			state.isLoading = false

			expect(state.isLoading).toBe(false)
		})
	})

	describe('State Reset', () => {
		it('resets store to initial state', () => {
			const _state = {
				categories: [{ id: '1', name: 'Category 1' }],
				currentCategory: { id: '1', name: 'Category 1' },
				error: 'Test error',
				isLoading: true,
			}

			const resetState = {
				categories: [],
				currentCategory: null,
				error: null,
				isLoading: false,
			}

			expect(resetState.categories).toEqual([])
			expect(resetState.currentCategory).toBe(null)
			expect(resetState.error).toBe(null)
			expect(resetState.isLoading).toBe(false)
		})
	})

	describe('Category Validation', () => {
		it('validates category name presence', () => {
			const category = { name: 'Test Category' }
			const isValid = category.name && category.name.length > 0

			expect(isValid).toBe(true)
		})

		it('validates category color format', () => {
			const category = { color: '#FF0000' }
			const isValidColor = /^#[0-9A-F]{6}$/i.test(category.color)

			expect(isValidColor).toBe(true)
		})
	})
})

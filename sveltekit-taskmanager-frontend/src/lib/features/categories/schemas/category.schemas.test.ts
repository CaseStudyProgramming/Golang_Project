import { describe, it, expect } from 'vitest'
import {
	createCategorySchema,
	updateCategorySchema,
	type CreateCategoryPayload,
	type UpdateCategoryPayload,
} from './category.schemas'

describe('Category Schemas', () => {
	describe('createCategorySchema', () => {
		it('should validate valid category creation data', () => {
			const validData = {
				name: 'Work',
				description: 'Work related tasks',
				color: '#FF0000',
				icon: 'briefcase',
			}
			const result = createCategorySchema.safeParse(validData)
			expect(result.success).toBe(true)
		})

		it('should validate category with only required fields', () => {
			const validData = {
				name: 'Work',
			}
			const result = createCategorySchema.safeParse(validData)
			expect(result.success).toBe(true)
		})

		it('should reject empty name', () => {
			const invalidData = {
				name: '',
			}
			const result = createCategorySchema.safeParse(invalidData)
			expect(result.success).toBe(false)
		})

		it('should reject name longer than 100 characters', () => {
			const invalidData = {
				name: 'A'.repeat(101),
			}
			const result = createCategorySchema.safeParse(invalidData)
			expect(result.success).toBe(false)
		})

		it('should reject description longer than 500 characters', () => {
			const invalidData = {
				name: 'Work',
				description: 'A'.repeat(501),
			}
			const result = createCategorySchema.safeParse(invalidData)
			expect(result.success).toBe(false)
		})

		it('should reject icon longer than 50 characters', () => {
			const invalidData = {
				name: 'Work',
				icon: 'A'.repeat(51),
			}
			const result = createCategorySchema.safeParse(invalidData)
			expect(result.success).toBe(false)
		})

		it('should reject invalid color format', () => {
			const invalidData = {
				name: 'Work',
				color: 'red',
			}
			const result = createCategorySchema.safeParse(invalidData)
			expect(result.success).toBe(false)
		})

		it('should accept valid hex color formats', () => {
			const validColors = ['#FF0000', '#00FF00', '#0000FF', '#ffffff', '#123ABC']
			validColors.forEach((color) => {
				const data = {
					name: 'Work',
					color,
				}
				const result = createCategorySchema.safeParse(data)
				expect(result.success).toBe(true)
			})
		})

		it('should trim name', () => {
			const data = {
				name: '  Work  ',
			}
			const result = createCategorySchema.safeParse(data)
			expect(result.success).toBe(true)
			if (result.success) {
				expect(result.data.name).toBe('Work')
			}
		})

		it('should trim description', () => {
			const data = {
				name: 'Work',
				description: '  Work description  ',
			}
			const result = createCategorySchema.safeParse(data)
			expect(result.success).toBe(true)
			if (result.success) {
				expect(result.data.description).toBe('Work description')
			}
		})

		it('should infer correct type', () => {
			const data: CreateCategoryPayload = {
				name: 'Work',
				description: 'Work related tasks',
			}
			expect(data.name).toBeDefined()
		})
	})

	describe('updateCategorySchema', () => {
		it('should validate valid category update data', () => {
			const validData = {
				name: 'Personal',
				description: 'Personal tasks',
				color: '#00FF00',
				icon: 'home',
			}
			const result = updateCategorySchema.safeParse(validData)
			expect(result.success).toBe(true)
		})

		it('should validate partial update with only description', () => {
			const validData = {
				description: 'Updated description',
			}
			const result = updateCategorySchema.safeParse(validData)
			expect(result.success).toBe(true)
		})

		it('should validate empty update (all fields optional)', () => {
			const validData = {}
			const result = updateCategorySchema.safeParse(validData)
			expect(result.success).toBe(true)
		})

		it('should reject invalid color format', () => {
			const invalidData = {
				color: 'red',
			}
			const result = updateCategorySchema.safeParse(invalidData)
			expect(result.success).toBe(false)
		})

		it('should reject empty name when provided', () => {
			const invalidData = {
				name: '',
			}
			const result = updateCategorySchema.safeParse(invalidData)
			expect(result.success).toBe(false)
		})

		it('should reject name longer than 100 characters', () => {
			const invalidData = {
				name: 'A'.repeat(101),
			}
			const result = updateCategorySchema.safeParse(invalidData)
			expect(result.success).toBe(false)
		})

		it('should reject description longer than 500 characters', () => {
			const invalidData = {
				description: 'A'.repeat(501),
			}
			const result = updateCategorySchema.safeParse(invalidData)
			expect(result.success).toBe(false)
		})

		it('should reject icon longer than 50 characters', () => {
			const invalidData = {
				icon: 'A'.repeat(51),
			}
			const result = updateCategorySchema.safeParse(invalidData)
			expect(result.success).toBe(false)
		})

		it('should infer correct type', () => {
			const data: UpdateCategoryPayload = {
				name: 'Updated',
				color: '#00FF00',
			}
			expect(data.name).toBeDefined()
			expect(data.color).toBeDefined()
		})
	})
})

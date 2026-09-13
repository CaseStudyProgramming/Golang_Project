import { describe, it, expect } from 'vitest'
import {
	createTagSchema,
	updateTagSchema,
	type CreateTagPayload,
	type UpdateTagPayload,
} from './tag.schemas'

describe('Tag Schemas', () => {
	describe('createTagSchema', () => {
		it('should validate valid tag creation data', () => {
			const validData = {
				name: 'Work',
				color: '#FF0000',
			}
			const result = createTagSchema.safeParse(validData)
			expect(result.success).toBe(true)
		})

		it('should validate tag without color', () => {
			const validData = {
				name: 'Work',
			}
			const result = createTagSchema.safeParse(validData)
			expect(result.success).toBe(true)
		})

		it('should reject empty name', () => {
			const invalidData = {
				name: '',
			}
			const result = createTagSchema.safeParse(invalidData)
			expect(result.success).toBe(false)
		})

		it('should reject name longer than 50 characters', () => {
			const invalidData = {
				name: 'A'.repeat(51),
			}
			const result = createTagSchema.safeParse(invalidData)
			expect(result.success).toBe(false)
		})

		it('should reject invalid color format', () => {
			const invalidData = {
				name: 'Work',
				color: 'red',
			}
			const result = createTagSchema.safeParse(invalidData)
			expect(result.success).toBe(false)
		})

		it('should reject invalid hex color format', () => {
			const invalidData = {
				name: 'Work',
				color: '#FF00',
			}
			const result = createTagSchema.safeParse(invalidData)
			expect(result.success).toBe(false)
		})

		it('should accept valid hex color formats', () => {
			const validColors = ['#FF0000', '#00FF00', '#0000FF', '#ffffff', '#123ABC']
			validColors.forEach((color) => {
				const data = {
					name: 'Work',
					color,
				}
				const result = createTagSchema.safeParse(data)
				expect(result.success).toBe(true)
			})
		})

		it('should trim name', () => {
			const data = {
				name: '  Work  ',
			}
			const result = createTagSchema.safeParse(data)
			expect(result.success).toBe(true)
			if (result.success) {
				expect(result.data.name).toBe('Work')
			}
		})

		it('should infer correct type', () => {
			const data: CreateTagPayload = {
				name: 'Work',
				color: '#FF0000',
			}
			expect(data.name).toBeDefined()
		})
	})

	describe('updateTagSchema', () => {
		it('should validate valid tag update data', () => {
			const validData = {
				name: 'Personal',
				color: '#00FF00',
			}
			const result = updateTagSchema.safeParse(validData)
			expect(result.success).toBe(true)
		})

		it('should validate partial update with only color', () => {
			const validData = {
				color: '#00FF00',
			}
			const result = updateTagSchema.safeParse(validData)
			expect(result.success).toBe(true)
		})

		it('should validate empty update (all fields optional)', () => {
			const validData = {}
			const result = updateTagSchema.safeParse(validData)
			expect(result.success).toBe(true)
		})

		it('should reject invalid color format', () => {
			const invalidData = {
				color: 'red',
			}
			const result = updateTagSchema.safeParse(invalidData)
			expect(result.success).toBe(false)
		})

		it('should reject empty name when provided', () => {
			const invalidData = {
				name: '',
			}
			const result = updateTagSchema.safeParse(invalidData)
			expect(result.success).toBe(false)
		})

		it('should reject name longer than 50 characters', () => {
			const invalidData = {
				name: 'A'.repeat(51),
			}
			const result = updateTagSchema.safeParse(invalidData)
			expect(result.success).toBe(false)
		})

		it('should infer correct type', () => {
			const data: UpdateTagPayload = {
				name: 'Updated',
				color: '#00FF00',
			}
			expect(data.name).toBeDefined()
			expect(data.color).toBeDefined()
		})
	})
})

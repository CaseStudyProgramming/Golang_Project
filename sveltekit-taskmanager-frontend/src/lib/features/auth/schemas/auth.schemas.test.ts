import { describe, it, expect } from 'vitest'
import {
	loginSchema,
	registerSchema,
	forgotPasswordSchema,
	resetPasswordSchema,
	type LoginCredentials,
	type RegistrationData,
	type ForgotPasswordRequest,
	type ResetPasswordRequest,
} from './auth.schemas'

describe('Auth Schemas', () => {
	describe('loginSchema', () => {
		it('should validate valid login credentials', () => {
			const validData = {
				email: 'test@example.com',
				password: 'password123',
			}
			const result = loginSchema.safeParse(validData)
			expect(result.success).toBe(true)
		})

		it('should reject invalid email format', () => {
			const invalidData = {
				email: 'invalid-email',
				password: 'password123',
			}
			const result = loginSchema.safeParse(invalidData)
			expect(result.success).toBe(false)
		})

		it('should reject empty email', () => {
			const invalidData = {
				email: '',
				password: 'password123',
			}
			const result = loginSchema.safeParse(invalidData)
			expect(result.success).toBe(false)
		})

		it('should reject empty password', () => {
			const invalidData = {
				email: 'test@example.com',
				password: '',
			}
			const result = loginSchema.safeParse(invalidData)
			expect(result.success).toBe(false)
		})

		it('should infer correct type', () => {
			const data: LoginCredentials = {
				email: 'test@example.com',
				password: 'password123',
			}
			expect(data.email).toBeDefined()
			expect(data.password).toBeDefined()
		})
	})

	describe('registerSchema', () => {
		it('should validate valid registration data', () => {
			const validData = {
				email: 'test@example.com',
				name: 'Test User',
				password: 'Password123!',
			}
			const result = registerSchema.safeParse(validData)
			expect(result.success).toBe(true)
		})

		it('should validate registration without name', () => {
			const validData = {
				email: 'test@example.com',
				password: 'Password123!',
			}
			const result = registerSchema.safeParse(validData)
			expect(result.success).toBe(true)
		})

		it('should reject weak password without uppercase', () => {
			const invalidData = {
				email: 'test@example.com',
				name: 'Test User',
				password: 'password123!',
			}
			const result = registerSchema.safeParse(invalidData)
			expect(result.success).toBe(false)
		})

		it('should reject weak password without lowercase', () => {
			const invalidData = {
				email: 'test@example.com',
				name: 'Test User',
				password: 'PASSWORD123!',
			}
			const result = registerSchema.safeParse(invalidData)
			expect(result.success).toBe(false)
		})

		it('should reject weak password without number', () => {
			const invalidData = {
				email: 'test@example.com',
				name: 'Test User',
				password: 'Password!',
			}
			const result = registerSchema.safeParse(invalidData)
			expect(result.success).toBe(false)
		})

		it('should reject weak password without special character', () => {
			const invalidData = {
				email: 'test@example.com',
				name: 'Test User',
				password: 'Password123',
			}
			const result = registerSchema.safeParse(invalidData)
			expect(result.success).toBe(false)
		})

		it('should reject password shorter than 8 characters', () => {
			const invalidData = {
				email: 'test@example.com',
				name: 'Test User',
				password: 'Pass1!',
			}
			const result = registerSchema.safeParse(invalidData)
			expect(result.success).toBe(false)
		})

		it('should reject password longer than 128 characters', () => {
			const invalidData = {
				email: 'test@example.com',
				name: 'Test User',
				password: 'A'.repeat(129),
			}
			const result = registerSchema.safeParse(invalidData)
			expect(result.success).toBe(false)
		})

		it('should trim name', () => {
			const data = {
				email: 'test@example.com',
				name: '  Test User  ',
				password: 'Password123!',
			}
			const result = registerSchema.safeParse(data)
			expect(result.success).toBe(true)
			if (result.success) {
				expect(result.data.name).toBe('Test User')
			}
		})

		it('should infer correct type', () => {
			const data: RegistrationData = {
				email: 'test@example.com',
				name: 'Test User',
				password: 'Password123!',
			}
			expect(data.email).toBeDefined()
			expect(data.password).toBeDefined()
		})
	})

	describe('forgotPasswordSchema', () => {
		it('should validate valid forgot password request', () => {
			const validData = {
				email: 'test@example.com',
			}
			const result = forgotPasswordSchema.safeParse(validData)
			expect(result.success).toBe(true)
		})

		it('should reject invalid email format', () => {
			const invalidData = {
				email: 'invalid-email',
			}
			const result = forgotPasswordSchema.safeParse(invalidData)
			expect(result.success).toBe(false)
		})

		it('should reject empty email', () => {
			const invalidData = {
				email: '',
			}
			const result = forgotPasswordSchema.safeParse(invalidData)
			expect(result.success).toBe(false)
		})

		it('should infer correct type', () => {
			const data: ForgotPasswordRequest = {
				email: 'test@example.com',
			}
			expect(data.email).toBeDefined()
		})
	})

	describe('resetPasswordSchema', () => {
		it('should validate valid reset password request', () => {
			const validData = {
				token: 'valid-token',
				password: 'Password123!',
				confirmPassword: 'Password123!',
			}
			const result = resetPasswordSchema.safeParse(validData)
			expect(result.success).toBe(true)
		})

		it('should reject mismatched passwords', () => {
			const invalidData = {
				token: 'valid-token',
				password: 'Password123!',
				confirmPassword: 'DifferentPassword123!',
			}
			const result = resetPasswordSchema.safeParse(invalidData)
			expect(result.success).toBe(false)
		})

		it('should reject empty token', () => {
			const invalidData = {
				token: '',
				password: 'Password123!',
				confirmPassword: 'Password123!',
			}
			const result = resetPasswordSchema.safeParse(invalidData)
			expect(result.success).toBe(false)
		})

		it('should reject empty password', () => {
			const invalidData = {
				token: 'valid-token',
				password: '',
				confirmPassword: '',
			}
			const result = resetPasswordSchema.safeParse(invalidData)
			expect(result.success).toBe(false)
		})

		it('should reject weak password', () => {
			const invalidData = {
				token: 'valid-token',
				password: 'weak',
				confirmPassword: 'weak',
			}
			const result = resetPasswordSchema.safeParse(invalidData)
			expect(result.success).toBe(false)
		})

		it('should infer correct type', () => {
			const data: ResetPasswordRequest = {
				token: 'valid-token',
				password: 'Password123!',
				confirmPassword: 'Password123!',
			}
			expect(data.token).toBeDefined()
			expect(data.password).toBeDefined()
			expect(data.confirmPassword).toBeDefined()
		})
	})
})

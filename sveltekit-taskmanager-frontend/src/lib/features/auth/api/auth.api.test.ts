import { beforeEach, describe, expect, it, vi } from 'vitest'

import type { User } from '../types/auth.types'

import { authApi } from './auth.api'

// Mock the HTTP client
vi.mock('$lib/shared/utils/api.utils', () => ({
	httpClient: {
		delete: vi.fn(),
		get: vi.fn(),
		patch: vi.fn(),
		post: vi.fn(),
	},
}))

const { httpClient } = (await import('$lib/shared/utils/api.utils')) as unknown as {
	httpClient: {
		delete: ReturnType<typeof vi.fn>
		get: ReturnType<typeof vi.fn>
		patch: ReturnType<typeof vi.fn>
		post: ReturnType<typeof vi.fn>
	}
}

describe('Auth API', () => {
	beforeEach(() => {
		vi.clearAllMocks()
	})

	describe('Login', () => {
		it('calls login endpoint with credentials', async () => {
			const mockUser: User = {
				email: 'test@example.com',
				id: '1',
				name: 'Test User',
			}
			const mockResponse = {
				expiresIn: 3600,
				refreshToken: 'refresh-token',
				token: 'auth-token',
				user: mockUser,
			}

			httpClient.post.mockResolvedValue(mockResponse)

			const result = await authApi.login({
				email: 'test@example.com',
				password: 'password',
			})

			expect(httpClient.post).toHaveBeenCalledWith('/auth/login', {
				email: 'test@example.com',
				password: 'password',
			})
			expect(result).toEqual(mockResponse)
		})

		it('handles login errors', async () => {
			httpClient.post.mockRejectedValue(new Error('Invalid credentials'))

			await expect(authApi.login({ email: 'test@example.com', password: 'wrong' })).rejects.toThrow(
				'Invalid credentials'
			)
		})
	})

	describe('Register', () => {
		it('calls register endpoint with user data', async () => {
			const mockUser: User = {
				email: 'test@example.com',
				id: '1',
				name: 'Test User',
			}
			const mockResponse = {
				expiresIn: 3600,
				refreshToken: 'refresh-token',
				token: 'auth-token',
				user: mockUser,
			}

			httpClient.post.mockResolvedValue(mockResponse)

			const result = await authApi.register({
				email: 'test@example.com',
				name: 'Test User',
				password: 'password',
			})

			expect(httpClient.post).toHaveBeenCalledWith('/auth/register', {
				email: 'test@example.com',
				name: 'Test User',
				password: 'password',
			})
			expect(result).toEqual(mockResponse)
		})

		it('handles registration errors', async () => {
			httpClient.post.mockRejectedValue(new Error('Email already exists'))

			await expect(
				authApi.register({ email: 'test@example.com', password: 'password' })
			).rejects.toThrow('Email already exists')
		})
	})

	describe('Logout', () => {
		it('calls logout endpoint', async () => {
			httpClient.post.mockResolvedValue(undefined)

			await authApi.logout()

			expect(httpClient.post).toHaveBeenCalledWith('/auth/logout')
		})

		it('handles logout errors gracefully', async () => {
			httpClient.post.mockRejectedValue(new Error('Network error'))

			await expect(authApi.logout()).rejects.toThrow('Network error')
		})
	})

	describe('Get Current User', () => {
		it('calls current user endpoint', async () => {
			const mockUser: User = {
				email: 'test@example.com',
				id: '1',
				name: 'Test User',
			}

			httpClient.get.mockResolvedValue(mockUser)

			const result = await authApi.getCurrentUser()

			expect(httpClient.get).toHaveBeenCalledWith('/auth/me')
			expect(result).toEqual(mockUser)
		})

		it('handles get current user errors', async () => {
			httpClient.get.mockRejectedValue(new Error('Unauthorized'))

			await expect(authApi.getCurrentUser()).rejects.toThrow('Unauthorized')
		})
	})

	describe('Refresh Token', () => {
		it('calls refresh token endpoint', async () => {
			const mockResponse = {
				expiresIn: 3600,
				refreshToken: 'new-refresh-token',
				token: 'new-token',
			}

			httpClient.post.mockResolvedValue(mockResponse)

			const result = await authApi.refreshToken('refresh-token')

			expect(httpClient.post).toHaveBeenCalledWith('/auth/refresh', {
				refreshToken: 'refresh-token',
			})
			expect(result).toEqual(mockResponse)
		})

		it('handles refresh token errors', async () => {
			httpClient.post.mockRejectedValue(new Error('Invalid refresh token'))

			await expect(authApi.refreshToken('invalid-token')).rejects.toThrow('Invalid refresh token')
		})
	})

	describe('Forgot Password', () => {
		it('calls forgot password endpoint', async () => {
			const mockResponse = { message: 'Password reset email sent' }

			httpClient.post.mockResolvedValue(mockResponse)

			const result = await authApi.forgotPassword({ email: 'test@example.com' })

			expect(httpClient.post).toHaveBeenCalledWith('/auth/forgot-password', {
				email: 'test@example.com',
			})
			expect(result).toEqual(mockResponse)
		})

		it('handles forgot password errors', async () => {
			httpClient.post.mockRejectedValue(new Error('Email not found'))

			await expect(authApi.forgotPassword({ email: 'nonexistent@example.com' })).rejects.toThrow(
				'Email not found'
			)
		})
	})

	describe('Reset Password', () => {
		it('calls reset password endpoint', async () => {
			const mockResponse = { message: 'Password reset successful' }

			httpClient.post.mockResolvedValue(mockResponse)

			const result = await authApi.resetPassword({
				confirmPassword: 'new-password',
				password: 'new-password',
				token: 'reset-token',
			})

			expect(httpClient.post).toHaveBeenCalledWith('/auth/reset-password', {
				confirmPassword: 'new-password',
				password: 'new-password',
				token: 'reset-token',
			})
			expect(result).toEqual(mockResponse)
		})

		it('handles reset password errors', async () => {
			httpClient.post.mockRejectedValue(new Error('Invalid token'))

			await expect(
				authApi.resetPassword({
					confirmPassword: 'new-password',
					password: 'new-password',
					token: 'invalid-token',
				})
			).rejects.toThrow('Invalid token')
		})
	})
})

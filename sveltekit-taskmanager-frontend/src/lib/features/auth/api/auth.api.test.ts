import { beforeEach, describe, expect, it, vi } from 'vitest';

import { authApi } from './auth.api';
import type { User } from '../types/auth.types';

// Mock the HTTP client
const mockHttpClient = {
	get: vi.fn(),
	post: vi.fn(),
	patch: vi.fn(),
	delete: vi.fn()
};

vi.mock('$lib/shared/utils/api.utils', () => ({
	httpClient: mockHttpClient
}));

describe('Auth API', () => {
	beforeEach(() => {
		vi.clearAllMocks();
	});

	describe('Login', () => {
		it('calls login endpoint with credentials', async () => {
			const mockUser: User = {
				email: 'test@example.com',
				id: '1',
				name: 'Test User'
			};
			const mockResponse = {
				user: mockUser,
				token: 'auth-token',
				refreshToken: 'refresh-token',
				expiresIn: 3600
			};

			mockHttpClient.post.mockResolvedValue(mockResponse);

			const result = await authApi.login({
				email: 'test@example.com',
				password: 'password'
			});

			expect(mockHttpClient.post).toHaveBeenCalledWith('/auth/login', {
				email: 'test@example.com',
				password: 'password'
			});
			expect(result).toEqual(mockResponse);
		});

		it('handles login errors', async () => {
			mockHttpClient.post.mockRejectedValue(new Error('Invalid credentials'));

			await expect(
				authApi.login({ email: 'test@example.com', password: 'wrong' })
			).rejects.toThrow('Invalid credentials');
		});
	});

	describe('Register', () => {
		it('calls register endpoint with user data', async () => {
			const mockUser: User = {
				email: 'test@example.com',
				id: '1',
				name: 'Test User'
			};
			const mockResponse = {
				user: mockUser,
				token: 'auth-token',
				refreshToken: 'refresh-token',
				expiresIn: 3600
			};

			mockHttpClient.post.mockResolvedValue(mockResponse);

			const result = await authApi.register({
				email: 'test@example.com',
				password: 'password',
				name: 'Test User'
			});

			expect(mockHttpClient.post).toHaveBeenCalledWith('/auth/register', {
				email: 'test@example.com',
				password: 'password',
				name: 'Test User'
			});
			expect(result).toEqual(mockResponse);
		});

		it('handles registration errors', async () => {
			mockHttpClient.post.mockRejectedValue(new Error('Email already exists'));

			await expect(
				authApi.register({ email: 'test@example.com', password: 'password' })
			).rejects.toThrow('Email already exists');
		});
	});

	describe('Logout', () => {
		it('calls logout endpoint', async () => {
			mockHttpClient.post.mockResolvedValue(undefined);

			await authApi.logout();

			expect(mockHttpClient.post).toHaveBeenCalledWith('/auth/logout');
		});

		it('handles logout errors gracefully', async () => {
			mockHttpClient.post.mockRejectedValue(new Error('Network error'));

			await expect(authApi.logout()).rejects.toThrow('Network error');
		});
	});

	describe('Get Current User', () => {
		it('calls current user endpoint', async () => {
			const mockUser: User = {
				email: 'test@example.com',
				id: '1',
				name: 'Test User'
			};

			mockHttpClient.get.mockResolvedValue(mockUser);

			const result = await authApi.getCurrentUser();

			expect(mockHttpClient.get).toHaveBeenCalledWith('/auth/me');
			expect(result).toEqual(mockUser);
		});

		it('handles get current user errors', async () => {
			mockHttpClient.get.mockRejectedValue(new Error('Unauthorized'));

			await expect(authApi.getCurrentUser()).rejects.toThrow('Unauthorized');
		});
	});

	describe('Refresh Token', () => {
		it('calls refresh token endpoint', async () => {
			const mockResponse = {
				token: 'new-token',
				refreshToken: 'new-refresh-token',
				expiresIn: 3600
			};

			mockHttpClient.post.mockResolvedValue(mockResponse);

			const result = await authApi.refreshToken('refresh-token');

			expect(mockHttpClient.post).toHaveBeenCalledWith('/auth/refresh', {
				refreshToken: 'refresh-token'
			});
			expect(result).toEqual(mockResponse);
		});

		it('handles refresh token errors', async () => {
			mockHttpClient.post.mockRejectedValue(new Error('Invalid refresh token'));

			await expect(authApi.refreshToken('invalid-token')).rejects.toThrow(
				'Invalid refresh token'
			);
		});
	});

	describe('Forgot Password', () => {
		it('calls forgot password endpoint', async () => {
			const mockResponse = { message: 'Password reset email sent' };

			mockHttpClient.post.mockResolvedValue(mockResponse);

			const result = await authApi.forgotPassword({ email: 'test@example.com' });

			expect(mockHttpClient.post).toHaveBeenCalledWith('/auth/forgot-password', {
				email: 'test@example.com'
			});
			expect(result).toEqual(mockResponse);
		});

		it('handles forgot password errors', async () => {
			mockHttpClient.post.mockRejectedValue(new Error('Email not found'));

			await expect(
				authApi.forgotPassword({ email: 'nonexistent@example.com' })
			).rejects.toThrow('Email not found');
		});
	});

	describe('Reset Password', () => {
		it('calls reset password endpoint', async () => {
			const mockResponse = { message: 'Password reset successful' };

			mockHttpClient.post.mockResolvedValue(mockResponse);

			const result = await authApi.resetPassword({
				token: 'reset-token',
				password: 'new-password',
				confirmPassword: 'new-password'
			});

			expect(mockHttpClient.post).toHaveBeenCalledWith('/auth/reset-password', {
				token: 'reset-token',
				password: 'new-password',
				confirmPassword: 'new-password'
			});
			expect(result).toEqual(mockResponse);
		});

		it('handles reset password errors', async () => {
			mockHttpClient.post.mockRejectedValue(new Error('Invalid token'));

			await expect(
				authApi.resetPassword({
					token: 'invalid-token',
					password: 'new-password',
					confirmPassword: 'new-password'
				})
			).rejects.toThrow('Invalid token');
		});
	});
});

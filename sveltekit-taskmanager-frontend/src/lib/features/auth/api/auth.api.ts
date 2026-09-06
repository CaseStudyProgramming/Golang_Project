/**
 * Authentication API client
 * Handles all authentication-related API calls
 */

import { httpClient } from '$lib/shared/utils/api.utils';
import type { User } from '../types/auth.types';

/**
 * Login response interface
 */
interface LoginResponse {
	user: User;
	token: string;
	refreshToken?: string;
	expiresIn?: number;
}

/**
 * Register response interface
 */
interface RegisterResponse {
	user: User;
	token: string;
	refreshToken?: string;
	expiresIn?: number;
}

/**
 * Login credentials interface
 */
interface LoginCredentials {
	email: string;
	password: string;
}

/**
 * Registration data interface
 */
interface RegistrationData {
	email: string;
	password: string;
	name?: string;
}

/**
 * Forgot password request interface
 */
interface ForgotPasswordRequest {
	email: string;
}

/**
 * Reset password request interface
 */
interface ResetPasswordRequest {
	token: string;
	password: string;
	confirmPassword: string;
}

/**
 * Authentication API client
 */
export const authApi = {
	/**
	 * Login user with credentials
	 */
	async login(credentials: LoginCredentials): Promise<LoginResponse> {
		return httpClient.post<LoginResponse>('/auth/login', credentials);
	},

	/**
	 * Register new user
	 */
	async register(data: RegistrationData): Promise<RegisterResponse> {
		return httpClient.post<RegisterResponse>('/auth/register', data);
	},

	/**
	 * Logout user
	 */
	async logout(): Promise<void> {
		return httpClient.post<void>('/auth/logout');
	},

	/**
	 * Refresh JWT token
	 */
	async refreshToken(refreshToken: string): Promise<{ token: string; refreshToken?: string; expiresIn?: number }> {
		return httpClient.post<{ token: string; refreshToken?: string; expiresIn?: number }>('/auth/refresh', { refreshToken });
	},

	/**
	 * Get current user
	 */
	async getCurrentUser(): Promise<User> {
		return httpClient.get<User>('/auth/me');
	},

	/**
	 * Request password reset
	 */
	async forgotPassword(data: ForgotPasswordRequest): Promise<{ message: string }> {
		return httpClient.post<{ message: string }>('/auth/forgot-password', data);
	},

	/**
	 * Reset password with token
	 */
	async resetPassword(data: ResetPasswordRequest): Promise<{ message: string }> {
		return httpClient.post<{ message: string }>('/auth/reset-password', data);
	}
};

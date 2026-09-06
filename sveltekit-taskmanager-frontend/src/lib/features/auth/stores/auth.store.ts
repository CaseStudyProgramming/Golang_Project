/**
 * Authentication store using Svelte 5 runes
 * Manages user session, JWT token, and authentication state
 */

import { setAuthToken, removeAuthToken, getAuthToken } from '$lib/shared/utils/auth.interceptors';
import { AuthenticationError, withErrorHandling } from '$lib/shared/utils/error.utils';
import type { User, AuthState, LoginCredentials, RegistrationData } from '../types/auth.types';

/**
 * Create authentication store with Svelte 5 runes
 */
function createAuthStore() {
	const state = $state<AuthState>({
		isAuthenticated: false,
		user: null,
		token: null,
		isLoading: false,
		error: null
	});

	/**
	 * Initialize authentication state from localStorage
	 */
	function initialize(): void {
		const token = getAuthToken();
		if (token) {
			state.token = token;
			state.isAuthenticated = true;
			// In a real app, you might want to validate the token
			// and fetch user data here
		}
	}

	/**
	 * Login user with credentials
	 */
	async function login(credentials: LoginCredentials): Promise<void> {
		state.isLoading = true;
		state.error = null;

		try {
			await withErrorHandling(async () => {
				// This would be replaced with actual API call
				// const response = await httpClient.post<{ user: User; token: string }>('/auth/login', credentials);
				
				// Mock response for development
				const mockUser: User = {
					id: '1',
					email: credentials.email,
					name: 'Test User',
					role: 'user'
				};
				const mockToken = 'mock_jwt_token_' + Date.now();

				state.user = mockUser;
				state.token = mockToken;
				state.isAuthenticated = true;

				setAuthToken(mockToken);
			}, 'Login failed');
		} catch (error) {
			state.error = error instanceof Error ? error.message : 'Login failed';
			throw error;
		} finally {
			state.isLoading = false;
		}
	}

	/**
	 * Register new user
	 */
	async function register(data: RegistrationData): Promise<void> {
		state.isLoading = true;
		state.error = null;

		try {
			await withErrorHandling(async () => {
				// This would be replaced with actual API call
				// const response = await httpClient.post<{ user: User; token: string }>('/auth/register', data);
				
				// Mock response for development
				const mockUser: User = {
					id: '1',
					email: data.email,
					name: data.name || 'New User',
					role: 'user'
				};
				const mockToken = 'mock_jwt_token_' + Date.now();

				state.user = mockUser;
				state.token = mockToken;
				state.isAuthenticated = true;

				setAuthToken(mockToken);
			}, 'Registration failed');
		} catch (error) {
			state.error = error instanceof Error ? error.message : 'Registration failed';
			throw error;
		} finally {
			state.isLoading = false;
		}
	}

	/**
	 * Logout user
	 */
	function logout(): void {
		state.user = null;
		state.token = null;
		state.isAuthenticated = false;
		state.error = null;

		removeAuthToken();
	}

	/**
	 * Update user data
	 */
	function updateUser(user: Partial<User>): void {
		if (state.user) {
			state.user = { ...state.user, ...user };
		}
	}

	/**
	 * Clear error state
	 */
	function clearError(): void {
		state.error = null;
	}

	return {
		get state() {
			return state;
		},
		initialize,
		login,
		register,
		logout,
		updateUser,
		clearError
	};
}

/**
 * Export authentication store instance
 */
export const authStore = createAuthStore();

/**
 * Initialize auth store on app load
 */
if (typeof window !== 'undefined') {
	authStore.initialize();
}

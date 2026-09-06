/**
 * Authentication store using Svelte 5 runes
 * Manages user session, JWT token, and authentication state
 */

import { setAuthToken, removeAuthToken, getAuthToken, setRefreshToken, getRefreshToken, isTokenExpired } from '$lib/shared/utils/auth.interceptors';
import { AuthenticationError, withErrorHandling, ValidationError } from '$lib/shared/utils/error.utils';
import { loginSchema, registerSchema } from '../schemas/auth.schemas';
import { authApi } from '../api/auth.api';
import type { User, AuthState } from '../types/auth.types';

let refreshTimer: ReturnType<typeof setTimeout> | null = null;

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
			
			// Set up automatic token refresh
			setupTokenRefresh();
			
			// Fetch current user data
			fetchCurrentUser().catch(() => {
				// If fetch fails, token might be invalid
				logout();
			});
		}
	}

	/**
	 * Setup automatic token refresh
	 */
	function setupTokenRefresh(): void {
		// Clear any existing timer
		if (refreshTimer) {
			clearTimeout(refreshTimer);
		}

		// Check token every minute
		refreshTimer = setInterval(() => {
			if (state.token && isTokenExpired()) {
				refreshToken().catch(() => {
					// If refresh fails, logout
					logout();
				});
			}
		}, 60000); // Check every minute
	}

	/**
	 * Cleanup authentication state
	 */
	function cleanup(): void {
		if (refreshTimer) {
			clearInterval(refreshTimer);
			refreshTimer = null;
		}
	}

	/**
	 * Login user with credentials
	 */
	async function login(credentials: { email: string; password: string }): Promise<void> {
		state.isLoading = true;
		state.error = null;

		try {
			// Validate input with Zod
			const validatedCredentials = loginSchema.parse(credentials);

			await withErrorHandling(async () => {
				const response = await authApi.login(validatedCredentials);

				state.user = response.user;
				state.token = response.token;
				state.isAuthenticated = true;

				setAuthToken(response.token, response.expiresIn);
				if (response.refreshToken) {
					setRefreshToken(response.refreshToken);
				}
			}, 'Login failed');
		} catch (error) {
			if (error instanceof Error && error.name === 'ZodError') {
				state.error = 'Invalid input: ' + error.message;
				throw new ValidationError('credentials', error.message);
			}
			state.error = error instanceof Error ? error.message : 'Login failed';
			throw error;
		} finally {
			state.isLoading = false;
		}
	}

	/**
	 * Register new user
	 */
	async function register(data: { email: string; password: string; name?: string }): Promise<void> {
		state.isLoading = true;
		state.error = null;

		try {
			// Validate input with Zod
			const validatedData = registerSchema.parse(data);

			await withErrorHandling(async () => {
				const response = await authApi.register(validatedData);

				state.user = response.user;
				state.token = response.token;
				state.isAuthenticated = true;

				setAuthToken(response.token, response.expiresIn);
				if (response.refreshToken) {
					setRefreshToken(response.refreshToken);
				}
			}, 'Registration failed');
		} catch (error) {
			if (error instanceof Error && error.name === 'ZodError') {
				state.error = 'Invalid input: ' + error.message;
				throw new ValidationError('registration', error.message);
			}
			state.error = error instanceof Error ? error.message : 'Registration failed';
			throw error;
		} finally {
			state.isLoading = false;
		}
	}

	/**
	 * Logout user
	 */
	async function logout(): Promise<void> {
		try {
			await authApi.logout();
		} catch (error) {
			console.error('Logout API call failed:', error);
		} finally {
			state.user = null;
			state.token = null;
			state.isAuthenticated = false;
			state.error = null;

			removeAuthToken();
			cleanup();
		}
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

	/**
	 * Refresh JWT token
	 */
	async function refreshToken(): Promise<void> {
		const refreshTokenValue = getRefreshToken();
		if (!refreshTokenValue) {
			throw new AuthenticationError('No refresh token available');
		}

		try {
			const response = await authApi.refreshToken(refreshTokenValue);
			state.token = response.token;
			setAuthToken(response.token, response.expiresIn);
			if (response.refreshToken) {
				setRefreshToken(response.refreshToken);
			}
		} catch (error) {
			console.error('Token refresh failed:', error);
			logout();
			throw error;
		}
	}

	/**
	 * Fetch current user data
	 */
	async function fetchCurrentUser(): Promise<void> {
		if (!state.token) return;

		state.isLoading = true;
		try {
			const user = await authApi.getCurrentUser();
			state.user = user;
			state.isAuthenticated = true;
		} catch (error) {
			console.error('Failed to fetch current user:', error);
			state.error = error instanceof Error ? error.message : 'Failed to fetch user';
			throw error;
		} finally {
			state.isLoading = false;
		}
	}

	/**
	 * Check if token needs refresh
	 */
	function needsTokenRefresh(): boolean {
		return isTokenExpired();
	}

	/**
	 * Ensure valid token (refresh if needed)
	 */
	async function ensureValidToken(): Promise<void> {
		if (!state.token) return;

		if (isTokenExpired()) {
			await refreshToken();
		}
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
		clearError,
		refreshToken,
		fetchCurrentUser,
		needsTokenRefresh,
		ensureValidToken,
		cleanup
	};
}

/**
 * Export authentication store instance
 * Only create store instance on client side to avoid SSR issues
 */
let authStoreInstance: ReturnType<typeof createAuthStore> | null = null;

export const authStore = new Proxy({} as ReturnType<typeof createAuthStore>, {
	get(_target, prop) {
		if (!authStoreInstance) {
			if (typeof window === 'undefined') {
				throw new Error('authStore can only be accessed on the client side');
			}
			authStoreInstance = createAuthStore();
			authStoreInstance.initialize();
		}
		return authStoreInstance[prop as keyof ReturnType<typeof createAuthStore>];
	}
});

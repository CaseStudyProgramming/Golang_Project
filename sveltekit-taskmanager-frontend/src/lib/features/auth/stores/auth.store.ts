/**
 * Authentication store using Svelte 5 runes
 * Manages user session, JWT token, and authentication state
 */

import { goto } from '$app/navigation';
import {
	getAuthToken,
	getRefreshToken,
	isTokenExpired,
	removeAuthToken,
	setAuthToken,
	setRefreshToken
} from '$lib/shared/utils/auth.interceptors';
import {
	AuthenticationError,
	ValidationError,
	withErrorHandling
} from '$lib/shared/utils/error.utils';

import type { AuthState, User } from '../types/auth.types';

import { authApi } from '../api/auth.api';
import { loginSchema, registerSchema } from '../schemas/auth.schemas';

let refreshTimer: null | ReturnType<typeof setTimeout> = null;
const REDIRECT_KEY = 'auth_redirect';

/**
 * Create authentication store with Svelte 5 runes
 */
function createAuthStore() {
	const state = $state<AuthState>({
		error: null,
		isAuthenticated: false,
		isLoading: false,
		token: null,
		user: null
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
	 * Save redirect URL for after authentication
	 */
	function saveRedirectUrl(url: string): void {
		if (typeof window === 'undefined') return;
		sessionStorage.setItem(REDIRECT_KEY, url);
	}

	/**
	 * Get and clear redirect URL
	 */
	function getRedirectUrl(): null | string {
		if (typeof window === 'undefined') return null;
		const url = sessionStorage.getItem(REDIRECT_KEY);
		sessionStorage.removeItem(REDIRECT_KEY);
		return url;
	}

	/**
	 * Perform redirect after authentication
	 */
	async function redirectAfterAuth(): Promise<void> {
		const redirectUrl = getRedirectUrl();
		if (redirectUrl) {
			await goto(redirectUrl);
		} else {
			await goto('/dashboard');
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
	async function register(data: { email: string; name?: string; password: string }): Promise<void> {
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
		cleanup,
		clearError,
		ensureValidToken,
		fetchCurrentUser,
		getRedirectUrl,
		initialize,
		login,
		logout,
		needsTokenRefresh,
		redirectAfterAuth,
		refreshToken,
		register,
		saveRedirectUrl,
		get state() {
			return state;
		},
		updateUser
	};
}

/**
 * Export authentication store instance
 * Only create store instance on client side to avoid SSR issues
 */
let authStoreInstance: null | ReturnType<typeof createAuthStore> = null;

// Create a safe SSR-compatible default store
function createSSRAuthStore() {
	const state: AuthState = {
		error: null,
		isAuthenticated: false,
		isLoading: false,
		token: null,
		user: null
	};

	return {
		cleanup: () => {},
		clearError: () => {},
		ensureValidToken: async () => {},
		fetchCurrentUser: async () => {},
		getRedirectUrl: () => null,
		initialize: () => {},
		login: async () => { throw new Error('Auth store not available during SSR'); },
		logout: async () => { throw new Error('Auth store not available during SSR'); },
		needsTokenRefresh: () => false,
		redirectAfterAuth: async () => {},
		refreshToken: async () => {},
		register: async () => { throw new Error('Auth store not available during SSR'); },
		saveRedirectUrl: () => {},
		get state() {
			return state;
		},
		updateUser: () => {}
	};
}

export const authStore = new Proxy({} as ReturnType<typeof createAuthStore>, {
	get(_target, prop) {
		if (!authStoreInstance) {
			if (typeof window === 'undefined') {
				// Return safe SSR-compatible store during server-side rendering
				authStoreInstance = createSSRAuthStore();
			} else {
				authStoreInstance = createAuthStore();
				authStoreInstance.initialize();
			}
		}
		return authStoreInstance![prop as keyof ReturnType<typeof createAuthStore>];
	}
});

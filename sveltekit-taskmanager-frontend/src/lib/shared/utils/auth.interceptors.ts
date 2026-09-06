/**
 * Authentication interceptors for HTTP client
 */

const TOKEN_KEY = 'auth_token';
const REFRESH_TOKEN_KEY = 'refresh_token';
const TOKEN_EXPIRY_KEY = 'token_expiry';

/**
 * Get authentication token from localStorage
 */
export function getAuthToken(): string | null {
	if (typeof window === 'undefined') return null;
	const token = localStorage.getItem(TOKEN_KEY);
	
	// Check if token is expired
	const expiry = localStorage.getItem(TOKEN_EXPIRY_KEY);
	if (expiry && token) {
		const expiryTime = parseInt(expiry, 10);
		if (Date.now() > expiryTime) {
			removeAuthToken();
			return null;
		}
	}
	
	return token;
}

/**
 * Get refresh token from localStorage
 */
export function getRefreshToken(): string | null {
	if (typeof window === 'undefined') return null;
	return localStorage.getItem(REFRESH_TOKEN_KEY);
}

/**
 * Set authentication token in localStorage
 */
export function setAuthToken(token: string, expiresIn: number = 3600): void {
	if (typeof window === 'undefined') return;
	localStorage.setItem(TOKEN_KEY, token);
	
	// Set token expiry (default 1 hour)
	const expiryTime = Date.now() + expiresIn * 1000;
	localStorage.setItem(TOKEN_EXPIRY_KEY, expiryTime.toString());
	
	// Also set as cookie for server-side access (httpOnly for security would be better)
	document.cookie = `auth_token=${token}; path=/; max-age=${expiresIn}; SameSite=Strict`;
}

/**
 * Set refresh token in localStorage
 */
export function setRefreshToken(token: string): void {
	if (typeof window === 'undefined') return;
	localStorage.setItem(REFRESH_TOKEN_KEY, token);
}

/**
 * Remove authentication token from localStorage
 */
export function removeAuthToken(): void {
	if (typeof window === 'undefined') return;
	localStorage.removeItem(TOKEN_KEY);
	localStorage.removeItem(REFRESH_TOKEN_KEY);
	localStorage.removeItem(TOKEN_EXPIRY_KEY);
	
	// Also remove cookie
	document.cookie = 'auth_token=; path=/; max-age=0; SameSite=Strict';
}

/**
 * Check if token is expired
 */
export function isTokenExpired(): boolean {
	const expiry = localStorage.getItem(TOKEN_EXPIRY_KEY);
	if (!expiry) return true;
	
	const expiryTime = parseInt(expiry, 10);
	return Date.now() > expiryTime;
}

/**
 * Request interceptor to add authentication headers
 */
export function authRequestInterceptor(request: RequestInit): RequestInit {
	const token = getAuthToken();
	if (token) {
		return {
			...request,
			headers: {
				...request.headers,
				Authorization: `Bearer ${token}`
			}
		};
	}
	return request;
}

/**
 * Response interceptor to handle authentication errors
 */
export async function authResponseInterceptor(response: Response): Promise<Response> {
	// Handle 401 Unauthorized - token might be expired
	if (response.status === 401) {
		removeAuthToken();
		// Redirect to login page if on client
		if (typeof window !== 'undefined') {
			window.location.href = '/auth/login';
		}
	}

	// Handle 403 Forbidden - insufficient permissions
	if (response.status === 403) {
		// Could redirect to a "not authorized" page
		if (typeof window !== 'undefined') {
			window.location.href = '/unauthorized';
		}
	}

	return response;
}

/**
 * Response interceptor for error logging
 */
export function errorLoggingInterceptor(response: Response): Response {
	if (!response.ok) {
		console.error(`API Error: ${response.status} ${response.statusText}`, {
			url: response.url,
			status: response.status,
			statusText: response.statusText
		});
	}
	return response;
}

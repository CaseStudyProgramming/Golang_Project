/**
 * Authentication interceptors for HTTP client
 */

/**
 * Get authentication token from localStorage
 */
export function getAuthToken(): string | null {
	if (typeof window === 'undefined') return null;
	return localStorage.getItem('auth_token');
}

/**
 * Set authentication token in localStorage
 */
export function setAuthToken(token: string): void {
	if (typeof window === 'undefined') return;
	localStorage.setItem('auth_token', token);
}

/**
 * Remove authentication token from localStorage
 */
export function removeAuthToken(): void {
	if (typeof window === 'undefined') return;
	localStorage.removeItem('auth_token');
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
			window.location.href = '/login';
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

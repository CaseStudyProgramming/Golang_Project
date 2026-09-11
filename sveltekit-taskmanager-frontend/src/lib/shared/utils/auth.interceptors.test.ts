import { beforeEach, describe, expect, it, vi } from 'vitest';

import {
	authResponseInterceptor,
	errorLoggingInterceptor
} from './auth.interceptors';

describe('Auth Interceptors', () => {
	beforeEach(() => {
		vi.clearAllMocks();
	});

	describe('authResponseInterceptor', () => {
		it('removes token on 401 response', async () => {
			const response = new Response(null, { status: 401, statusText: 'Unauthorized' });
			const result = await authResponseInterceptor(response);
			
			if (globalThis.localStorage) {
				expect(globalThis.localStorage.removeItem).toHaveBeenCalledWith('auth_token');
				expect(globalThis.localStorage.removeItem).toHaveBeenCalledWith('refresh_token');
				expect(globalThis.localStorage.removeItem).toHaveBeenCalledWith('token_expiry');
			}
			expect(result).toBe(response);
		});

		it('handles 403 Forbidden response', async () => {
			const response = new Response(null, { status: 403, statusText: 'Forbidden' });
			const result = await authResponseInterceptor(response);
			
			if (globalThis.window?.location) {
				expect(globalThis.window.location.href).toBe('/unauthorized');
			}
			expect(result).toBe(response);
		});

		it('passes through successful responses', async () => {
			const response = new Response(null, { status: 200, statusText: 'OK' });
			const result = await authResponseInterceptor(response);
			
			expect(result).toBe(response);
		});
	});

	describe('errorLoggingInterceptor', () => {
		it('logs error responses', () => {
			const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {});
			
			const response = new Response('{"error": "test"}', { 
				status: 500, 
				statusText: 'Internal Server Error' 
			});
			Object.defineProperty(response, 'url', { value: 'http://test.com/api', writable: true });
			
			const result = errorLoggingInterceptor(response);
			
			expect(consoleSpy).toHaveBeenCalledWith(
				'API Error: 500 Internal Server Error',
				expect.objectContaining({
					status: 500,
					statusText: 'Internal Server Error',
					url: 'http://test.com/api'
				})
			);
			expect(result).toBe(response);
			
			consoleSpy.mockRestore();
		});

		it('does not log successful responses', () => {
			const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {});
			
			const response = new Response(null, { status: 200, statusText: 'OK' });
			const result = errorLoggingInterceptor(response);
			
			expect(consoleSpy).not.toHaveBeenCalled();
			expect(result).toBe(response);
			
			consoleSpy.mockRestore();
		});
	});
});
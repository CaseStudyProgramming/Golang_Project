import { serverEnv } from '$lib/server/env';
import { redirect } from '@sveltejs/kit';
import type { Handle } from '@sveltejs/kit';

// Import serverEnv at the top to trigger fail-fast validation on app startup
console.log('Server environment validated successfully', serverEnv);

/**
 * Handle incoming requests
 * Implements authentication checks and session management
 */
export const handle: Handle = async ({ event, resolve }) => {
	const authToken = event.cookies.get('auth_token');
	const path = event.url.pathname;

	// Protect dashboard routes
	if (path.startsWith('/dashboard') && !authToken) {
		redirect(302, '/auth/login');
	}

	// Redirect authenticated users away from auth pages
	if ((path === '/auth/login' || path === '/auth/register') && authToken) {
		redirect(302, '/dashboard');
	}

	const response = await resolve(event);
	return response;
};

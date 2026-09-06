/**
 * Protected route layout server load function
 * Implements navigation guard for authentication
 */

import { redirect } from '@sveltejs/kit';
import type { LayoutServerLoad } from './$types';

export const load: LayoutServerLoad = async ({ cookies, url }) => {
	// Check for authentication token in cookies
	const authToken = cookies.get('auth_token');

	// If no token, redirect to login with return URL
	if (!authToken) {
		throw redirect(302, `/auth/login?redirectTo=${url.pathname}`);
	}

	// You could also validate the token here by calling your backend
	// For now, we'll just check for its existence

	return {
		isAuthenticated: true
	};
};

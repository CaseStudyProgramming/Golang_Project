/**
 * Root layout server load function
 * Provides authentication state to all pages
 */

import type { LayoutServerLoad } from './$types';

export const load: LayoutServerLoad = async ({ cookies }) => {
	// Check for auth token in cookies (set by client-side auth)
	const authToken = cookies.get('auth_token');
	
	return {
		isAuthenticated: !!authToken
	};
};

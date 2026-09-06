/**
 * Root layout server load function
 * Provides authentication state to all pages
 */

import type { LayoutServerLoad } from './$types';

export const load: LayoutServerLoad = async () => {
	// Note: This runs on the server, so we can't access localStorage
	// The actual auth state will be initialized on the client in +layout.svelte
	
	return {
		// You can pass server-side data here if needed
	};
};

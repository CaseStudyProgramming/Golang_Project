import tailwindcss from '@tailwindcss/vite';
import adapter from '@sveltejs/adapter-auto';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig({
	plugins: [
		tailwindcss(),
		sveltekit({
			compilerOptions: {
				// Force runes mode for the project, except for libraries. Can be removed in svelte 6.
				runes: ({ filename }) =>
					filename.split(/[/\\]/).includes('node_modules') ? undefined : true
			},

			// adapter-auto only supports some environments, see https://svelte.dev/docs/kit/adapter-auto for a list.
			// If your environment is not supported, or you settled on a specific environment, switch out the adapter.
			// See https://svelte.dev/docs/kit/adapters for more information about adapters.
			adapter: adapter()
		})
	],
	
	// Build optimizations for production
	build: {
		// Enable source maps for production debugging
		sourcemap: true,
		
		// Optimize chunk splitting
		rollupOptions: {
			output: {
				manualChunks: {
					// Vendor chunks for better caching
					'chart-vendor': ['chart.js', 'svelte-chartjs'],
					'zod': ['zod']
				}
			}
		},
		
		// Minify output
		minify: 'terser',
		terserOptions: {
			compress: {
				drop_console: true, // Remove console.log in production
				drop_debugger: true
			},
			format: {
				comments: false
			}
		}
	},
	
	// Optimize dependencies
	optimizeDeps: {
		include: ['chart.js', 'svelte-chartjs', 'zod']
	},
	
	// Server configuration
	server: {
		port: 5173,
		host: true,
		strictPort: false,
		hmr: {
			overlay: true
		}
	},
	
	// Preview configuration
	preview: {
		port: 4173,
		host: true
	}
});

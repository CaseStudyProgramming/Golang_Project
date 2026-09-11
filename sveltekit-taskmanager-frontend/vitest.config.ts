import { defineConfig } from 'vitest/config';
import { svelte } from '@sveltejs/vite-plugin-svelte';

export default defineConfig({
	plugins: [
		svelte({
			compilerOptions: {
				runes: true
			}
		})
	],
	test: {
		globals: true,
		environment: 'jsdom',
		include: ['src/**/*.{test,spec}.{js,ts}'],
		exclude: ['tests/**', 'playwright-tests/**', '**/*.e2e.ts'],
		typecheck: {
			tsconfig: './tsconfig.json'
		},
		coverage: {
			provider: 'v8',
			reporter: ['text', 'json', 'html', 'lcov'],
			exclude: [
				'node_modules/',
				'src/lib/server/',
				'tests/',
				'*.config.*',
				'vitest.config.ts',
				'svelte.config.js'
			],
			thresholds: {
				lines: 80,
				functions: 80,
				branches: 80,
				statements: 80
			}
		},
		setupFiles: ['./tests/setup.ts']
	}
});

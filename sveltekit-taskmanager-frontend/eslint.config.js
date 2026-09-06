import js from '@eslint/js';
import ts from 'typescript-eslint';
import svelte from 'eslint-plugin-svelte';
import prettier from 'eslint-config-prettier';
import perfectionist from 'eslint-plugin-perfectionist';
import globals from 'globals';

/** @type {import('eslint').Linter.Config[]} */
const config = [
	...svelte.configs['flat/recommended'],
	js.configs.recommended,
	...ts.configs.recommended,
	prettier,
	...perfectionist.configs['recommended-natural'],
	{
		languageOptions: {
			globals: {
				...globals.browser,
				...globals.node
			}
		},
		rules: {
			'@typescript-eslint/no-unused-vars': [
				'error',
				{
					argsIgnorePattern: '^_',
					varsIgnorePattern: '^_'
				}
			],
			'@typescript-eslint/no-explicit-any': 'warn',
			'perfectionist/sort-interfaces': 'off',
			'perfectionist/sort-named-properties': 'off'
		}
	},
	{
		ignores: [
			'.svelte-kit/',
			'build/',
			'dist/',
			'node_modules/',
			'*.config.js',
			'*.config.ts'
		]
	}
];

export default config;

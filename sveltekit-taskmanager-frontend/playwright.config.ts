/// <reference types="node" />
import { defineConfig, devices } from '@playwright/test';

// Detect if running in CI environment
const isCI = process.env.CI === 'true';

// Configuration for different environments
const config = {
	baseURL: process.env.BASE_URL || 'http://localhost:5173',
	backendURL: process.env.PUBLIC_API_BASE_URL || 'http://localhost:8080',
	dbHost: process.env.DB_HOST || 'localhost',
	dbPort: process.env.DB_PORT || '5433',
	dbUser: process.env.DB_USER || 'postgres',
	dbPassword: process.env.DB_PASSWORD || 'Berjuang#382',
	dbName: process.env.DB_NAME || 'taskmanager_test',
};

export default defineConfig({
	testDir: './playwright-tests',
	testMatch: '**/*.e2e.ts',
	fullyParallel: true,
	forbidOnly: !!isCI,
	retries: isCI ? 2 : 0,
	workers: isCI ? 2 : undefined,

	reporter: [
		['html'],
		['list'],
		['junit', { outputFile: 'test-results/junit.xml' }]
	],
	use: {
		baseURL: config.baseURL,
		trace: 'on-first-retry',
		screenshot: 'only-on-failure',
		video: 'retain-on-failure'
	},
	projects: [
		{
			name: 'chromium',
			use: { ...devices['Desktop Chrome'] }
		},
		{
			name: 'firefox',
			use: { ...devices['Desktop Firefox'] }
		},
		{
			name: 'webkit',
			use: { ...devices['Desktop Safari'] }
		},
		{
			name: 'Mobile Chrome',
			use: { ...devices['Pixel 5'] }
		},
		{
			name: 'Mobile Safari',
			use: { ...devices['iPhone 12'] }
		}
	],
	webServer: isCI ? undefined : [
		{
			command: 'bun run dev',
			url: 'http://localhost:5173',
			reuseExistingServer: !isCI,
			timeout: 120000,
			env: {
				PUBLIC_API_BASE_URL: config.backendURL
			}
		},
		{
			command: 'cd ../backendGoVanilaTaskmanager && go run main.go',
			url: 'http://localhost:8080',
			reuseExistingServer: !isCI,
			timeout: 120000,
			env: {
				DB_HOST: config.dbHost,
				DB_PORT: config.dbPort,
				DB_USER: config.dbUser,
				DB_PASSWORD: config.dbPassword,
				DB_NAME: config.dbName,
				SSL_MODE: 'disable',
				JWT_SECRET: 'test-secret-key-for-e2e-testing-only',
				PORT: '8080'
			}
		}
	]
});

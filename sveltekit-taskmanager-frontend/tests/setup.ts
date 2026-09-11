import { vi, afterEach, beforeAll, afterAll } from 'vitest';
import { cleanup } from '@testing-library/svelte';
import '@testing-library/jest-dom';
import { http, HttpResponse } from 'msw';
import { setupServer } from 'msw/node';

// Setup Testing Library cleanup
afterEach(() => {
	cleanup();
});

// Mock environment variables
vi.mock('$env/static/public', () => ({
	PUBLIC_API_BASE_URL: 'http://localhost:8080'
}));

vi.mock('$env/static/private', () => ({}));

// Mock localStorage
const localStorageMock = {
	getItem: vi.fn(() => null),
	removeItem: vi.fn(),
	setItem: vi.fn(),
	clear: vi.fn()
};

// Ensure localStorage is available globally
if (!globalThis.localStorage) {
	Object.defineProperty(globalThis, 'localStorage', {
		value: localStorageMock,
		writable: true
	});
}

// Mock window.location
const locationMock = {
	href: ''
};

// Ensure window.location is available globally
if (!globalThis.window?.location) {
	Object.defineProperty(globalThis, 'window', {
		value: {
			...globalThis.window,
			location: locationMock,
			document: {
				cookie: ''
			}
		},
		writable: true
	});
}

// Mock Svelte 5 runes for testing (simplified version)
// Note: Using type casting to avoid TypeScript errors with Svelte 5 rune types
(globalThis as any).$state = function <T>(initial: T): T {
	return initial;
};
(globalThis as any).$derived = function <T>(fn: () => T): T {
	return fn();
};
(globalThis as any).$props = function <T extends Record<string, unknown>>(): T {
	return {} as T;
};
(globalThis as any).$effect = function (fn: () => void): void {
	// No-op for tests
};

// Setup MSW for API mocking
const server = setupServer(
	// Example API mocks - add more as needed
	http.get('http://localhost:8080/api/tasks', () => {
		return HttpResponse.json({
			tasks: [],
			total: 0,
			page: 1,
			limit: 10
		});
	}),
	http.post('http://localhost:8080/api/auth/login', () => {
		return HttpResponse.json({
			token: 'mock-token',
			user: {
				id: '1',
				email: 'test@example.com',
				name: 'Test User'
			}
		});
	})
);

beforeAll(() => {
	server.listen({ onUnhandledRequest: 'error' });
});

afterAll(() => {
	server.close();
});

afterEach(() => {
	server.resetHandlers();
});

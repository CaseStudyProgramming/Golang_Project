import { cleanup } from '@testing-library/svelte'
import { HttpResponse, http } from 'msw'
import '@testing-library/jest-dom'
import { setupServer } from 'msw/node'
import { afterAll, afterEach, beforeAll, vi } from 'vitest'

// Setup Testing Library cleanup
afterEach(() => {
	cleanup()
})

// Mock environment variables
vi.mock('$env/static/public', () => ({
	PUBLIC_API_BASE_URL: 'http://localhost:8080',
}))

vi.mock('$env/static/private', () => ({}))

// Mock $lib/env
vi.mock('$lib/env', () => ({
	publicEnv: {
		PUBLIC_API_BASE_URL: 'http://localhost:8080',
	},
}))

// Mock localStorage
const localStorageMock = {
	clear: vi.fn(),
	getItem: vi.fn(() => null),
	removeItem: vi.fn(),
	setItem: vi.fn(),
}

// Ensure localStorage is available globally
if (!globalThis.localStorage) {
	Object.defineProperty(globalThis, 'localStorage', {
		value: localStorageMock,
		writable: true,
	})
}

// Mock window.location
const locationMock = {
	href: '',
}

// Ensure window.location is available globally
if (!globalThis.window?.location) {
	Object.defineProperty(globalThis, 'window', {
		value: {
			...globalThis.window,
			document: {
				cookie: '',
			},
			location: locationMock,
		},
		writable: true,
	})
}
// Mock Svelte 5 runes for testing (simplified version)
// Note: Using type casting to avoid TypeScript errors with Svelte 5 rune types
;(globalThis as unknown as { $state: <T>(initial: T) => T }).$state = <T>(initial: T): T => initial
;(globalThis as unknown as { $derived: <T>(fn: () => T) => T }).$derived = <T>(fn: () => T): T =>
	fn()
;(globalThis as unknown as { $props: <T extends Record<string, unknown>>() => T }).$props = <
	T extends Record<string, unknown>,
>(): T => ({}) as T
;(globalThis as unknown as { $effect: (_fn: () => void) => void }).$effect = (
	_fn: () => void
): void => {
	// No-op for tests
}

// Setup MSW for API mocking
const server = setupServer(
	// Example API mocks - add more as needed
	http.get('http://localhost:8080/api/tasks', () => {
		return HttpResponse.json({
			limit: 10,
			page: 1,
			tasks: [],
			total: 0,
		})
	}),
	http.post('http://localhost:8080/api/auth/login', () => {
		return HttpResponse.json({
			token: 'mock-token',
			user: {
				email: 'test@example.com',
				id: '1',
				name: 'Test User',
			},
		})
	})
)

beforeAll(() => {
	server.listen({ onUnhandledRequest: 'error' })
})

afterAll(() => {
	server.close()
})

afterEach(() => {
	server.resetHandlers()
})

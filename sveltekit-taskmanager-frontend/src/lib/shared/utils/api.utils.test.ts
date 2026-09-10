import { beforeEach, describe, expect, it, vi } from 'vitest';

import { HttpClient } from './api.utils';

// Mock window object for tests
const localStorageMock = {
	getItem: vi.fn(() => null),
	removeItem: vi.fn(),
	setItem: vi.fn()
};

Object.defineProperty(globalThis, 'window', {
	value: {
		...globalThis.window,
		document: {
			cookie: ''
		},
		localStorage: localStorageMock
	},
	writable: true
});

describe('HttpClient', () => {
	beforeEach(() => {
		vi.clearAllMocks();
	});

	it('handles successful API response', async () => {
		globalThis.fetch = vi.fn(() =>
			Promise.resolve({
				json: () => Promise.resolve({ data: 'test' }),
				ok: true
			} as Response)
		);

		const client = new HttpClient('http://localhost:8080');
		const result = await client.get<{ data: string }>('/test');
		expect(result).toEqual({ data: 'test' });
	});

	it('throws error on failed response', async () => {
		globalThis.fetch = vi.fn(() =>
			Promise.resolve({
				ok: false,
				status: 404,
				statusText: 'Not Found'
			} as Response)
		);

		const client = new HttpClient('http://localhost:8080');
		await expect(client.get('/test')).rejects.toThrow('API Error: 404 Not Found');
	});

	it('creates client instance with default options', () => {
		const client = new HttpClient('http://localhost:8080');
		expect(client).toBeDefined();
		expect(typeof client.get).toBe('function');
		expect(typeof client.post).toBe('function');
		expect(typeof client.put).toBe('function');
		expect(typeof client.delete).toBe('function');
	});
});

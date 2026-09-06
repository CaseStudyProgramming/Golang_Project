import { describe, it, expect, vi } from 'vitest';
import { fetchJson } from './api.utils';

describe('fetchJson', () => {
	it('handles successful API response', async () => {
		global.fetch = vi.fn(() =>
			Promise.resolve({
				ok: true,
				json: () => Promise.resolve({ data: 'test' })
			} as Response)
		);

		const result = await fetchJson<{ data: string }>('/test');
		expect(result).toEqual({ data: 'test' });
	});

	it('throws error on failed response', async () => {
		global.fetch = vi.fn(() =>
			Promise.resolve({
				ok: false,
				status: 404
			} as Response)
		);

		await expect(fetchJson('/test')).rejects.toThrow('HTTP error! status: 404');
	});
});

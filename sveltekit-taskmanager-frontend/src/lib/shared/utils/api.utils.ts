/**
 * Base API configuration and utilities
 */
import { publicEnv } from '$lib/env';

const API_BASE_URL = publicEnv.PUBLIC_API_BASE_URL;

/**
 * Performs a fetch request with error handling and JSON parsing
 */
export async function fetchJson<T>(url: string, options?: RequestInit): Promise<T> {
	const response = await fetch(`${API_BASE_URL}${url}`, {
		headers: {
			'Content-Type': 'application/json',
			...options?.headers
		},
		...options
	});

	if (!response.ok) {
		throw new Error(`HTTP error! status: ${response.status}`);
	}

	return response.json();
}

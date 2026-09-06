/**
 * Base API configuration and utilities
 */
import { publicEnv } from '$lib/env';
import type { ApiResponse, PaginatedResponse, PaginationParams } from '$lib/shared/types/api.types';
import { authRequestInterceptor, authResponseInterceptor, errorLoggingInterceptor } from './auth.interceptors';
import { ApiError } from './error.utils';

const API_BASE_URL = publicEnv.PUBLIC_API_BASE_URL;

/**
 * HTTP client configuration options
 */
interface HttpClientOptions {
	retries?: number;
	retryDelay?: number;
	timeout?: number;
	headers?: Record<string, string>;
}

/**
 * Request options for HTTP methods
 */
interface RequestOptions {
	headers?: Record<string, string>;
	signal?: AbortSignal;
	method?: string;
	body?: string;
}


/**
 * Request interceptor function type
 */
type RequestInterceptor = (request: RequestInit) => RequestInit | Promise<RequestInit>;

/**
 * Response interceptor function type
 */
type ResponseInterceptor = (response: Response) => Response | Promise<Response>;

/**
 * HTTP client class with interceptors and retry logic
 */
export class HttpClient {
	private baseUrl: string;
	private defaultOptions: HttpClientOptions;
	private requestInterceptors: RequestInterceptor[] = [];
	private responseInterceptors: ResponseInterceptor[] = [];

	constructor(baseUrl: string, options: HttpClientOptions = {}) {
		this.baseUrl = baseUrl;
		this.defaultOptions = {
			retries: 3,
			retryDelay: 1000,
			timeout: 30000,
			...options
		};
	}

	/**
	 * Add a request interceptor
	 */
	addRequestInterceptor(interceptor: RequestInterceptor): void {
		this.requestInterceptors.push(interceptor);
	}

	/**
	 * Add a response interceptor
	 */
	addResponseInterceptor(interceptor: ResponseInterceptor): void {
		this.responseInterceptors.push(interceptor);
	}

	/**
	 * Apply request interceptors
	 */
	private async applyRequestInterceptors(options: RequestInit): Promise<RequestInit> {
		let processedOptions = options;
		for (const interceptor of this.requestInterceptors) {
			processedOptions = await interceptor(processedOptions);
		}
		return processedOptions;
	}

	/**
	 * Apply response interceptors
	 */
	private async applyResponseInterceptors(response: Response): Promise<Response> {
		let processedResponse = response;
		for (const interceptor of this.responseInterceptors) {
			processedResponse = await interceptor(processedResponse);
		}
		return processedResponse;
	}

	/**
	 * Create an abort controller with timeout
	 */
	private createTimeoutController(timeout: number): AbortController {
		const controller = new AbortController();
		setTimeout(() => controller.abort(), timeout);
		return controller;
	}

	/**
	 * Sleep for retry delay
	 */
	private async sleep(ms: number): Promise<void> {
		return new Promise((resolve) => setTimeout(resolve, ms));
	}

	/**
	 * Perform HTTP request with retry logic
	 */
	private async request<T>(
		endpoint: string,
		options: RequestOptions = {}
	): Promise<T> {
		const url = `${this.baseUrl}${endpoint}`;
		const maxRetries = this.defaultOptions.retries || 0;
		let lastError: Error | null = null;

		for (let attempt = 0; attempt <= maxRetries; attempt++) {
			try {
				const timeoutController = this.createTimeoutController(
					this.defaultOptions.timeout || 30000
				);

				const processedOptions = await this.applyRequestInterceptors({
					...options,
					signal: timeoutController.signal,
					headers: {
						'Content-Type': 'application/json',
						...this.defaultOptions.headers,
						...options.headers
					}
				});

				let response = await fetch(url, processedOptions);
				response = await this.applyResponseInterceptors(response);

				if (!response.ok) {
					const errorData = await this.parseErrorData(response);
					throw new ApiError(response.status, response.statusText, errorData);
				}

				return await response.json();
			} catch (error) {
				lastError = error as Error;

				// Don't retry on abort or 4xx errors (except 429)
				if (
					error instanceof Error &&
					(error.name === 'AbortError' ||
						(error instanceof ApiError && error.status >= 400 && error.status < 500 && error.status !== 429))
				) {
					break;
				}

				// Retry with exponential backoff
				if (attempt < maxRetries) {
					const delay = (this.defaultOptions.retryDelay || 1000) * Math.pow(2, attempt);
					await this.sleep(delay);
				}
			}
		}

		throw lastError || new Error('Request failed');
	}

	/**
	 * Parse error data from response
	 */
	private async parseErrorData(response: Response): Promise<unknown> {
		try {
			return await response.json();
		} catch {
			return { message: response.statusText };
		}
	}

	/**
	 * GET request
	 */
	async get<T>(endpoint: string, options?: RequestOptions): Promise<T> {
		return this.request<T>(endpoint, { ...options, method: 'GET' });
	}

	/**
	 * POST request
	 */
	async post<T>(endpoint: string, data?: unknown, options?: RequestOptions): Promise<T> {
		return this.request<T>(endpoint, {
			...options,
			method: 'POST',
			body: JSON.stringify(data)
		});
	}

	/**
	 * PUT request
	 */
	async put<T>(endpoint: string, data?: unknown, options?: RequestOptions): Promise<T> {
		return this.request<T>(endpoint, {
			...options,
			method: 'PUT',
			body: JSON.stringify(data)
		});
	}

	/**
	 * PATCH request
	 */
	async patch<T>(endpoint: string, data?: unknown, options?: RequestOptions): Promise<T> {
		return this.request<T>(endpoint, {
			...options,
			method: 'PATCH',
			body: JSON.stringify(data)
		});
	}

	/**
	 * DELETE request
	 */
	async delete<T>(endpoint: string, options?: RequestOptions): Promise<T> {
		return this.request<T>(endpoint, { ...options, method: 'DELETE' });
	}
}

/**
 * Create HTTP client instance with interceptors
 * Only create client instance on client side to avoid SSR issues
 */
let httpClientInstance: HttpClient | null = null;

export const httpClient = new Proxy({} as HttpClient, {
	get(_target, prop) {
		if (!httpClientInstance) {
			httpClientInstance = new HttpClient(API_BASE_URL);
			// Add authentication interceptors
			httpClientInstance.addRequestInterceptor(authRequestInterceptor);
			httpClientInstance.addResponseInterceptor(authResponseInterceptor);
			httpClientInstance.addResponseInterceptor(errorLoggingInterceptor);
		}
		return httpClientInstance[prop as keyof HttpClient];
	}
});

/**
 * Legacy fetchJson function for backward compatibility
 * @deprecated Use httpClient methods instead
 */
export async function fetchJson<T>(url: string, options?: RequestOptions): Promise<T> {
	const client = httpClient as HttpClient;
	return client.get<T>(url, options);
}

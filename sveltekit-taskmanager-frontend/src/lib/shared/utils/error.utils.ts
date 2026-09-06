/**
 * Error handling utilities and wrappers
 */

/**
 * Custom error class for API errors
 */
export class ApiError extends Error {
	constructor(
		public status: number,
		public statusText: string,
		public data?: unknown
	) {
		super(`API Error: ${status} ${statusText}`);
		this.name = 'ApiError';
	}
}

/**
 * Custom error class for validation errors
 */
export class ValidationError extends Error {
	constructor(public field: string, message: string) {
		super(message);
		this.name = 'ValidationError';
	}
}

/**
 * Custom error class for authentication errors
 */
export class AuthenticationError extends Error {
	constructor(message: string = 'Authentication failed') {
		super(message);
		this.name = 'AuthenticationError';
	}
}

/**
 * Error handler function type
 */
type ErrorHandler = (error: Error) => void;

/**
 * Global error handlers
 */
const errorHandlers: ErrorHandler[] = [];

/**
 * Register global error handler
 */
export function registerErrorHandler(handler: ErrorHandler): void {
	errorHandlers.push(handler);
}

/**
 * Execute all registered error handlers
 */
function executeErrorHandlers(error: Error): void {
	errorHandlers.forEach((handler) => {
		try {
			handler(error);
		} catch (handlerError) {
			console.error('Error in error handler:', handlerError);
		}
	});
}

/**
 * Wrap async function with error handling
 */
export async function withErrorHandling<T>(
	fn: () => Promise<T>,
	context?: string
): Promise<T> {
	try {
		return await fn();
	} catch (error) {
		const processedError = error instanceof Error ? error : new Error(String(error));
		
		if (context) {
			processedError.message = `${context}: ${processedError.message}`;
		}
		
		executeErrorHandlers(processedError);
		throw processedError;
	}
}

/**
 * Check if error is a specific type
 */
export function isApiError(error: unknown): error is ApiError {
	return error instanceof ApiError;
}

export function isValidationError(error: unknown): error is ValidationError {
	return error instanceof ValidationError;
}

export function isAuthenticationError(error: unknown): error is AuthenticationError {
	return error instanceof AuthenticationError;
}

/**
 * Get user-friendly error message
 */
export function getErrorMessage(error: unknown): string {
	if (isApiError(error)) {
		switch (error.status) {
			case 400:
				return 'Invalid request. Please check your input.';
			case 401:
				return 'You need to log in to access this resource.';
			case 403:
				return 'You don\'t have permission to access this resource.';
			case 404:
				return 'The requested resource was not found.';
			case 429:
				return 'Too many requests. Please try again later.';
			case 500:
				return 'Server error. Please try again later.';
			default:
				return error.statusText || 'An error occurred.';
		}
	}

	if (isValidationError(error)) {
		return `Validation error: ${error.message}`;
	}

	if (isAuthenticationError(error)) {
		return error.message;
	}

	if (error instanceof Error) {
		return error.message;
	}

	return 'An unexpected error occurred.';
}

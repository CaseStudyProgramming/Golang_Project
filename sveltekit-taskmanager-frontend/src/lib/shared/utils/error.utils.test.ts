import { beforeEach, describe, expect, it, vi } from 'vitest';

import {
	ApiError,
	AuthenticationError,
	clearErrorHandlers,
	getErrorMessage,
	isApiError,
	isAuthenticationError,
	isValidationError,
	registerErrorHandler,
	ValidationError,
	withErrorHandling
} from './error.utils';

describe('Error Utils', () => {
	beforeEach(() => {
		vi.clearAllMocks();
		clearErrorHandlers();
	});

	describe('Custom Error Classes', () => {
		it('creates ApiError with correct properties', () => {
			const error = new ApiError(404, 'Not Found', { details: 'Resource missing' });
			
			expect(error).toBeInstanceOf(Error);
			expect(error.name).toBe('ApiError');
			expect(error.status).toBe(404);
			expect(error.statusText).toBe('Not Found');
			expect(error.data).toEqual({ details: 'Resource missing' });
			expect(error.message).toBe('API Error: 404 Not Found');
		});

		it('creates AuthenticationError with default message', () => {
			const error = new AuthenticationError();
			
			expect(error).toBeInstanceOf(Error);
			expect(error.name).toBe('AuthenticationError');
			expect(error.message).toBe('Authentication failed');
		});

		it('creates AuthenticationError with custom message', () => {
			const error = new AuthenticationError('Token expired');
			
			expect(error).toBeInstanceOf(Error);
			expect(error.name).toBe('AuthenticationError');
			expect(error.message).toBe('Token expired');
		});

		it('creates ValidationError with field and message', () => {
			const error = new ValidationError('email', 'Invalid email format');
			
			expect(error).toBeInstanceOf(Error);
			expect(error.name).toBe('ValidationError');
			expect(error.field).toBe('email');
			expect(error.message).toBe('Invalid email format');
		});
	});

	describe('Error Type Guards', () => {
		it('identifies ApiError correctly', () => {
			const apiError = new ApiError(400, 'Bad Request');
			const authError = new AuthenticationError();
			const validationError = new ValidationError('field', 'error');
			const genericError = new Error('Generic error');

			expect(isApiError(apiError)).toBe(true);
			expect(isApiError(authError)).toBe(false);
			expect(isApiError(validationError)).toBe(false);
			expect(isApiError(genericError)).toBe(false);
			expect(isApiError(null)).toBe(false);
			expect(isApiError(undefined)).toBe(false);
		});

		it('identifies AuthenticationError correctly', () => {
			const authError = new AuthenticationError();
			const apiError = new ApiError(401, 'Unauthorized');
			const genericError = new Error('Generic error');

			expect(isAuthenticationError(authError)).toBe(true);
			expect(isAuthenticationError(apiError)).toBe(false);
			expect(isAuthenticationError(genericError)).toBe(false);
			expect(isAuthenticationError(null)).toBe(false);
		});

		it('identifies ValidationError correctly', () => {
			const validationError = new ValidationError('field', 'error');
			const apiError = new ApiError(400, 'Bad Request');
			const genericError = new Error('Generic error');

			expect(isValidationError(validationError)).toBe(true);
			expect(isValidationError(apiError)).toBe(false);
			expect(isValidationError(genericError)).toBe(false);
			expect(isValidationError(null)).toBe(false);
		});
	});

	describe('getErrorMessage', () => {
		it('returns user-friendly message for 400 error', () => {
			const error = new ApiError(400, 'Bad Request');
			expect(getErrorMessage(error)).toBe('Invalid request. Please check your input.');
		});

		it('returns user-friendly message for 401 error', () => {
			const error = new ApiError(401, 'Unauthorized');
			expect(getErrorMessage(error)).toBe('You need to log in to access this resource.');
		});

		it('returns user-friendly message for 403 error', () => {
			const error = new ApiError(403, 'Forbidden');
			expect(getErrorMessage(error)).toBe("You don't have permission to access this resource.");
		});

		it('returns user-friendly message for 404 error', () => {
			const error = new ApiError(404, 'Not Found');
			expect(getErrorMessage(error)).toBe('The requested resource was not found.');
		});

		it('returns user-friendly message for 429 error', () => {
			const error = new ApiError(429, 'Too Many Requests');
			expect(getErrorMessage(error)).toBe('Too many requests. Please try again later.');
		});

		it('returns user-friendly message for 500 error', () => {
			const error = new ApiError(500, 'Internal Server Error');
			expect(getErrorMessage(error)).toBe('Server error. Please try again later.');
		});

		it('returns statusText for unknown API error status', () => {
			const error = new ApiError(418, "I'm a teapot");
			expect(getErrorMessage(error)).toBe("I'm a teapot");
		});

		it('returns validation error message', () => {
			const error = new ValidationError('email', 'Invalid format');
			expect(getErrorMessage(error)).toBe('Validation error: Invalid format');
		});

		it('returns authentication error message', () => {
			const error = new AuthenticationError('Session expired');
			expect(getErrorMessage(error)).toBe('Session expired');
		});

		it('returns generic error message for standard Error', () => {
			const error = new Error('Something went wrong');
			expect(getErrorMessage(error)).toBe('Something went wrong');
		});

		it('returns default message for non-error values', () => {
			expect(getErrorMessage(null)).toBe('An unexpected error occurred.');
			expect(getErrorMessage(undefined)).toBe('An unexpected error occurred.');
			expect(getErrorMessage('string error')).toBe('An unexpected error occurred.');
		});
	});

	describe('Error Handlers', () => {
		it('registers and executes error handler', () => {
			const handler = vi.fn();
			registerErrorHandler(handler);

			const error = new Error('Test error');
			withErrorHandling(() => {
				throw error;
			}).catch(() => {});

			expect(handler).toHaveBeenCalledWith(error);
		});

		it('executes multiple registered error handlers', () => {
			const handler1 = vi.fn();
			const handler2 = vi.fn();
			registerErrorHandler(handler1);
			registerErrorHandler(handler2);

			const error = new Error('Test error');
			withErrorHandling(() => {
				throw error;
			}).catch(() => {});

			expect(handler1).toHaveBeenCalledWith(error);
			expect(handler2).toHaveBeenCalledWith(error);
		});

		it('handles errors in error handlers gracefully', () => {
			const badHandler = vi.fn(() => {
				throw new Error('Handler error');
			});
			const goodHandler = vi.fn();
			registerErrorHandler(badHandler);
			registerErrorHandler(goodHandler);

			const error = new Error('Test error');
			withErrorHandling(() => {
				throw error;
			}).catch(() => {});

			expect(badHandler).toHaveBeenCalled();
			expect(goodHandler).toHaveBeenCalled();
		});

		it('does not throw when error handler fails', () => {
			const badHandler = vi.fn(() => {
				throw new Error('Handler error');
			});
			registerErrorHandler(badHandler);

			expect(() => {
				withErrorHandling(() => {
					throw new Error('Test error');
				}).catch(() => {});
			}).not.toThrow();
		});
	});

	describe('withErrorHandling', () => {
		it('wraps successful function execution', async () => {
			const fn = vi.fn(() => Promise.resolve('success'));
			const result = await withErrorHandling(fn);

			expect(result).toBe('success');
			expect(fn).toHaveBeenCalled();
		});

		it('wraps async function execution', async () => {
			const fn = vi.fn(async () => 'async success');
			const result = await withErrorHandling(fn);

			expect(result).toBe('async success');
			expect(fn).toHaveBeenCalled();
		});

		it('throws error when function fails', async () => {
			const fn = vi.fn(() => {
				throw new Error('Function error');
			});

			await expect(withErrorHandling(fn)).rejects.toThrow('Function error');
		});

		it('adds context to error message when provided', async () => {
			const fn = vi.fn(() => {
				throw new Error('Function error');
			});

			try {
				await withErrorHandling(fn, 'User operation');
			} catch (error) {
				expect(error).toBeInstanceOf(Error);
				expect((error as Error).message).toBe('User operation: Function error');
			}
		});

		it('converts non-error values to Error', async () => {
			const fn = vi.fn(() => {
				throw 'string error';
			});

			try {
				await withErrorHandling(fn);
			} catch (error) {
				expect(error).toBeInstanceOf(Error);
				expect((error as Error).message).toBe('string error');
			}
		});

		it('executes error handlers before throwing', async () => {
			const handler = vi.fn();
			registerErrorHandler(handler);

			const fn = vi.fn(() => {
				throw new Error('Test error');
			});

			try {
				await withErrorHandling(fn);
			} catch {
				expect(handler).toHaveBeenCalled();
			}
		});
	});
});
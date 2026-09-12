import { beforeEach, describe, expect, it, vi } from 'vitest'
import type { Mock } from 'vitest'

import {
	authRequestInterceptor,
	authResponseInterceptor,
	errorLoggingInterceptor,
	getAuthToken,
	getRefreshToken,
	isTokenExpired,
	removeAuthToken,
	setAuthToken,
	setRefreshToken,
} from './auth.interceptors'

describe('Auth Interceptors', () => {
	beforeEach(() => {
		vi.clearAllMocks()
		// Reset localStorage mock
		const localStorageMock = {
			getItem: vi.fn(() => null) as Mock,
			setItem: vi.fn() as Mock,
			removeItem: vi.fn() as Mock,
			clear: vi.fn() as Mock,
		}
		Object.defineProperty(globalThis, 'localStorage', {
			value: localStorageMock,
			writable: true,
		})
		Object.defineProperty(globalThis, 'window', {
			value: {},
			writable: true,
		})
		Object.defineProperty(globalThis, 'document', {
			value: { cookie: '' },
			writable: true,
		})
	})

	describe('authResponseInterceptor', () => {
		it('removes token on 401 response', async () => {
			const windowMock = {
				location: { href: '' },
				document: { cookie: '' },
			}
			Object.defineProperty(globalThis, 'window', {
				value: windowMock,
				writable: true,
			})

			const response = new Response(null, { status: 401, statusText: 'Unauthorized' })
			const result = await authResponseInterceptor(response)

			expect(result).toBe(response)
			expect(globalThis.localStorage.removeItem as Mock).toHaveBeenCalledWith('auth_token')
			expect(globalThis.localStorage.removeItem as Mock).toHaveBeenCalledWith('refresh_token')
			expect(globalThis.localStorage.removeItem as Mock).toHaveBeenCalledWith('token_expiry')
		})

		it('handles 403 Forbidden response', async () => {
			const windowMock = {
				location: { href: '' },
				document: { cookie: '' },
			}
			Object.defineProperty(globalThis, 'window', {
				value: windowMock,
				writable: true,
			})

			const response = new Response(null, { status: 403, statusText: 'Forbidden' })
			const result = await authResponseInterceptor(response)

			expect(result).toBe(response)
		})

		it('passes through successful responses', async () => {
			const response = new Response(null, { status: 200, statusText: 'OK' })
			const result = await authResponseInterceptor(response)

			expect(result).toBe(response)
		})
	})

	describe('errorLoggingInterceptor', () => {
		it('logs error responses', () => {
			const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {})

			const response = new Response('{"error": "test"}', {
				status: 500,
				statusText: 'Internal Server Error',
			})
			Object.defineProperty(response, 'url', { value: 'http://test.com/api', writable: true })

			const result = errorLoggingInterceptor(response)

			expect(consoleSpy).toHaveBeenCalledWith(
				'API Error: 500 Internal Server Error',
				expect.objectContaining({
					status: 500,
					statusText: 'Internal Server Error',
					url: 'http://test.com/api',
				})
			)
			expect(result).toBe(response)

			consoleSpy.mockRestore()
		})

		it('does not log successful responses', () => {
			const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {})

			const response = new Response(null, { status: 200, statusText: 'OK' })
			const result = errorLoggingInterceptor(response)

			expect(consoleSpy).not.toHaveBeenCalled()
			expect(result).toBe(response)

			consoleSpy.mockRestore()
		})
	})

	describe('authRequestInterceptor', () => {
		it('adds Authorization header when token exists', () => {
			;(globalThis.localStorage.getItem as Mock).mockImplementation((key: string) => {
				if (key === 'auth_token') return 'test-token'
				if (key === 'token_expiry') return (Date.now() + 3600000).toString()
				return null
			})

			const request: RequestInit = {
				headers: { 'Content-Type': 'application/json' },
			}
			const result = authRequestInterceptor(request)

			expect(result.headers).toEqual({
				'Content-Type': 'application/json',
				Authorization: 'Bearer test-token',
			})
		})

		it('does not add Authorization header when no token', () => {
			;(globalThis.localStorage.getItem as Mock).mockReturnValue(null)

			const request: RequestInit = {
				headers: { 'Content-Type': 'application/json' },
			}
			const result = authRequestInterceptor(request)

			expect(result.headers).toEqual({
				'Content-Type': 'application/json',
			})
		})

		it('handles request without existing headers', () => {
			;(globalThis.localStorage.getItem as Mock).mockImplementation((key: string) => {
				if (key === 'auth_token') return 'test-token'
				if (key === 'token_expiry') return (Date.now() + 3600000).toString()
				return null
			})

			const request: RequestInit = {}
			const result = authRequestInterceptor(request)

			expect(result.headers).toEqual({
				Authorization: 'Bearer test-token',
			})
		})
	})

	describe('getAuthToken', () => {
		it('returns null when window is undefined', () => {
			Object.defineProperty(globalThis, 'window', {
				value: undefined,
				writable: true,
			})

			const token = getAuthToken()
			expect(token).toBeNull()
		})

		it('returns null when no token in localStorage', () => {
			;(globalThis.localStorage.getItem as Mock).mockReturnValue(null)

			const token = getAuthToken()
			expect(token).toBeNull()
		})

		it('returns token when valid token exists', () => {
			;(globalThis.localStorage.getItem as Mock).mockImplementation((key: string) => {
				if (key === 'auth_token') return 'valid-token'
				if (key === 'token_expiry') return (Date.now() + 3600000).toString()
				return null
			})

			const token = getAuthToken()
			expect(token).toBe('valid-token')
		})

		it('returns null and removes token when expired', () => {
			;(globalThis.localStorage.getItem as Mock).mockImplementation((key: string) => {
				if (key === 'auth_token') return 'expired-token'
				if (key === 'token_expiry') return (Date.now() - 1000).toString()
				return null
			})

			const token = getAuthToken()
			expect(token).toBeNull()
			expect(globalThis.localStorage.removeItem as Mock).toHaveBeenCalledWith('auth_token')
			expect(globalThis.localStorage.removeItem as Mock).toHaveBeenCalledWith('refresh_token')
			expect(globalThis.localStorage.removeItem as Mock).toHaveBeenCalledWith('token_expiry')
		})
	})

	describe('getRefreshToken', () => {
		it('returns null when window is undefined', () => {
			Object.defineProperty(globalThis, 'window', {
				value: undefined,
				writable: true,
			})

			const token = getRefreshToken()
			expect(token).toBeNull()
		})

		it('returns refresh token when exists', () => {
			;(globalThis.localStorage.getItem as Mock).mockImplementation((key: string) => {
				if (key === 'refresh_token') return 'refresh-token'
				return null
			})

			const token = getRefreshToken()
			expect(token).toBe('refresh-token')
		})

		it('returns null when no refresh token', () => {
			;(globalThis.localStorage.getItem as Mock).mockReturnValue(null)

			const token = getRefreshToken()
			expect(token).toBeNull()
		})
	})

	describe('isTokenExpired', () => {
		it('returns true when no expiry in localStorage', () => {
			;(globalThis.localStorage.getItem as Mock).mockReturnValue(null)

			const expired = isTokenExpired()
			expect(expired).toBe(true)
		})

		it('returns true when token is expired', () => {
			;(globalThis.localStorage.getItem as Mock).mockReturnValue((Date.now() - 1000).toString())

			const expired = isTokenExpired()
			expect(expired).toBe(true)
		})

		it('returns false when token is not expired', () => {
			;(globalThis.localStorage.getItem as Mock).mockReturnValue((Date.now() + 3600000).toString())

			const expired = isTokenExpired()
			expect(expired).toBe(false)
		})
	})

	describe('removeAuthToken', () => {
		it('does nothing when window is undefined', () => {
			Object.defineProperty(globalThis, 'window', {
				value: undefined,
				writable: true,
			})

			removeAuthToken()
			// Should not throw
		})

		it('removes all auth tokens from localStorage', () => {
			removeAuthToken()

			expect(globalThis.localStorage.removeItem).toHaveBeenCalledWith('auth_token')
			expect(globalThis.localStorage.removeItem).toHaveBeenCalledWith('refresh_token')
			expect(globalThis.localStorage.removeItem).toHaveBeenCalledWith('token_expiry')
		})

		it('clears auth token cookie', () => {
			removeAuthToken()

			expect(globalThis.document.cookie).toBe('auth_token=; path=/; max-age=0; SameSite=Strict')
		})
	})

	describe('setAuthToken', () => {
		it('does nothing when window is undefined', () => {
			Object.defineProperty(globalThis, 'window', {
				value: undefined,
				writable: true,
			})

			setAuthToken('test-token')
			// Should not throw
		})

		it('sets token in localStorage with default expiry', () => {
			setAuthToken('test-token')

			expect(globalThis.localStorage.setItem).toHaveBeenCalledWith('auth_token', 'test-token')
			expect(globalThis.localStorage.setItem).toHaveBeenCalledWith(
				'token_expiry',
				expect.any(String)
			)
			expect(globalThis.document.cookie).toContain('auth_token=test-token')
		})

		it('sets token in localStorage with custom expiry', () => {
			setAuthToken('test-token', 7200)

			expect(globalThis.localStorage.setItem).toHaveBeenCalledWith('auth_token', 'test-token')
			expect(globalThis.localStorage.setItem).toHaveBeenCalledWith(
				'token_expiry',
				expect.any(String)
			)
			expect(globalThis.document.cookie).toContain('max-age=7200')
		})
	})

	describe('setRefreshToken', () => {
		it('does nothing when window is undefined', () => {
			Object.defineProperty(globalThis, 'window', {
				value: undefined,
				writable: true,
			})

			setRefreshToken('test-refresh-token')
			// Should not throw
		})

		it('sets refresh token in localStorage', () => {
			setRefreshToken('test-refresh-token')

			expect(globalThis.localStorage.setItem).toHaveBeenCalledWith(
				'refresh_token',
				'test-refresh-token'
			)
		})
	})

	describe('authRequestInterceptor', () => {
		it('adds Authorization header when token exists', () => {
			const localStorageMock = {
				getItem: vi.fn((key) => {
					if (key === 'auth_token') return 'test-token'
					if (key === 'token_expiry') return (Date.now() + 3600000).toString()
					return null
				}),
			}

			Object.defineProperty(globalThis, 'localStorage', {
				value: localStorageMock,
				writable: true,
			})

			const request: RequestInit = {
				headers: { 'Content-Type': 'application/json' },
			}
			const result = authRequestInterceptor(request)

			expect(result.headers).toEqual({
				'Content-Type': 'application/json',
				Authorization: 'Bearer test-token',
			})
		})

		it('does not add Authorization header when no token', () => {
			const localStorageMock = {
				getItem: vi.fn(() => null),
			}

			Object.defineProperty(globalThis, 'localStorage', {
				value: localStorageMock,
				writable: true,
			})

			const request: RequestInit = {
				headers: { 'Content-Type': 'application/json' },
			}
			const result = authRequestInterceptor(request)

			expect(result.headers).toEqual({
				'Content-Type': 'application/json',
			})
		})

		it('handles request without existing headers', () => {
			const localStorageMock = {
				getItem: vi.fn((key) => {
					if (key === 'auth_token') return 'test-token'
					if (key === 'token_expiry') return (Date.now() + 3600000).toString()
					return null
				}),
			}

			Object.defineProperty(globalThis, 'localStorage', {
				value: localStorageMock,
				writable: true,
			})

			const request: RequestInit = {}
			const result = authRequestInterceptor(request)

			expect(result.headers).toEqual({
				Authorization: 'Bearer test-token',
			})
		})
	})

	describe('getAuthToken', () => {
		it('returns null when window is undefined', () => {
			Object.defineProperty(globalThis, 'window', {
				value: undefined,
				writable: true,
			})

			const token = getAuthToken()
			expect(token).toBeNull()
		})

		it('returns null when no token in localStorage', () => {
			const localStorageMock = {
				getItem: vi.fn(() => null),
			}

			Object.defineProperty(globalThis, 'localStorage', {
				value: localStorageMock,
				writable: true,
			})
			Object.defineProperty(globalThis, 'window', {
				value: {},
				writable: true,
			})

			const token = getAuthToken()
			expect(token).toBeNull()
		})

		it('returns token when valid token exists', () => {
			const localStorageMock = {
				getItem: vi.fn((key) => {
					if (key === 'auth_token') return 'valid-token'
					if (key === 'token_expiry') return (Date.now() + 3600000).toString()
					return null
				}),
			}

			Object.defineProperty(globalThis, 'localStorage', {
				value: localStorageMock,
				writable: true,
			})
			Object.defineProperty(globalThis, 'window', {
				value: {},
				writable: true,
			})

			const token = getAuthToken()
			expect(token).toBe('valid-token')
		})

		it('returns null and removes token when expired', () => {
			const removeSpy = vi.fn()
			const localStorageMock = {
				getItem: vi.fn((key) => {
					if (key === 'auth_token') return 'expired-token'
					if (key === 'token_expiry') return (Date.now() - 1000).toString()
					return null
				}),
				removeItem: removeSpy,
			}

			Object.defineProperty(globalThis, 'localStorage', {
				value: localStorageMock,
				writable: true,
			})
			Object.defineProperty(globalThis, 'window', {
				value: {},
				writable: true,
			})

			const token = getAuthToken()
			expect(token).toBeNull()
			expect(removeSpy).toHaveBeenCalledWith('auth_token')
			expect(removeSpy).toHaveBeenCalledWith('refresh_token')
			expect(removeSpy).toHaveBeenCalledWith('token_expiry')
		})
	})

	describe('getRefreshToken', () => {
		it('returns null when window is undefined', () => {
			Object.defineProperty(globalThis, 'window', {
				value: undefined,
				writable: true,
			})

			const token = getRefreshToken()
			expect(token).toBeNull()
		})

		it('returns refresh token when exists', () => {
			const localStorageMock = {
				getItem: vi.fn((key) => {
					if (key === 'refresh_token') return 'refresh-token'
					return null
				}),
			}

			Object.defineProperty(globalThis, 'localStorage', {
				value: localStorageMock,
				writable: true,
			})
			Object.defineProperty(globalThis, 'window', {
				value: {},
				writable: true,
			})

			const token = getRefreshToken()
			expect(token).toBe('refresh-token')
		})

		it('returns null when no refresh token', () => {
			const localStorageMock = {
				getItem: vi.fn(() => null),
			}

			Object.defineProperty(globalThis, 'localStorage', {
				value: localStorageMock,
				writable: true,
			})
			Object.defineProperty(globalThis, 'window', {
				value: {},
				writable: true,
			})

			const token = getRefreshToken()
			expect(token).toBeNull()
		})
	})

	describe('isTokenExpired', () => {
		it('returns true when no expiry in localStorage', () => {
			const localStorageMock = {
				getItem: vi.fn(() => null),
			}

			Object.defineProperty(globalThis, 'localStorage', {
				value: localStorageMock,
				writable: true,
			})

			const expired = isTokenExpired()
			expect(expired).toBe(true)
		})

		it('returns true when token is expired', () => {
			const localStorageMock = {
				getItem: vi.fn(() => (Date.now() - 1000).toString()),
			}

			Object.defineProperty(globalThis, 'localStorage', {
				value: localStorageMock,
				writable: true,
			})

			const expired = isTokenExpired()
			expect(expired).toBe(true)
		})

		it('returns false when token is not expired', () => {
			const localStorageMock = {
				getItem: vi.fn(() => (Date.now() + 3600000).toString()),
			}

			Object.defineProperty(globalThis, 'localStorage', {
				value: localStorageMock,
				writable: true,
			})

			const expired = isTokenExpired()
			expect(expired).toBe(false)
		})
	})

	describe('removeAuthToken', () => {
		it('does nothing when window is undefined', () => {
			Object.defineProperty(globalThis, 'window', {
				value: undefined,
				writable: true,
			})

			removeAuthToken()
			// Should not throw
		})

		it('removes all auth tokens from localStorage', () => {
			const removeSpy = vi.fn()
			const localStorageMock = {
				removeItem: removeSpy,
			}

			Object.defineProperty(globalThis, 'localStorage', {
				value: localStorageMock,
				writable: true,
			})
			Object.defineProperty(globalThis, 'window', {
				value: {},
				writable: true,
			})

			removeAuthToken()

			expect(removeSpy).toHaveBeenCalledWith('auth_token')
			expect(removeSpy).toHaveBeenCalledWith('refresh_token')
			expect(removeSpy).toHaveBeenCalledWith('token_expiry')
		})

		it('clears auth token cookie', () => {
			const removeSpy = vi.fn()
			const localStorageMock = {
				removeItem: removeSpy,
			}
			const documentMock = {
				cookie: '',
			}

			Object.defineProperty(globalThis, 'localStorage', {
				value: localStorageMock,
				writable: true,
			})
			Object.defineProperty(globalThis, 'window', {
				value: { document: documentMock },
				writable: true,
			})
			Object.defineProperty(globalThis, 'document', {
				value: documentMock,
				writable: true,
			})

			removeAuthToken()

			expect(documentMock.cookie).toBe('auth_token=; path=/; max-age=0; SameSite=Strict')
		})
	})

	describe('setAuthToken', () => {
		it('does nothing when window is undefined', () => {
			Object.defineProperty(globalThis, 'window', {
				value: undefined,
				writable: true,
			})

			setAuthToken('test-token')
			// Should not throw
		})

		it('sets token in localStorage with default expiry', () => {
			const setSpy = vi.fn()
			const localStorageMock = {
				setItem: setSpy,
			}
			const documentMock = {
				cookie: '',
			}

			Object.defineProperty(globalThis, 'localStorage', {
				value: localStorageMock,
				writable: true,
			})
			Object.defineProperty(globalThis, 'window', {
				value: { document: documentMock },
				writable: true,
			})
			Object.defineProperty(globalThis, 'document', {
				value: documentMock,
				writable: true,
			})

			setAuthToken('test-token')

			expect(setSpy).toHaveBeenCalledWith('auth_token', 'test-token')
			expect(setSpy).toHaveBeenCalledWith('token_expiry', expect.any(String))
			expect(documentMock.cookie).toContain('auth_token=test-token')
		})

		it('sets token in localStorage with custom expiry', () => {
			const setSpy = vi.fn()
			const localStorageMock = {
				setItem: setSpy,
			}
			const documentMock = {
				cookie: '',
			}

			Object.defineProperty(globalThis, 'localStorage', {
				value: localStorageMock,
				writable: true,
			})
			Object.defineProperty(globalThis, 'window', {
				value: { document: documentMock },
				writable: true,
			})
			Object.defineProperty(globalThis, 'document', {
				value: documentMock,
				writable: true,
			})

			setAuthToken('test-token', 7200)

			expect(setSpy).toHaveBeenCalledWith('auth_token', 'test-token')
			expect(setSpy).toHaveBeenCalledWith('token_expiry', expect.any(String))
			expect(documentMock.cookie).toContain('max-age=7200')
		})
	})

	describe('setRefreshToken', () => {
		it('does nothing when window is undefined', () => {
			Object.defineProperty(globalThis, 'window', {
				value: undefined,
				writable: true,
			})

			setRefreshToken('test-refresh-token')
			// Should not throw
		})

		it('sets refresh token in localStorage', () => {
			const setSpy = vi.fn()
			const localStorageMock = {
				setItem: setSpy,
			}

			Object.defineProperty(globalThis, 'localStorage', {
				value: localStorageMock,
				writable: true,
			})
			Object.defineProperty(globalThis, 'window', {
				value: {},
				writable: true,
			})

			setRefreshToken('test-refresh-token')

			expect(setSpy).toHaveBeenCalledWith('refresh_token', 'test-refresh-token')
		})
	})
})

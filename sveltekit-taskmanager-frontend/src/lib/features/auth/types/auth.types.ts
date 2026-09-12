/**
 * Authentication feature types
 */

/**
 * Authentication state interface
 */
export interface AuthState {
	isAuthenticated: boolean
	user: null | User
	token: null | string
	isLoading: boolean
	error: null | string
}

/**
 * User interface
 */
export interface User {
	id: string
	email: string
	name?: string
	role?: string
	createdAt?: string
	updatedAt?: string
}

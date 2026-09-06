/**
 * Authentication feature types
 */

/**
 * User interface
 */
export interface User {
	id: string;
	email: string;
	name?: string;
	role?: string;
	createdAt?: string;
	updatedAt?: string;
}

/**
 * Authentication state interface
 */
export interface AuthState {
	isAuthenticated: boolean;
	user: User | null;
	token: string | null;
	isLoading: boolean;
	error: string | null;
}

/**
 * Login credentials
 */
export interface LoginCredentials {
	email: string;
	password: string;
}

/**
 * Registration data
 */
export interface RegistrationData {
	email: string;
	password: string;
	name?: string;
}

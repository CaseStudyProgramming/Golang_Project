/**
 * Authentication schemas with Zod validation for OWASP compliance
 */

import { z } from 'zod';

/**
 * Email validation schema
 */
const emailSchema = z
	.string()
	.min(1, 'Email is required')
	.email('Invalid email format')
	.max(255, 'Email is too long')
	.transform((val) => val.toLowerCase().trim());

/**
 * Password validation schema (OWASP password guidelines)
 */
const passwordSchema = z
	.string()
	.min(8, 'Password must be at least 8 characters')
	.max(128, 'Password is too long')
	.regex(/[A-Z]/, 'Password must contain at least one uppercase letter')
	.regex(/[a-z]/, 'Password must contain at least one lowercase letter')
	.regex(/[0-9]/, 'Password must contain at least one number')
	.regex(/[^A-Za-z0-9]/, 'Password must contain at least one special character');

/**
 * Name validation schema
 */
const nameSchema = z
	.string()
	.min(1, 'Name is required')
	.max(100, 'Name is too long')
	.transform((val) => val.trim());

/**
 * Login credentials validation schema
 */
export const loginSchema = z.object({
	email: emailSchema,
	password: z.string().min(1, 'Password is required')
});

/**
 * Registration data validation schema
 */
export const registerSchema = z.object({
	email: emailSchema,
	password: passwordSchema,
	name: nameSchema.optional()
});

/**
 * Forgot password request validation schema
 */
export const forgotPasswordSchema = z.object({
	email: emailSchema
});

/**
 * Reset password request validation schema
 */
export const resetPasswordSchema = z.object({
	token: z.string().min(1, 'Reset token is required'),
	password: passwordSchema,
	confirmPassword: z.string().min(1, 'Please confirm your password')
}).refine((data) => data.password === data.confirmPassword, {
	message: 'Passwords do not match',
	path: ['confirmPassword']
});

/**
 * Type inference for login credentials
 */
export type LoginCredentials = z.infer<typeof loginSchema>;

/**
 * Type inference for registration data
 */
export type RegistrationData = z.infer<typeof registerSchema>;

/**
 * Type inference for forgot password request
 */
export type ForgotPasswordRequest = z.infer<typeof forgotPasswordSchema>;

/**
 * Type inference for reset password request
 */
export type ResetPasswordRequest = z.infer<typeof resetPasswordSchema>;

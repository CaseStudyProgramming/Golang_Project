import { z } from 'zod';

/**
 * Public environment variables schema
 * All public variables must start with PUBLIC_
 */
const publicEnvSchema = z.object({
	PUBLIC_API_BASE_URL: z.string().url().default('http://localhost:8080')
});

/**
 * Validate and export public environment variables
 */
export const publicEnv = publicEnvSchema.parse({
	PUBLIC_API_BASE_URL: import.meta.env.PUBLIC_API_BASE_URL || 'http://localhost:8080'
});

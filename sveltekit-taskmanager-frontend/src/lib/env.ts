import { env } from '$env/static/public';
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
export const publicEnv = publicEnvSchema.parse(env);

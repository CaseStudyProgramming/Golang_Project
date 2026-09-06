import { env } from '$env/static/private';
import { z } from 'zod';

/**
 * Server-only environment variables schema
 * These are never exposed to the client
 */
const serverEnvSchema = z.object({
	// Add server-side environment variables here
	// Example: DATABASE_URL: z.string().url()
});

/**
 * Validate and export server environment variables
 * This should be imported at the top of hooks.server.ts for fail-fast validation
 */
export const serverEnv = serverEnvSchema.parse(env);

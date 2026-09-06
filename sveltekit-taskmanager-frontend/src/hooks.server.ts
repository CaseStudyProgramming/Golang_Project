import { serverEnv } from '$lib/server/env';

// Import serverEnv at the top to trigger fail-fast validation on app startup
console.log('Server environment validated successfully', serverEnv);

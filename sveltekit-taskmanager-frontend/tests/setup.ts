import { vi } from 'vitest';

// Mock environment variables
vi.mock('$env/static/public', () => ({
	env: {
		PUBLIC_API_BASE_URL: 'http://localhost:8080'
	}
}));

vi.mock('$env/static/private', () => ({
	env: {}
}));

# Task Manager Frontend

A modern, responsive task management application built with SvelteKit, TypeScript, and Tailwind CSS.

## Overview

This is the frontend component of the Task Manager application, providing a user-friendly interface for managing tasks, categories, tags, and analytics. It features a modern architecture with Svelte 5 runes, comprehensive testing, and production-ready build optimization.

## Tech Stack

- **Framework**: SvelteKit 2.x with Svelte 5
- **Language**: TypeScript
- **Styling**: Tailwind CSS 4.x
- **State Management**: Svelte 5 runes ($state, $derived, $props)
- **Testing**: Vitest (unit), Playwright (E2E)
- **Build Tool**: Vite 8.x
- **Package Manager**: Bun
- **Charts**: Chart.js with svelte-chartjs
- **Validation**: Zod

## Features

### Core Functionality

- User authentication (login, register, logout)
- Task CRUD operations with advanced filtering
- Task completion tracking and progress monitoring
- Category management with color coding
- Tag system for task organization
- Subtask management with checklist items
- Activity logging and audit trail
- Search and pagination
- Bulk operations (delete, complete)

### Advanced Features

- Analytics dashboard with visualizations
- CSV export functionality
- Responsive design (mobile-first)
- Real-time performance monitoring
- Web Vitals tracking
- Offline support preparation
- Multi-timezone support

### UI/UX Features

- Modern, clean interface design
- Loading states and skeleton screens
- Error handling with toast notifications
- Confirmation dialogs for destructive actions
- Empty states with helpful messaging
- Accessible components (ARIA labels, keyboard navigation)

## Prerequisites

- **Node.js**: 18+ or Bun latest
- **Backend API**: Task Manager Go backend running on port 8080
- **Git**: Latest version

## Installation

### 1. Clone the Repository

```bash
git clone https://github.com/CaseStudyProgramming/Golang_Project.git
cd taskmanager/sveltekit-taskmanager-frontend
```

### 2. Install Dependencies

```bash
bun install
```

Or with npm:

```bash
npm install
```

### 3. Configure Environment Variables

Create a `.env` file based on `.env.example`:

```bash
cp .env.example .env
```

Edit `.env` and set your backend API URL:

```env
PUBLIC_API_URL=http://localhost:8080
```

### 4. Start Development Server

```bash
bun run dev
```

Or with npm:

```bash
npm run dev
```

The application will be available at `http://localhost:5173`

## Development Commands

### Core Commands

```bash
# Start development server
bun run dev

# Build for production
bun run build

# Preview production build
bun run preview

# Type checking
bun run check

# Type checking in watch mode
bun run check:watch
```

### Code Quality

```bash
# Run linter
bun run lint

# Format code
bun run format

# Format check
bun run format:check
```

### Testing

```bash
# Run unit tests
bun test

# Run tests with UI
bun test:ui

# Run tests with coverage
bun test:coverage

# Run E2E tests
bun test:e2e

# Run E2E tests with UI
bun test:e2e:ui

# Run E2E tests in debug mode
bun test:e2e:debug
```

### Docker

```bash
# Build Docker image
bun run docker:build

# Run Docker container
bun run docker:run
```

## Project Structure

```
sveltekit-taskmanager-frontend/
├── src/
│   ├── lib/
│   │   ├── features/           # Feature-based architecture
│   │   │   ├── analytics/      # Analytics feature
│   │   │   ├── auth/           # Authentication feature
│   │   │   ├── categories/     # Category management
│   │   │   ├── tags/           # Tag management
│   │   │   └── tasks/          # Task management
│   │   ├── shared/             # Shared utilities and components
│   │   │   ├── components/     # Reusable UI components
│   │   │   ├── stores/         # Shared state management
│   │   │   ├── types/          # Shared TypeScript types
│   │   │   └── utils/          # Helper functions
│   │   ├── server/             # Server-side code
│   │   │   └── env.ts          # Environment validation
│   │   ├── assets/             # Static assets
│   │   └── index.ts           # Library exports
│   ├── routes/                 # SvelteKit file-based routing
│   │   ├── +layout.svelte      # Root layout
│   │   ├── +page.svelte        # Home page
│   │   ├── auth/               # Authentication routes
│   │   ├── dashboard/          # Dashboard routes
│   │   ├── tasks/              # Task management routes
│   │   ├── categories/         # Category routes
│   │   └── tags/               # Tag routes
│   ├── app.html                # HTML template
│   ├── app.d.ts               # Type declarations
│   └── hooks.server.ts        # Server hooks
├── static/                     # Static assets
├── tests/                      # Test setup
├── playwright-tests/           # E2E tests
├── public/                     # Public assets
├── package.json                # Dependencies and scripts
├── tsconfig.json              # TypeScript configuration
├── vite.config.ts             # Vite configuration
├── vitest.config.ts           # Vitest configuration
├── playwright.config.ts      # Playwright configuration
├── tailwind.config.js        # Tailwind configuration
├── .env.example              # Environment variables template
└── README.md                  # This file
```

## Component Usage

### Shared Components

Shared components are located in `src/lib/shared/components/` and can be imported from `$lib/shared/components`:

```svelte
<script>
  import { LoadingSpinner, Toast, EmptyState } from '$lib/shared/components';
</script>

<LoadingSpinner size="medium" />
<EmptyState title="No tasks found" icon="📝" />
```

Available shared components:

- `LoadingSpinner` - Loading indicator
- `Skeleton` - Loading placeholder
- `ProgressBar` - Progress indicator
- `EmptyState` - Empty state display
- `Toast` - Notification system
- `ConfirmDialog` - Confirmation dialogs
- `PasswordStrength` - Password strength indicator

### Feature Components

Feature-specific components are organized by domain:

```svelte
<script>
  import { TaskList, TaskForm, TaskFilters } from '$lib/features/tasks';
  import { CategoryList, CategoryForm } from '$lib/features/categories';
  import { TagInput } from '$lib/features/tags';
</script>
```

For detailed component documentation, see [COMPONENT_DOCUMENTATION.md](../COMPONENT_DOCUMENTATION.md).

## API Integration

The frontend communicates with the backend API through a centralized API utility:

```typescript
import { api } from '$lib/shared/utils/api.utils';

// Get tasks
const tasks = await api.getTasks({ page: 1, limit: 10 });

// Create task
const newTask = await api.createTask({
	title: 'New task',
	priority: 'HIGH'
});
```

For complete API documentation, see [API_INTEGRATION_GUIDE.md](../API_INTEGRATION_GUIDE.md).

## Environment Variables

### Required Variables

- `PUBLIC_API_URL` - Backend API URL (must start with PUBLIC_)

### Optional Variables

- `PUBLIC_APP_NAME` - Application name
- `PUBLIC_APP_VERSION` - Application version
- `PUBLIC_ENABLE_ANALYTICS` - Enable analytics feature
- `PUBLIC_ENABLE_EXPORT` - Enable export feature
- `PUBLIC_ENABLE_BULK_OPERATIONS` - Enable bulk operations
- `PUBLIC_SENTRY_DSN` - Sentry DSN for error tracking
- `PUBLIC_GOOGLE_ANALYTICS_ID` - Google Analytics ID
- `PUBLIC_DEFAULT_PAGE_SIZE` - Default pagination size
- `PUBLIC_MAX_FILE_SIZE` - Maximum file size in bytes
- `PUBLIC_DEFAULT_TIMEZONE` - Default timezone
- `PUBLIC_ENABLE_TIMEZONE_DETECTION` - Enable timezone detection

For detailed environment configuration, see [ENVIRONMENT_CONFIGURATION.md](../ENVIRONMENT_CONFIGURATION.md).

## Testing

### Unit Tests

Unit tests are co-located with implementation files using the `*.test.ts` pattern:

```typescript
// Example: src/lib/shared/utils/api.utils.test.ts
import { describe, it, expect } from 'vitest';
import { api } from './api.utils';

describe('API Utilities', () => {
	it('should make GET request', async () => {
		const result = await api.get('/test');
		expect(result).toBeDefined();
	});
});
```

### E2E Tests

E2E tests are located in `playwright-tests/`:

```typescript
// Example: playwright-tests/task-management.e2e.ts
import { test, expect } from '@playwright/test';

test('should create a new task', async ({ page }) => {
	await page.goto('/tasks');
	await page.click('[data-testid="add-task-button"]');
	await page.fill('[data-testid="task-title"]', 'Test task');
	await page.click('[data-testid="save-task"]');
	await expect(page.locator('text=Test task')).toBeVisible();
});
```

### Running Tests

```bash
# Unit tests
bun test

# E2E tests
bun test:e2e

# Coverage report
bun test:coverage
```

## Build and Deployment

### Production Build

```bash
bun run build
```

The optimized production build will be in the `build/` directory.

### Build Optimization

The production build includes:

- Code splitting and tree shaking
- Minification with Terser
- Asset optimization
- Source maps for debugging
- Bundle size optimization

For detailed build optimization information, see [BUILD_OPTIMIZATION.md](../BUILD_OPTIMIZATION.md).

### Deployment Options

#### Vercel

```bash
npm install -g vercel
vercel --prod
```

#### Netlify

```bash
npm install -g netlify-cli
netlify deploy --prod --dir=build
```

#### Docker

```bash
docker build -t taskmanager-frontend .
docker run -p 80:80 taskmanager-frontend
```

#### Static Hosting

Upload the `build/` directory to any static hosting service (AWS S3, CloudFront, etc.).

For detailed deployment instructions, see [DEPLOYMENT_GUIDE.md](../DEPLOYMENT_GUIDE.md).

## Performance Monitoring

The application includes built-in performance monitoring:

- Web Vitals tracking (FCP, LCP, FID, CLS, TTFB)
- Custom performance metrics (API response time, component render time)
- Performance data collection and reporting

```typescript
import { getPerformanceMonitor } from '$lib/shared/performance-monitoring';

const monitor = getPerformanceMonitor();
// Performance metrics are automatically collected
```

For detailed performance monitoring information, see [PERFORMANCE_MONITORING.md](../PERFORMANCE_MONITORING.md).

## Troubleshooting

### Common Issues

#### Development Server Won't Start

```bash
# Clear cache and reinstall
rm -rf node_modules .svelte-kit build
bun install
bun run dev
```

#### Type Errors

```bash
# Run type checking
bun run check

# If errors persist, check tsconfig.json
```

#### API Connection Issues

```bash
# Verify backend is running
curl http://localhost:8080/health

# Check environment variables
echo $PUBLIC_API_URL
```

#### Build Failures

```bash
# Clear build cache
rm -rf build .svelte-kit
bun run build
```

### Getting Help

- Check the [main README.md](../README.md) for project overview
- Review [API_INTEGRATION_GUIDE.md](../API_INTEGRATION_GUIDE.md) for API issues
- See [DEPLOYMENT_GUIDE.md](../DEPLOYMENT_GUIDE.md) for deployment help
- Check [AGENTS_FrontEnd.md](./AGENTS_FrontEnd.md) for development guidelines

## Development Guidelines

Follow the development guidelines specified in [AGENTS_FrontEnd.md](./AGENTS_FrontEnd.md):

- Use Svelte 5 runes for reactivity
- Follow vertical slice architecture
- Implement comprehensive testing
- Use TypeScript strict mode
- Follow OWASP security guidelines
- Maintain code quality standards

## Contributing

1. Follow the existing code style and patterns
2. Write tests for new features
3. Update documentation as needed
4. Follow the commit message conventions
5. Ensure all tests pass before submitting

## License

This project is licensed under the MIT License.

## Additional Resources

- [SvelteKit Documentation](https://kit.svelte.dev/docs)
- [Svelte 5 Documentation](https://svelte.dev/docs)
- [Tailwind CSS Documentation](https://tailwindcss.com/docs)
- [Vite Documentation](https://vitejs.dev)
- [Vitest Documentation](https://vitest.dev)
- [Playwright Documentation](https://playwright.dev)
- [Main Project README](../README.md)
- [Backend Documentation](../backendGoVanilaTaskmanager/readme.md)

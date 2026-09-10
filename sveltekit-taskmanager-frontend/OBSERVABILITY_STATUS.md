# Observability Status & Migration Guide

## Current Implementation Status

### ✅ Implemented

- **Error Handling Utilities**: Custom error classes in `src/lib/shared/utils/error.utils.ts`
  - `ApiError` for API-specific errors
  - `ValidationError` for validation errors
  - `AuthenticationError` for authentication failures
  - Global error handler registration system
  - User-friendly error message generation
- **API Interceptors**: Auth interceptors in `src/lib/shared/utils/auth.interceptors.ts`
  - Token management (automatic refresh)
  - Request/response interception
  - Error handling for API calls
- **Activity Logging**: Task activity logging in `src/lib/features/tasks/stores/activity.store.ts`
  - Task creation, update, deletion logging
  - User action tracking for task operations
- **SvelteKit Built-in Features**:
  - Server-side rendering with error handling
  - Basic performance metrics
  - Built-in error pages
- **Authentication State Management**: JWT-based auth with HTTP-only cookies
  - Auth state in `src/lib/features/auth/stores/auth.store.ts`
  - Automatic token refresh
  - Session management

### ❌ Not Implemented (Production Gaps)

- **Production Error Tracking**: No integration with error tracking services (Sentry, LogRocket, etc.)
- **Structured Logging**: No structured logging library or JSON format for logs
- **Performance Monitoring**: No dedicated performance monitoring service (Web Vitals, bundle analysis)
- **User Analytics**: No user analytics or behavior tracking (Google Analytics, Mixpanel, etc.)
- **Error Boundaries**: No custom error boundaries beyond SvelteKit defaults
- **Real User Monitoring (RUM)**: No RUM for real-time performance monitoring
- **Session Replay**: No session replay for debugging user issues
- **API Performance Monitoring**: No dedicated API response time tracking
- **Network Error Tracking**: No network failure monitoring and alerting
- **Feature Flag Management**: No feature flag system for gradual rollouts

## Migration Guide (Recommended for Production)

### Phase 1: Error Tracking

**Priority**: High
**Estimated Effort**: 2-3 days

**Goals**:

- Integrate production error tracking service
- Configure source maps for better error debugging
- Track JavaScript errors, API failures, and user interactions
- Set up alerting for critical errors

**Implementation Steps**:

1. Choose error tracking service:
   - **Sentry** (recommended, excellent for JavaScript/TypeScript)
   - **LogRocket** (includes session replay)
   - **Bugsnag** (good error tracking with release tracking)
   - **Rollbar** (real-time error tracking)

2. Install Sentry (example):

   ```bash
   bun add @sentry/sveltekit
   ```

3. Configure Sentry in SvelteKit:

   ```typescript
   // src/hooks.client.ts
   import * as Sentry from '@sentry/sveltekit';

   Sentry.init({
   	dsn: import.meta.env.VITE_PUBLIC_SENTRY_DSN,
   	tracesSampleRate: 1.0,
   	environment: import.meta.env.MODE,
   	release: import.meta.env.VITE_APP_VERSION
   });
   ```

   ```typescript
   // src/hooks.server.ts
   import * as Sentry from '@sentry/sveltekit';

   Sentry.init({
   	dsn: import.meta.env.VITE_PUBLIC_SENTRY_DSN,
   	environment: import.meta.env.MODE,
   	tracesSampleRate: 1.0
   });
   ```

4. Configure source maps:

   ```javascript
   // vite.config.ts
   export default defineConfig({
   	build: {
   		sourcemap: true // Enable source maps for production
   	}
   });
   ```

5. Add user context to errors:

   ```typescript
   // In error handlers
   Sentry.setUser({
   	id: user.id,
   	email: user.email
   });

   Sentry.setContext('task_operation', {
   	task_id: task.id,
   	operation: 'create'
   });
   ```

6. Set up error alerting in Sentry dashboard:
   - Alert on error rate > threshold
   - Alert on new errors introduced
   - Alert on specific error types (authentication failures, etc.)

**Files to Create**:

- `src/hooks.client.ts` (if not exists)
- Update `src/hooks.server.ts`

**Files to Modify**:

- `vite.config.ts`
- `src/lib/shared/utils/error.utils.ts` (integrate Sentry)
- `.env` (add Sentry DSN)

### Phase 2: Performance Monitoring

**Priority**: High
**Estimated Effort**: 2-3 days

**Goals**:

- Implement Web Vitals monitoring (LCP, FID, CLS)
- Track page load times and API response times
- Monitor bundle size and loading performance
- Set up performance budgets and alerts

**Implementation Steps**:

1. Add Web Vitals library:

   ```bash
   bun add web-vitals
   ```

2. Implement Web Vitals tracking:

   ```typescript
   // src/lib/shared/utils/performance.utils.ts
   import { onCLS, onFID, onLCP, onFCP, onTTFB } from 'web-vitals';

   export function initWebVitals() {
   	onCLS(metric => {
   		sendToAnalytics('CLS', metric);
   	});

   	onFID(metric => {
   		sendToAnalytics('FID', metric);
   	});

   	onLCP(metric => {
   		sendToAnalytics('LCP', metric);
   	});

   	onFCP(metric => {
   		sendToAnalytics('FCP', metric);
   	});

   	onTTFB(metric => {
   		sendToAnalytics('TTFB', metric);
   	});
   }

   function sendToAnalytics(name: string, metric: any) {
   	// Send to analytics service (Sentry, Google Analytics, etc.)
   	if (typeof window !== 'undefined' && window.Sentry) {
   		window.Sentry.captureMessage(`Web Vital: ${name}`, {
   			level: 'info',
   			extra: {
   				value: metric.value,
   				rating: metric.rating
   			}
   		});
   	}
   }
   ```

3. Initialize Web Vitals in app:

   ```typescript
   // src/routes/+layout.svelte
   import { onMount } from 'svelte';
   import { initWebVitals } from '$lib/shared/utils/performance.utils';

   onMount(() => {
   	initWebVitals();
   });
   ```

4. Add bundle size monitoring:

   ```bash
   bun add -D rollup-plugin-visualizer
   ```

   ```javascript
   // vite.config.ts
   import { visualizer } from 'rollup-plugin-visualizer';

   export default defineConfig({
   	plugins: [
   		visualizer({
   			open: true,
   			gzipSize: true,
   			brotliSize: true
   		})
   	]
   });
   ```

5. Set up performance budgets:

   ```javascript
   // vite.config.ts
   export default defineConfig({
   	build: {
   		rollupOptions: {
   			output: {
   				manualChunks: {
   					vendor: ['svelte', '@sveltejs/kit'],
   					ui: ['bits-ui', 'lucide-svelte']
   				}
   			}
   		}
   	}
   });
   ```

6. Monitor API response times:
   ```typescript
   // In api.utils.ts
   export async function apiRequest<T>(url: string, options?: RequestInit): Promise<T> {
   	const start = performance.now();

   	try {
   		const response = await fetch(url, options);
   		const duration = performance.now() - start;

   		// Track API performance
   		if (duration > 1000) {
   			console.warn(`Slow API request: ${url} took ${duration}ms`);
   			// Send to analytics
   		}

   		return response.json();
   	} catch (error) {
   		const duration = performance.now() - start;
   		// Track failed requests
   		throw error;
   	}
   }
   ```

**Files to Create**:

- `src/lib/shared/utils/performance.utils.ts`

**Files to Modify**:

- `src/routes/+layout.svelte`
- `vite.config.ts`
- `src/lib/shared/utils/api.utils.ts`

### Phase 3: User Analytics

**Priority**: Medium
**Estimated Effort**: 3-4 days

**Goals**:

- Implement user action tracking for critical flows
- Track feature usage and user engagement
- Set up funnel analysis for key user journeys

**Implementation Steps**:

1. Choose analytics service:
   - **Google Analytics 4** (free, widely used)
   - **Mixpanel** (excellent for product analytics)
   - **Amplitude** (good for user behavior analysis)
   - **PostHog** (open-source alternative)

2. Install analytics library (example with Google Analytics):

   ```bash
   bun add @sveltejs/analytics
   ```

3. Configure analytics:

   ```typescript
   // src/hooks.client.ts
   import { analytics } from '@sveltejs/analytics/google';

   analytics('G-XXXXXXXXXX', {
   	pageview: true
   });
   ```

4. Track custom events:

   ```typescript
   // src/lib/shared/utils/analytics.utils.ts
   export function trackEvent(eventName: string, properties?: Record<string, any>) {
   	if (typeof window !== 'undefined' && window.gtag) {
   		window.gtag('event', eventName, properties);
   	}
   }

   export function trackUserAction(action: string, context?: Record<string, any>) {
   	trackEvent('user_action', {
   		action,
   		...context
   	});
   }
   ```

5. Track critical user flows:

   ```typescript
   // In auth.store.ts
   async login(email: string, password: string) {
     try {
       trackUserAction('login_attempt', { email });
       const response = await api.login(email, password);
       trackUserAction('login_success', { user_id: response.user.id });
       return response;
     } catch (error) {
       trackUserAction('login_failure', { error: error.message });
       throw error;
     }
   }
   ```

6. Track feature usage:

   ```typescript
   // In task store
   async createTask(taskData: TaskData) {
     trackUserAction('task_create', {
       priority: taskData.priority,
       category: taskData.category_id,
     });
     // ... create task logic
   }
   ```

7. Set up funnel analysis:
   - Track registration funnel: visit → signup → email verification → first task
   - Track task completion funnel: create → start → complete
   - Track feature adoption: tags, categories, subtasks usage

**Files to Create**:

- `src/lib/shared/utils/analytics.utils.ts`

**Files to Modify**:

- `src/hooks.client.ts`
- `src/lib/features/auth/stores/auth.store.ts`
- `src/lib/features/tasks/stores/task.store.ts`

### Phase 4: Error Boundaries & UX Improvements

**Priority**: Medium
**Estimated Effort**: 2-3 days

**Goals**:

- Implement custom error boundaries for better error UX
- Add error recovery mechanisms
- Implement fallback UI for critical failures

**Implementation Steps**:

1. Create custom error pages:

   ```svelte
   <!-- src/routes/+error.svelte -->
   <script lang="ts">
     import { page } from '$app/stores';
     import { trackError } from '$lib/shared/utils/analytics.utils';

     $: {
       trackError('page_error', {
         status: $page.status,
         url: $page.url.pathname,
       });
     }
   </script>

   <div class="error-container">
     <h1>{$page.status}</h1>
     <p>{$page.error?.message || 'Something went wrong'}</p>
     <button on:click={() => window.location.reload()}>Try Again</button>
     <a href="/">Go Home</a>
   </div>

   <style>
     .error-container {
       display: flex;
       flex-direction: column;
       align-items: center;
       justify-content: center;
       min-height: 100vh;
       padding: 2rem;
     }
   </style>
   ```

2. Implement component-level error boundaries:

   ```svelte
   <!-- src/lib/shared/components/ErrorBoundary.svelte -->
   <script lang="ts">
     import { onMount } from 'svelte';

     let hasError = false;
     let errorMessage = '';

     onMount(() => {
       window.addEventListener('error', handleError);
       window.addEventListener('unhandledrejection', handlePromiseRejection);
     });

     function handleError(event: ErrorEvent) {
       hasError = true;
       errorMessage = event.error?.message || 'Unknown error';
       // Track error
     }

     function handlePromiseRejection(event: PromiseRejectionEvent) {
       hasError = true;
       errorMessage = event.reason?.message || 'Unknown error';
       // Track error
     }

     function retry() {
       hasError = false;
       errorMessage = '';
       // Retry logic
     }
   </script>

   {#if hasError}
     <div class="error-fallback">
       <p>{errorMessage}</p>
       <button on:click={retry}>Retry</button>
     </div>
   {:else}
     <slot />
   {/if}
   ```

3. Implement loading states and skeletons:

   ```svelte
   <!-- Example skeleton loading -->
   {#if loading}
     <div class="skeleton-loader">
       <div class="skeleton-header" />
       <div class="skeleton-content" />
     </div>
   {:else}
     <div class="content">
       <!-- Actual content -->
     </div>
   {/if}

   <style>
     .skeleton-loader {
       animation: pulse 1.5s ease-in-out infinite;
     }

     @keyframes pulse {
       0%, 100% { opacity: 1; }
       50% { opacity: 0.5; }
     }
   </style>
   ```

4. Implement offline support:
   ```typescript
   // src/lib/shared/utils/offline.utils.ts
   export function initOfflineSupport() {
   	window.addEventListener('online', () => {
   		// Sync offline changes
   	});

   	window.addEventListener('offline', () => {
   		// Show offline notification
   	});
   }

   export function isOnline(): boolean {
   	return navigator.onLine;
   }
   ```

**Files to Create**:

- `src/routes/+error.svelte`
- `src/lib/shared/components/ErrorBoundary.svelte`
- `src/lib/shared/utils/offline.utils.ts`

**Files to Modify**:

- `src/routes/+layout.svelte` (init offline support)

## Production Checklist

Before deploying to production, ensure:

- [ ] Error tracking service integrated (Sentry, LogRocket, etc.)
- [ ] Source maps configured for error debugging
- [ ] Web Vitals monitoring implemented (LCP, FID, CLS)
- [ ] Performance budgets defined and monitored
- [ ] Bundle size analysis automated
- [ ] API response time tracking implemented
- [ ] User analytics integrated (GA4, Mixpanel, etc.)
- [ ] Critical user flows tracked (login, task operations)
- [ ] Custom error pages implemented
- [ ] Error recovery mechanisms in place
- [ ] Loading states and skeletons implemented
- [ ] Offline support and graceful degradation
- [ ] Performance monitoring dashboards set up
- [ ] Alerting configured for critical errors
- [ ] User funnel analysis implemented
- [ ] Session replay (if using LogRocket or similar)

## Monitoring & Alerting Recommendations

### Key Metrics to Monitor

- **Error Rate**: JavaScript errors, API failures
- **Performance**: LCP < 2.5s, FID < 100ms, CLS < 0.1
- **User Engagement**: Session duration, page views per session
- **Conversion**: Registration rate, task completion rate
- **API Performance**: Response times, failure rates
- **Bundle Size**: Initial load size, chunk sizes

### Alert Thresholds (Example)

- Error rate > 1% for 5 minutes
- LCP > 4s for 10% of users
- API failure rate > 5% for 5 minutes
- Registration conversion rate drops > 20%
- New error introduced in production

### Recommended Tools

- **Error Tracking**: Sentry, LogRocket, Bugsnag
- **Performance Monitoring**: Google Analytics Web Vitals, Lighthouse CI
- **User Analytics**: Google Analytics 4, Mixpanel, Amplitude
- **Session Replay**: LogRocket, FullStory
- **Bundle Analysis**: Bundlephobia, webpack-bundle-analyzer
- **Real User Monitoring**: SpeedCurve, Calibre

## Privacy & Compliance Considerations

- **User Consent**: Implement cookie consent for analytics
- **Data Minimization**: Only collect necessary analytics data
- **Anonymization**: Anonymize user IDs and PII in analytics
- **GDPR Compliance**: Provide data export and deletion options
- **Regional Restrictions**: Respect regional data regulations
- **Cookie Management**: Allow users to opt-out of tracking

# AGENTS.md

## Core Directives & Priority (Precedence Order)

1. **Security First**: OWASP guidelines override all formatting and style rules.
2. **Architecture Integrity**: Vertical Slice rules override file-level preferences.
3. **Functionality > Style**: Working implementation takes priority over code style/refactoring.
4. **Circuit Breaker**: Stop and ask after 3 consecutive failed test/type-check attempts.
5. **When in Doubt, Ask**: Stop and request clarification if user instructions contradict these rules or if a requirement is ambiguous.

## Workflow & Task Execution

- **Branch-Based Development**:
  - Always work on a new branch for each sub-issue using git branch standard naming convention (e.g., `feat/setup-sveltekit-frontend`).
- **Task Breakdown via Granular Checklists / Issues**:
  - Read and parse task lists (`- [ ]`) or sub-issues sequentially.
  - Work on only ONE granular checklist item at a time. Complete it fully before moving to the next.
- **Atomic Commits & Verification**:
  - Run verification commands (`bun run check`, `bun test`) after completing each individual checklist item.
  - Immediately create a dedicated Git commit upon completing each item (e.g., `feat(scope): complete task X - description`).
  - Update or check off the task list item upon a successful commit.
- **Rollback Readiness**:
  - Each completed task MUST correspond to a clean, isolated Git commit to allow single-step rollbacks (`git revert`) without losing previous progress.

## Token Efficiency

- Skip recaps and conversational summaries unless the result is ambiguous or requires further input.

## Principles

- **OWASP Security Standard**: Validate all inputs with Zod/Valibot, prevent XSS/injection, and leverage SvelteKit CSRF/CORS protections, OWASP Top 10
- **Clarity and Consistency**: Clarity over cleverness. Match existing code patterns. Minimal changes unless refactoring is explicitly requested.
- **Modularity**: Keep functions under 50 lines and components under 200 lines. Break down when improves readability and structure.
- **Type Safety**: TypeScript everywhere including test files. Use `any` ONLY when dealing with external libraries without types, and document why with `@ts-expect-error` comments.
- **Error Handling**: Avoid unnecessary `try/catch`. Use type narrowing over type casting. Use SvelteKit error boundaries for critical failures.
- **Exports & Routing**: Use named exports for TypeScript modules. Follow SvelteKit file-based routing conventions (`+page.svelte`, `+page.server.ts`, `+layout.svelte`, etc.).
- **Path Aliases**: Use absolute imports via `$lib/...` (SvelteKit standard path alias).
- **Tooling**: Follow existing ESLint, Prettier, or Biome setups; do not reformat unrelated code.
- **Schema Validation:** Import Zod as a value (`import { z } from 'zod';`) when building/parsing runtime schemas (e.g., `env.ts`, form actions). Use `import type` ONLY when importing inferred TypeScript types (`type User = z.infer<typeof userSchema>`).
- **Type Inference**: Let the compiler infer return types unless: (1) function is exported, (2) return type is complex, or (3) explicit annotation improves clarity.
- **Function Parameters**: Use an options object for functions with 3+ parameters, optional flags, or ambiguous arguments.
- **Debugging**: Hypothesis-driven debugging—formulate 1–3 most likely causes first, then validate incrementally.

## Type Safety Tooling Stack

- **Three-Layer Type Safety Approach**:
  1. **Incremental TypeScript Compilation** (`tsconfig.json`): Performance foundation with caching for fast type checking (70-90% faster than full compilation).
  2. **Vitest Type Checking Integration** (`vitest.config.ts`): Developer experience with instant type feedback during test runs.
  3. **CI/CD Type Checking Gate**: Final safety net before production deployment.

- **TypeScript Configuration**:
  - Enable `strict: true` for maximum type safety across all files including tests.
  - Enable `incremental: true` with `tsBuildInfoFile: ".tsbuildinfo"` for performance optimization.
  - Include test files in type checking: remove test file exclusions from `tsconfig.json`.
  - Use tiered approach if needed: strict for production code, relaxed for test files with separate `tsconfig.test.json`.

- **Vitest Type Checking**:
  - Enable `typecheck` configuration in `vitest.config.ts` for integrated type checking during test runs.
  - Use same `tsconfig.json` for consistency between development and CI/CD.
  - Type errors in tests will fail test runs, providing immediate feedback.

- **CI/CD Integration**:
  - Always run type checking (`bun run check`) before test execution in CI/CD pipelines.
  - Type checking failures should block deployment to production.
  - Use caching strategies (GitHub Actions cache, TurboRepo, Nx) for faster CI/CD builds.

- **Performance Best Practices**:
  - Target type checking time < 500ms for medium projects using incremental compilation.
  - Monitor type checking performance; if > 2s, investigate project references or monorepo optimization.
  - Leverage modern tooling (Vite, esbuild, swc, Bun) for faster compilation.

- **TypeScript Configuration Best Practices**:
  - **tsconfig.json**: Base configuration with strict mode enabled for all files including tests. Remove test file exclusions from `exclude` array.
  - **Incremental Compilation**: Always enable for performance - set `incremental: true` and `tsBuildInfoFile: ".tsbuildinfo"`.
  - **No Test Exclusions**: Never exclude test files from type checking in base tsconfig.json. If relaxed type checking for tests is needed, use separate `tsconfig.test.json` that extends base config.
  - **Path Aliases**: Use SvelteKit's built-in `$lib` alias and configure additional aliases in `tsconfig.json` and `vite.config.ts` for consistency.
  - **Module Resolution**: Use `"moduleResolution": "bundler"` for modern ESM projects (SvelteKit default).

- **Handling Edge Cases**:
  - For external libraries without types: use `@ts-expect-error` with explanatory comments, never globally disable type checking.
  - For complex mocking scenarios: use relaxed type config for test files if strict mode causes significant DX friction.
  - Document any type compromises with clear rationale in code comments.

## System Quality & Reliability (Frontend/Client-Side Scope)

- **System Observability & Incident Response**: Implement client-side error tracking (e.g., Sentry, LogRocket), logging for user actions, and performance monitoring. Ensure auditability of all critical operations in server-side SvelteKit code.
- **High Availability & Fault Tolerance (HA/FT)**: Implement graceful degradation for API failures, offline support via service workers, and retry mechanisms with exponential backoff for network requests.
- **API Defensive Design (Defensive Programming) & Code Quality**: Validate all API responses, implement error boundaries, use defensive coding practices. Follow SOLID principles and maintain high code coverage with meaningful tests.
- **Performance & Concurrency**: Optimize bundle size, implement lazy loading and code splitting, handle concurrent state updates safely, use debouncing/throttling for expensive operations.
- **Security**: Follow OWASP Top 10 guidelines (XSS, CSRF, injection prevention), implement content security policy, use secure communication (HTTPS), and regularly update dependencies. Never expose secrets or sensitive data in client code.
- **Usability / Robustness**: Design intuitive error messages, implement graceful error handling, and ensure the system provides helpful feedback. Design for edge cases and unexpected user behavior.
- **Data Privacy & Information Disclosure Protection**: Implement data minimization in client storage, avoid logging sensitive information, and comply with privacy regulations (GDPR, CCPA). Never store PII in localStorage/sessionStorage without encryption.

**Note**: Infrastructure, database resiliency, traffic control at server level, and data consistency rules apply to backend systems. Frontend should handle UI-level error states and retry logic for API calls.


## Commands (Bun)
- Always use `bun` as the package manager and test/runtime runner:
  - `bun run dev` - Start dev server
  - `bun run build` - Build for production
  - `bun run preview` - Preview production build
  - `bun run check` - Svelte & TypeScript type checking for all files including tests (`svelte-check`)
  - `bun run check:watch` - Watch mode for type checking with instant feedback
  - `bun run lint` - Run linter
  - `bun run format` - Format code
  - `bun test` - Run unit tests with Vitest (includes type checking integration)
  - `bun test:ui` - Run tests with Vitest UI interface
  - `bun test:coverage` - Run tests with coverage report

## Git Commits

- **Conventional Commits**: Format as `type: summary without scope`.
- Summary must be a short, specific sentence explaining what changed and why.
- Valid types: `feat|fix|docs|style|refactor|perf|test|build|ci|chore|revert`.
- Include `BREAKING CHANGE:` in the commit footer when applicable.

## Environment Variables

- **No Direct Reads:** NEVER read `process.env` directly. Do NOT import raw `$env/*` in UI components, routes (`+page.svelte`), or feature modules.
- **Centralized Validation:** Route all environment variables through Zod validation files before app consumption:
  - **Server Secrets:** Define and validate in `$lib/server/env.ts` using `$env/static/private` (or `$env/dynamic/private` for Docker runtime deployments).
  - **Public Variables:** Define and validate in `$lib/env.ts` using `$env/static/public` (must start with `PUBLIC_`).
- **Feature-Specific Secrets Exception:** If a secret belongs strictly to a single feature slice, place its validation inside `src/lib/features/{feature}/server/env.ts` importing from `$lib/server/env.ts`.
- **Startup Boot Validation:** Import `$lib/server/env` at the top of `src/hooks.server.ts` to trigger fail-fast validation upon app startup.
- **Consumption Standard:** Always import `serverEnv` or `publicEnv` exported objects in your feature logic.

## Styling

- Use latest Tailwind CSS utility classes.
- Reuse shared components (e.g., shadcn-svelte, Bits UI, or custom UI components).
- Mobile-first, fully responsive design.
- Avoid redundant, unused, or conflicting utility classes.

## Architecture & Vertical Slices

- **Vertical Slice Architecture**:
  - Organize code by business domain features (`src/lib/features/`) rather than technical layers.
  - Feature structure: `src/lib/features/{feature}/{api,components,stores,types,index.ts}`.
  - Each feature slice must be independent and self-contained.
- **Shared vs Feature-Specific**:
  - **Shared (`src/lib/shared/` or `$lib/`)**: Native `fetch` wrappers/helpers (Do NOT use `axios`), authentication, error handling, validation schemas, UI primitives.
  - **Feature-Specific**: Business logic, domain types, feature UI components, local stores.
- **Dependency Direction & Communication**:
  - Features can depend on shared utilities, but features must NOT depend directly on other features.
  - Cross-slice communication must use well-defined TypeScript interfaces or shared stores in `$lib/stores/`.
- **Slice Boundaries & Extraction**:
  - A slice should contain 1–3 related business concepts. If a slice exceeds 10 files in subdirectories, split it into smaller slices.
  - **Extraction Rule**: Do NOT estimate code percentages. If 3 or more functions, types, or UI components are identical across 2 different slices, extract them directly into `$lib/shared/`.
- **Error Handling**:
  - Handle domain-specific errors inside the slice. Use SvelteKit error boundaries for critical failures.

## Svelte & SvelteKit

- **Svelte 5 Runes & Reactive State**:
  - Prefer modern Svelte 5 runes (`$state`, `$derived`, `$props`, `$effect`) over legacy reactive statements (`$: ...`).
  - Keep state minimal; compute dependent values using `$derived`.
  - Avoid using `$effect` for state synchronization; use it exclusively for side effects and direct DOM interactions.
  - Component props syntax: `let { foo, bar }: Props = $props();` or `const { foo, bar } = $props();`.
- **SvelteKit Architecture & Routing**:
  - Import generated types directly from `./$types` for load functions and form actions (`PageServerLoad`, `PageData`, `ActionData`, `Actions`, `LayoutServerLoad`).
  - Use SvelteKit **Form Actions** (`export const actions = { ... }`) for data mutations instead of raw custom POST API endpoints.
  - Keep server-only business logic and database queries strictly isolated in `$lib/server/`.
  - Use native `fetch` with custom wrappers; do NOT use or install `axios`.
  - Use SvelteKit built-in navigation utilities (`goto`, `redirect`, `error`).
- **Global/Shared State**:
  - Extract reusable reactive logic into .svelte.ts or .svelte.js files using $state classes or functions instead of using legacy Svelte stores (writable, readable)
- **External Backend Communication**:
  - All communications with external backend services (e.g., Go microservices or custom APIs) that require sensitive credentials must be routed through SvelteKit's server-side functions (+page.server.ts or +server.ts) to ensure secrets are never exposed to the client

## JSDoc

- Start each block with `/**` directly above the symbol.
- Write short, sentence-case, present-tense descriptions of intent.
- Tag order: description → `@param` → `@returns` → `@throws` (only if it can throw).

## Tests

- **Type-Safe Testing**: All test files (`*.test.ts`, `*.spec.ts`) are included in TypeScript type checking. Tests must pass type checking to be considered valid.
- **Test File Organization**:
  - Co-locate unit and integration tests (`*.test.ts`) with implementation files.
  - Place Playwright E2E tests (`*.spec.ts` or `*.e2e.ts`) in the `tests/` directory.
- **Test Structure**: Top `describe` = subject; nested `describe` = scenarios/contexts.
- **Test Naming**: `it` titles: short, third-person present, `verb + object + context` (sentence case, no period). Omit words like "should/works/handles".
- **Mocking Strategy**: Avoid unnecessary mocking unless dealing with external network or hardware I/O. When mocking is required, use type-safe mocks that maintain type safety.
- **Type Safety in Tests**:
  - Use `@ts-expect-error` with explanatory comments for intentional type violations (e.g., mocking untyped libraries).
  - Leverage TypeScript for test data validation and API response typing.
  - Ensure test doubles, mocks, and fixtures maintain type contracts with production code.
- **Type-Safe Test Helpers**:
  - Extract reusable logic into helper functions to avoid TypeScript literal type narrowing issues
  - Helper functions accept union types as parameters to enable proper type inference without type assertions
  - This maintains full type safety while solving compile-time comparison errors with literal values
  - Helpers are test-only utilities (excluded from production builds) and do not represent production code logic
- **Negative Testing**:
  - Test edge cases and invalid inputs to ensure functions reject illegal data appropriately
  - Use domain data validation (Zod schemas) to validate boundary conditions for external data (API responses, form payloads, database schemas)
  - Type safety provides compile-time guardrails; runtime validation (Zod) provides business logic guardrails
  - Type safety does NOT prevent random testing - it ensures random data is type-valid while still allowing variation
- **Data Generators for Mock Data**:
  - **Trigger Condition**: Create data generators/factories ONLY when an entity is used in >5 different test files AND has >5 properties
  - **Implementation Options**: 
    - **@faker-js/faker** (Recommended for multi-users, multi-regions, realistic data): Use for user profiles, regional data (timezones, formats), collaboration scenarios
    - **Fishery pattern** (Factory pattern): Use for business logic states, known scenarios, deterministic behavior
    - **Hybrid approach**: Combine both - use Faker for realistic data (names, emails, regions) and Fishery for known states (pending, completed, assigned)
    - **Custom factories**: Use for simple needs without library dependency
  - **Gradual Migration**: Do not migrate all entities at once—prioritize entities that change most frequently
  - **Manual Construction**: For entities below threshold, continue using manual mock construction
  - **Benefits**: Reduces boilerplate, ensures consistency, easier to update when entity structure changes, supports realistic multi-user/multi-region testing
- **Comprehensive Testing Approach**:
  - Type Safety + Test Helpers + Negative Testing + (Conditional Data Generators) = Type-safe comprehensive testing with negative case coverage
  - This combination ensures compile-time type safety, runtime validation, comprehensive edge case coverage, and maintainable test code
- **Vitest Integration**: Run `bun test` to execute tests with integrated type checking. Type errors in tests will fail the test suite.

## CI/CD & Type Safety Gates

- **Type Checking Gate**: All CI/CD pipelines MUST run `bun run check` before test execution and deployment. Type checking failures should block the pipeline.
- **Test Execution Gate**: Run `bun test` after successful type checking. Both gates must pass for deployment.
- **Pipeline Order**:
  1. Type check (all files including tests)
  2. Lint check
  3. Unit tests (with type checking)
  4. E2E tests
  5. Build
- **Caching Strategy**: Implement TypeScript build cache and dependency caching in CI/CD for faster builds (GitHub Actions cache, TurboRepo, or similar).
- **Branch Protection**: Enable branch protection rules requiring type checking and test status to pass before merging to main branch.
- **Rollback Safety**: Maintain atomic commits with clear type checking verification to enable safe rollbacks if type-related issues are discovered post-deployment.
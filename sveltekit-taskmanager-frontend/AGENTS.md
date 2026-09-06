# AGENTS.md

## Core Directives & Priority
1. **Security First**: OWASP guidelines override all formatting and style rules.
2. **Architecture Integrity**: Vertical Slice rules override file-level preferences.
3. **Circuit Breaker**: Stop and ask after 3 consecutive failed test/type-check attempts.
4. **When in Doubt, Ask**: Stop and request clarification if user instructions contradict these rules or if a requirement is ambiguous.

## Principles
- **OWASP Security Standard**: Validate all inputs with Zod/Valibot, prevent XSS/injection, and leverage SvelteKit CSRF/CORS protections, OWASP Top 10
- **Clarity and Consistency**: Clarity over cleverness. Minimal changes. Match existing code patterns.
- **Modularity**: Keep components and functions short; break down when it improves readability and structure.
- **Type Safety**: TypeScript everywhere. No `any` unless isolated, documented, and strictly necessary.
- **Error Handling**: Avoid unnecessary `try/catch`. Use type narrowing over type casting.
- **Exports & Routing**: Use named exports for TypeScript modules. Follow SvelteKit file-based routing conventions (`+page.svelte`, `+page.server.ts`, `+layout.svelte`, etc.).
- **Path Aliases**: Use absolute imports via `$lib/...` (SvelteKit standard path alias).
- **Tooling**: Follow existing ESLint, Prettier, or Biome setups; do not reformat unrelated code.
- **Schema Validation:** Import Zod as a value (`import { z } from 'zod';`) when building/parsing runtime schemas (e.g., `env.ts`, form actions). Use `import type` ONLY when importing inferred TypeScript types (`type User = z.infer<typeof userSchema>`).
- **Type Inference**: Let the compiler infer return types unless explicit annotation improves clarity.
- **Function Parameters**: Use an options object for functions with 3+ parameters, optional flags, or ambiguous arguments.
- **Debugging**: Hypothesis-driven debugging—formulate 1–3 most likely causes first, then validate incrementally.

## Workflow & Task Execution
- **Branch-Based Development**:
  - Always work on a new branch for each sub-issue
- **Task Breakdown via Checklists / Issues**:
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

## Commands (Bun)
- Always use `bun` as the package manager and test/runtime runner:
  - `bun run dev` - Start dev server
  - `bun run build` - Build for production
  - `bun run preview` - Preview production build
  - `bun run check` - Svelte & TypeScript type checking (`svelte-check`)
  - `bun run lint` - Run linter
  - `bun run format` - Format code
  - `bun test` - Run unit tests with Bun / Vitest

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

## JSDoc
- Start each block with `/**` directly above the symbol.
- Write short, sentence-case, present-tense descriptions of intent.
- Tag order: description → `@param` → `@returns` → `@throws` (only if it can throw).

## Tests
- Co-locate unit and integration tests (`*.test.ts`) with implementation files.
- Place Playwright E2E tests (`*.spec.ts` or `*.e2e.ts`) in the `tests/` directory.
- Test structure: Top `describe` = subject; nested `describe` = scenarios/contexts.
- `it` titles: short, third-person present, `verb + object + context` (sentence case, no period). Omit words like "should/works/handles".
- Avoid unnecessary mocking unless dealing with external network or hardware I/O.
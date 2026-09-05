# AGENTS.md

## Principles

- OWASP Approach
- Clarity and consistency over cleverness. Minimal changes. Match existing patterns.
- Keep components/functions short; break down when it improves structure.
- TypeScript everywhere; no `any` unless isolated and necessary.
- No unnecessary `try/catch`.
- Avoid casting; use narrowing.
- Named exports for TS modules. Follow SvelteKit file-based routing conventions (`+page.svelte`, `+page.server.ts`, `+layout.svelte`, etc.).
- Absolute imports via `$lib/...` (SvelteKit standard path alias).
- Follow existing ESLint / Prettier / Biome setup; don't reformat unrelated code.
- Zod / Valibot type-only import when schema parsing: `import type * as z from 'zod';`.
- Let compiler infer return types unless annotation adds clarity.
- Options object for 3+ params, optional flags, or ambiguous args.
- Hypothesis-driven debugging: 1–3 causes, validate most likely first.

## Token Efficiency

- Skip recaps unless the result is ambiguous or you need more input.

## Commands (Bun)

- Always use `bun` as the package manager and test/runtime runner.
- Core commands:
  - `bun run dev` - Start dev server
  - `bun run build` - Build for production
  - `bun run preview` - Preview production build
  - `bun run check` - Svelte & TypeScript type checking (`svelte-check`)
  - `bun run lint` - Run linter
  - `bun run format` - Format code
  - `bun test` - Run unit tests with Bun / Vitest

## Git Commits

- Conventional Commits: `type: summary without scope`.
- The summary should be a short, specific sentence that explains what changed and where or why, not a vague phrase.
- Types: `feat|fix|docs|style|refactor|perf|test|build|ci|chore|revert`.
- `BREAKING CHANGE:` in footer when needed.

## Environment Variables

- Always use SvelteKit's built-in environment modules (`$env/static/private`, `$env/static/public`, `$env/dynamic/private`, `$env/dynamic/public`).
- Never read `process.env` directly.
- Ensure private server secrets are placed inside `$env/static/private` or `$lib/server/` to prevent accidental client exposure.

## Styling

- Tailwind CSS utility classes latest.
- Reuse shared components (e.g., shadcn-svelte, Bits UI, or custom UI components).
- Mobile-first and fully responsive layout design.
- Avoid redundant, unused, or conflicting utility classes.

## Architecture & Vertical Slices

- **Vertical Slice Architecture**:
  - Organize code by business features (`src/features/`) rather than technical layers
  - Each feature slice should be independent and self-contained
  - Shared infrastructure (HTTP clients, utilities) goes in `src/shared/` or `src/lib/`
  - Business logic stays within feature slices, infrastructure can be shared
  - Avoid duplicating infrastructure code across slices - compose/extend shared utilities instead
  - Feature structure: `src/features/{feature}/{api,components,stores,types,index.ts}`
- **Shared vs Feature-Specific**:
  - Shared: HTTP wrappers, authentication, error handling, validation schemas, date utilities
  - Feature-specific: Business logic, feature-specific types, UI components, feature stores
- **Dependency Direction**:
  - Features can depend on shared utilities
  - Features should not depend on other features (unless well-defined interfaces)
  - Shared utilities should not depend on features

## Svelte & SvelteKit

- **Svelte 5 or latest Runes & Reactive State**:
  - Prefer modern Svelte 5 runes (`$state`, `$derived`, `$props`, `$effect`) over legacy reactive statements (`$: ...`).
  - Keep state minimal; compute dependent values using `$derived` instead of manually syncing `$state`.
  - Avoid using `$effect` for state synchronization; use it solely for side effects and direct DOM interactions.
  - Component props syntax: `let { foo, bar }: Props = $props();` or `const { foo, bar } = $props();`.
- **SvelteKit Architecture & Routing**:
  - Import generated types directly from `./$types` for load functions and form actions (`PageServerLoad`, `PageData`, `ActionData`, `Actions`, `LayoutServerLoad`).
  - Leverage SvelteKit **Form Actions** (`export const actions = { ... }`) for data mutations instead of writing raw custom POST API endpoints.
  - Keep server-only business logic and database queries strictly isolated in `$lib/server/`.
  - Use SvelteKit built-in navigation utilities (`goto`, `redirect`, `error`) instead of manual fetch/window routing hacks.

## JSDoc

- Start each block with `/**` directly above the symbol.
- Short, sentence-case, present-tense description of intent.
- Order: description → `@param` → `@returns` → `@throws` (only if it can throw).

## Tests

- `*.test.ts` for unit/integration tests co-located with their implementation files.
- `*.spec.ts` or `*.e2e.ts` in `tests/` directory for Playwright E2E tests.
- Top `describe` = subject; nested `describe` to group scenarios or contexts.
- `it` titles: short, third-person present, `verb + object + context` (sentence case, no period). Omit "should/works/handles". State _what_, not _how_.
- Avoid unnecessary mocking unless dealing with external network or hardware I/O.

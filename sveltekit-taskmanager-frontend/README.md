# Task Manager Frontend

A production-ready SvelteKit task manager application with TypeScript, Tailwind CSS, and comprehensive testing setup.

## Prerequisites

- [Bun](https://bun.sh/) - JavaScript runtime and package manager
- Node.js >= 18 (optional, as Bun is the primary runtime)

## Installation

```bash
# Install dependencies
bun install

# Copy environment variables
cp .env.example .env
```

## Development

```bash
# Start development server
bun run dev

# The app will be available at http://localhost:5173
```

## Available Scripts

```bash
# Development
bun run dev          # Start development server
bun run build        # Build for production
bun run preview      # Preview production build

# Code Quality
bun run check        # Run TypeScript type checking
bun run lint         # Run ESLint
bun run format       # Format code with Prettier

# Testing
bun test             # Run unit tests
bun test:ui          # Run tests with UI
```

## Project Structure

```
src/
├── features/        # Feature-based modules
│   └── tasks/       # Task management feature
│       ├── api/     # API client functions
│       ├── components/  # UI components
│       ├── stores/  # Svelte stores
│       └── types/   # TypeScript types
├── lib/             # Shared libraries
│   ├── client/      # Client-side utilities
│   └── server/      # Server-side utilities
├── routes/          # SvelteKit routes
├── shared/          # Shared utilities and types
│   ├── types/       # Shared TypeScript types
│   └── utils/       # Shared utility functions
└── static/          # Static assets
```

## Environment Variables

Copy `.env.example` to `.env` and configure:

```env
VITE_API_BASE_URL=http://localhost:8080/api
VITE_APP_NAME=Task Manager
VITE_APP_VERSION=1.0.0
```

## Technology Stack

- **Framework**: SvelteKit with Svelte 5 (Runes)
- **Language**: TypeScript
- **Styling**: Tailwind CSS
- **State Management**: Svelte 5 Runes
- **Testing**: Vitest + Testing Library
- **Code Quality**: ESLint + Prettier
- **Build Tool**: Vite
- **Package Manager**: Bun

## Development Guidelines

- Follow the coding standards defined in `AGENTS.md`
- Use Svelte 5 runes (`$state`, `$derived`, `$props`) instead of legacy reactive statements
- Implement OWASP security best practices
- Write unit tests for new features
- Run `bun run lint` and `bun run check` before committing
- Use conventional commits for commit messages

## License

MIT

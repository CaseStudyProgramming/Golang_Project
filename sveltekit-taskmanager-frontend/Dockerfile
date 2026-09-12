# Multi-stage build for optimized production image
FROM node:24-alpine AS builder

# Set working directory
WORKDIR /app

# Install bun for faster builds
RUN npm install -g bun

# Copy package files for dependency caching
COPY package.json bun.lock ./

# Install dependencies
RUN bun install --frozen-lockfile

# Copy source code
COPY . .

# Build the application
RUN bun run build

# Runtime stage - nginx for serving static files
FROM nginx:alpine

# Copy custom nginx config
COPY nginx.conf /etc/nginx/nginx.conf

# Copy built assets from builder
COPY --from=builder /app/.svelte-kit/output/client /usr/share/nginx/html

# Copy .htaccess for routing (if needed)
COPY --from=builder /app/.svelte-kit/output/client/.htaccess /usr/share/nginx/html/.htaccess 2>/dev/null || true

# Expose port
EXPOSE 80

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost/ || exit 1

# Start nginx
CMD ["nginx", "-g", "daemon off;"]

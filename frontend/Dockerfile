# Multi-stage build for Next.js frontend with Bun
FROM oven/bun:latest AS builder

WORKDIR /app/frontend

# Copy dependency files first (layer caching)
COPY frontend/package.json frontend/bun.lock* ./
RUN bun install --frozen-lockfile

# Copy app code and build
COPY frontend .
RUN bun run build

# Runtime stage
FROM oven/bun:latest

WORKDIR /app/frontend

# Copy only built artifacts and node_modules from builder
COPY --from=builder /app/frontend/.next ./.next
COPY --from=builder /app/frontend/public ./public
COPY --from=builder /app/frontend/node_modules ./node_modules
COPY --from=builder /app/frontend/package.json ./

# Create non-root user
# RUN useradd -m -u 1000 appuser && chown -R appuser:appuser /app
# USER appuser

EXPOSE 3000
CMD ["bun", "run", "start"]

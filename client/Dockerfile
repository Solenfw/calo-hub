FROM node:22-alpine AS base
WORKDIR /app

COPY package.json package-lock.json ./
RUN npm ci

# ==========================================
FROM base AS dev
ENV NEXT_TELEMETRY_DISABLED=1
COPY . .
CMD ["npm", "run", "dev", "--", "--hostname", "0.0.0.0"]

# ==========================================
FROM base AS builder
ENV NEXT_TELEMETRY_DISABLED=1
COPY . .
RUN npm run build

# ==========================================
FROM node:22-alpine AS prod
WORKDIR /app
ENV NODE_ENV=production
ENV NEXT_TELEMETRY_DISABLED=1

# 1. Copy the public directory (for images, favicons, etc.)
COPY --from=builder /app/public ./public

# 2. Copy the standalone directory (This contains the minimal required node_modules and server.js)
COPY --from=builder --chown=node:node /app/.next/standalone ./

# 3. Copy the static build files (Standalone mode deliberately leaves these out, so we copy them manually)
COPY --from=builder --chown=node:node /app/.next/static ./.next/static

USER node
EXPOSE 3000

# HOST=0.0.0.0 ensures the standalone server accepts connections from outside the container
ENV HOST=0.0.0.0 

# 4. Instead of 'bun run start', we run the auto-generated standalone server
CMD ["node", "server.js"]
# Simple Dockerfile for the monorepo frontend
FROM node:18-alpine AS base
WORKDIR /app

# Install dependencies
COPY package.json package-lock.json* ./
RUN npm ci --production

# Copy frontend
COPY frontend ./frontend
WORKDIR /app/frontend
RUN npm run build

EXPOSE 3000
CMD ["npm", "run", "start"]

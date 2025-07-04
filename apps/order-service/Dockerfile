# Build stage
FROM node:20-alpine AS builder

WORKDIR /app

# Copy all necessary files
COPY package*.json ./
COPY tsconfig*.json ./
COPY libs/ ./libs/
COPY proto/ ./proto/
COPY apps/order-service/ ./apps/order-service/

# Install all dependencies (including dev dependencies for build)
RUN npm install

# Build the application from the order-service directory
WORKDIR /app/apps/order-service
RUN npm run build

# Production stage
FROM node:20-alpine

WORKDIR /app

# Copy package files and install production dependencies
COPY apps/order-service/package*.json ./
RUN npm install --omit=dev

# Copy built application from builder stage
COPY --from=builder /app/apps/order-service/dist ./dist

# Copy proto files to runtime (ensure they're available at /proto path)
COPY --from=builder /app/proto /proto

# Copy libs if needed at runtime
COPY --from=builder /app/libs ./libs

# Debug: List files to verify structure
RUN ls -la /proto/ || echo "Proto directory not found"
RUN ls -la ./dist/ || echo "Dist directory not found"

EXPOSE 3006

CMD ["node", "dist/main.js"]
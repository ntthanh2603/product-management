# Build stage
FROM node:20-alpine AS builder

WORKDIR /app

COPY ./apps/user-service/package*.json ./
RUN npm install

COPY ./apps/user-service ./
RUN npm run build

# Production stage
FROM node:20-alpine

WORKDIR /app

COPY ./apps/user-service/package*.json ./
RUN npm install --omit=dev

COPY --from=builder /app/dist ./dist

EXPOSE 3008

CMD ["node", "dist/main.js"]
# Build stage
FROM node:20-alpine AS builder

WORKDIR /app

COPY ./tasks-service/package*.json ./
RUN npm install

COPY ./tasks-service ./
RUN npm run build

# Production stage
FROM node:20-alpine

WORKDIR /app

COPY ./tasks-service/package*.json ./
RUN npm install --omit=dev

COPY --from=builder /app/dist ./dist

EXPOSE 3002

CMD ["node", "dist/main.js"]
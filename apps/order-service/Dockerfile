# Build stage
FROM node:20-alpine AS builder

WORKDIR /app

COPY ./order-service/package*.json ./
RUN npm install

COPY ./order-service ./
RUN npm run build

# Production stage
FROM node:20-alpine

WORKDIR /app

COPY ./order-service/package*.json ./
RUN npm install --omit=dev

COPY --from=builder /app/dist ./dist

EXPOSE 3006

CMD ["node", "dist/main.js"]
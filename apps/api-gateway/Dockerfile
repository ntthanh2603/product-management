# Build stage
FROM node:20-alpine AS builder

WORKDIR /app

COPY ./api-gateway/package*.json ./
RUN npm install

COPY ./api-gateway ./
RUN npm run build

# Production stage
FROM node:20-alpine

WORKDIR /app

COPY ./api-gateway/package*.json ./
RUN npm install --omit=dev

COPY --from=builder /app/dist ./dist

EXPOSE 3000

CMD ["node", "dist/main.js"]
# Build stage
FROM node:20-alpine AS builder

WORKDIR /app

COPY ./workspaces/package*.json ./
RUN npm install

COPY ./workspaces ./
RUN npm run build

# Production stage
FROM node:20-alpine

WORKDIR /app

COPY ./workspaces/package*.json ./
RUN npm install --omit=dev

COPY --from=builder /app/dist ./dist

EXPOSE 3001

CMD ["node", "dist/main.js"]
# Stage 1: Frontend Builder
FROM node:20-alpine AS frontend-builder
WORKDIR /app/src/frontend
COPY src/frontend/package.json ./
RUN npm install -g pnpm && pnpm install --frozen-lockfile
COPY src/frontend/ ./
RUN pnpm build

# Stage 2: Backend Builder
FROM golang:1.21-alpine AS backend-builder
WORKDIR /app/src/backend
COPY src/backend/go.mod src/backend/go.sum ./
RUN go mod download
COPY src/backend/ ./
# Frontend dist is built into src/backend/internal/web/dist via Vite outDir
COPY --from=frontend-builder /app/src/backend/internal/web/dist ./internal/web/dist
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o /app/oversteplab ./cmd/server/main.go

# Stage 3: Runtime
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=backend-builder /app/oversteplab .
RUN mkdir -p /app/data
ENV OVERSTEPLAB_DB_PATH=/app/data/oversteplab.db
ENV OVERSTEPLAB_PORT=5000
EXPOSE 5000
ENTRYPOINT ["./oversteplab"]

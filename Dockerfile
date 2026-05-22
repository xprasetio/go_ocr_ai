# Build stage untuk Go backend
FROM golang:1.24.7-alpine AS go-builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build aplikasi
RUN CGO_ENABLED=0 GOOS=linux go build -o ocr-api ./cmd/api

# Build stage untuk frontend Next.js (optional - uncomment jika perlu include frontend)
# FROM node:20-alpine AS web-builder
# WORKDIR /web
# COPY web/package*.json ./
# RUN npm ci
# COPY web/ .
# RUN npm run build

# Final stage - runtime
FROM alpine:3.20

WORKDIR /app

# Install ca-certificates untuk HTTPS
RUN apk add --no-cache ca-certificates

# Copy binary dari builder stage
COPY --from=go-builder /app/ocr-api .

# Expose port default
EXPOSE 8088

# Default environment
ENV PORT=8088

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:${PORT}/health || exit 1

# Run aplikasi
CMD ["./ocr-api"]

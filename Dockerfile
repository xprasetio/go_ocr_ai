# =========================
# Build Stage
# =========================
FROM golang:1.24.7-alpine AS builder

WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -o ocr-api ./cmd/api

# =========================
# Runtime Stage
# =========================
FROM alpine:3.20

WORKDIR /app

# HTTPS certificates
RUN apk add --no-cache ca-certificates

# Copy binary
COPY --from=builder /app/ocr-api .

# Railway inject PORT automatically
EXPOSE 8080

# Run app
CMD ["./ocr-api"]
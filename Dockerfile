# ============================================
# Stage 1: Builder
# ============================================
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy dependency manifests first (layer caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# ============================================
# Stage 2: Final (minimal production image)
# ============================================
FROM alpine:latest

# Install TLS certificates and timezone data
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user for security
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

# Copy compiled binary from builder
COPY --from=builder /app/main .

# Create uploads directory
RUN mkdir -p /var/uploads/notaria178 && chown -R appuser:appgroup /var/uploads/notaria178

# Switch to non-root user
USER appuser

EXPOSE 8080

ENTRYPOINT ["./main"]

# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install git for fetching dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o server ./cmd/server

# Test stage - runs all unit tests with verbose output
FROM builder AS tester
CMD ["go", "test", "-v", "-count=1", "-cover", "./..."]

# Final stage
FROM alpine:3.19

# Install ca-certificates for HTTPS and tzdata for timezones
RUN apk --no-cache add ca-certificates tzdata

# Set timezone
ENV TZ=America/Sao_Paulo

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/server .

# Copy .env.example as default .env
COPY --from=builder /app/.env.example .env

# Create uploads directory
RUN mkdir -p /app/uploads

# Expose port
EXPOSE 8000

# Run the server
CMD ["./server"]

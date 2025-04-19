# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod file
COPY go.mod ./

# Download dependencies (this will create go.sum if needed)
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o validator-key-manager ./cmd/validator-key-manager

# Final stage
FROM alpine:latest

WORKDIR /app

# Create non-root user
RUN adduser -D -g '' appuser

# Copy binary from builder
COPY --from=builder /app/validator-key-manager .

# Use non-root user
USER appuser

# Expose port 8080
EXPOSE 8080

# Run the application
CMD ["./validator-key-manager"]

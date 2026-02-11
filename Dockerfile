# Build stage
FROM golang:1.24-alpine AS builder

# Install build dependencies for SQLite
RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o weather-api .

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache ca-certificates sqlite

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/weather-api .

# Expose port
EXPOSE 3000

# Run the application
CMD ["./weather-api"]
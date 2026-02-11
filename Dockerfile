# Build stage for Go backend
FROM golang:1.24-alpine AS go-builder

RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o weather-api .

# Build stage for frontend using Bun
FROM oven/bun:1-alpine AS frontend-builder

WORKDIR /app

COPY frontend/package.json ./
RUN bun install

COPY frontend/ .
RUN bun run build

# Final stage
FROM alpine:latest

RUN apk add --no-cache ca-certificates sqlite

WORKDIR /root/

# Copy Go binary
COPY --from=go-builder /app/weather-api .

# Copy frontend build
COPY --from=frontend-builder /app/dist ../dist/frontend

EXPOSE 3000

CMD ["./weather-api"]
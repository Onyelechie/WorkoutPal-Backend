# Stage 1: Build
FROM golang:1.25.1-alpine AS builder
WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server ./src/cmd/api/main.go

# Stage 2: Runtime
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app

# Copy binary
COPY --from=builder /app/server .

# Expose port expected by Azure
EXPOSE 8080

# Set environment variable for Azure App Service
ENV PORT=8080

# Start the server
CMD ["./server"]

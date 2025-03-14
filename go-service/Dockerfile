FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the Go service
RUN CGO_ENABLED=0 GOOS=linux go build -o go-service ./go-service

# Start a new stage from scratch
FROM alpine:latest

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/go-service .

# Expose port 50051 for gRPC
EXPOSE 50051

# Run the service
CMD ["./go-service"]
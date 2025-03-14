FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the API Gateway
RUN CGO_ENABLED=0 GOOS=linux go build -o api-gateway ./api-gateway

# Start a new stage from scratch
FROM alpine:latest

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/api-gateway .

# Expose port 8080 for REST API
EXPOSE 8080

# Run the API Gateway
CMD ["./api-gateway"]
# Stage 1: Build the Go app
FROM golang:1.21 AS builder

WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker cache
COPY backend/go.mod backend/go.sum ./

# Download dependencies
RUN go mod tidy

# Copy the rest of the Go source code
COPY backend/. .

# Build a statically linked binary for Alpine
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o main . || { echo 'Build failed'; exit 1; }

# Stage 2: Create a minimal image for running the app
FROM alpine:latest

# Install necessary certificates to make requests
RUN apk --no-cache add ca-certificates

# Set the working directory
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/main /root/

# Expose the application port
EXPOSE 8080

# Run the binary
CMD ["./main"]

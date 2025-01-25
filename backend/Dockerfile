# Stage 1: Build the Go app
FROM golang:1.21 AS builder

WORKDIR /app

# Copy the Go source code
COPY . .

# Download dependencies and build the binary
RUN go mod tidy

# Build for the correct architecture
RUN GOARCH=amd64 go build -o main .

# Stage 2: Create a minimal image for running the app
FROM alpine:latest

RUN apk --no-cache add ca-certificates

# Set the working directory
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/main /root/

# Expose the application port
EXPOSE 8080

# Run the binary
CMD ["./main"]

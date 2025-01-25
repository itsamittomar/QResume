FROM golang:1.20 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod tidy
COPY . .
RUN go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
RUN apk add --no-cache libc6-compat  # Install glibc compatibility

WORKDIR /app
COPY --from=builder /app/main /app/main
RUN chmod +x /app/main

EXPOSE 8080
CMD ["./main"]

# Stage 1: Build
FROM golang:1.22 AS builder

WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary for Linux
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

# Stage 2: Runtime
FROM alpine:3.18

WORKDIR /root/

# Install necessary runtime libraries (if needed)
RUN apk --no-cache add libc6-compat

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Ensure the binary is executable
RUN chmod +x /root/main

# Expose the application port
EXPOSE 8080

# Run the binary
CMD ["./main"]

# ------------------------
# 1) Build stage
# ------------------------
FROM golang:alpine AS builder

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum first (for caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of your backend source code
COPY . .

# Build the Go binary (disable CGO for a static binary if desired)
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./main.go

# ------------------------
# 2) Runtime stage
# ------------------------
FROM alpine:latest

# Create a non-root user (optional, but recommended)
RUN adduser -D webapp
USER webapp

# Copy binary from builder stage
COPY --from=builder /app/main /usr/local/bin/main

# Expose the port your Go app listens on
EXPOSE 8080

# Run the binary
CMD ["main"]

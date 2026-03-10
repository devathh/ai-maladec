# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install dependencies for git and ca-certificates
RUN apk add --no-cache git ca-certificates

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o xcoder ./cmd/main.go

# Final stage
FROM alpine:latest

WORKDIR /root/

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Copy the binary from builder
COPY --from=builder /app/xcoder .
COPY --from=builder /app/start.md .

# Expose port if needed (currently console app, but good practice)
# EXPOSE 8080

CMD ["./xcoder"]

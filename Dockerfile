# syntax=docker/dockerfile:1

########################################
# Build Stage
########################################
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Dependency caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build a statically linked Linux binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app .

########################################
# Final (Runtime) Stage
########################################
FROM alpine:latest

WORKDIR /app

# Install CA certs (optional, needed for HTTPS)
RUN apk --no-cache add ca-certificates

# Copy binary from builder stage
COPY --from=builder /app/app .

# Run as non-root (optional)
USER nobody

ENTRYPOINT ["./app"]

# syntax=docker/dockerfile:1

########################################
# Build Stage
########################################
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Enable cross-platform builds using BuildKit-provided args
ARG TARGETOS
ARG TARGETARCH

# Dependency caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build for the correct target platform
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o app .

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

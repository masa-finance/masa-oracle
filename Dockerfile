# Build the Go binary in a separate stage
FROM golang:1.22 AS builder

# Install Node.js and Yarn
RUN curl -fsSL https://deb.nodesource.com/setup_20.x | bash - && \
    apt-get install -y nodejs && \
    npm install -g npm@latest yarn

WORKDIR /app
# Only copy files needed for go mod download
COPY go.mod go.sum ./
RUN go mod download

# Only copy the necessary source directories
COPY cmd/ cmd/
COPY pkg/ pkg/
COPY internal/ internal/
COPY contracts/ contracts/
COPY Makefile ./

# Set specific GOARCH and ensure we're building with basic CPU features
ENV GOARCH=amd64
ENV CGO_ENABLED=0
ENV GOMAXPROCS=4

# Accept version as build arg
ARG VERSION=dev
# Build with minimal CPU feature set
RUN VERSION=${VERSION} GOAMD64=v1 make build

# Use the official Ubuntu 22.04 image as a base for the final image
FROM ubuntu:22.04 AS base

# Install ca-certificates to ensure TLS verification works
RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates && \
    update-ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/bin/masa-node /usr/bin/masa-node
RUN chmod +x /usr/bin/masa-node

# Create app directory and .masa subdirectory
RUN mkdir -p /app/.masa && \
    # Create masa user and set ownership
    useradd -m -s /bin/bash masa && \
    chown -R masa:masa /app

# Switch to user 'masa' for following commands
USER masa
WORKDIR /app

# Declare the volume for persistence
VOLUME ["/app/.masa"]

# Expose necessary ports
EXPOSE 4001 8080

# Set default command to start the Go application
ENTRYPOINT [ "/usr/bin/masa-node", "--masaDir", "/app/.masa" ]
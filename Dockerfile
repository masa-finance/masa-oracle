# Build the Go binary in a separate stage
FROM golang:1.22 AS builder

# Install Node.js and Yarn
RUN curl -fsSL https://deb.nodesource.com/setup_20.x | bash - && \
    apt-get install -y nodejs && \
    npm install -g npm@latest yarn && \
    yarn config set --home enableTelemetry 0

WORKDIR /app
# Only copy files needed for go mod download
COPY go.mod go.sum ./
RUN go mod download

# Copy all contract files first
COPY contracts/ contracts/
WORKDIR /app/contracts
RUN yarn install --frozen-lockfile --non-interactive --no-git
WORKDIR /app

# Copy remaining source directories
COPY cmd/ cmd/
COPY pkg/ pkg/
COPY internal/ internal/
COPY node/ node/
COPY docs/ docs/
COPY Makefile ./

# Set specific GOARCH and ensure static binary
ENV GOARCH=amd64
ENV CGO_ENABLED=0

# Set VERSION for the build
ARG VERSION=dev
ENV VERSION=${VERSION}

# Build
RUN mkdir -p bin && make build

# Use the official Ubuntu 22.04 image as a base for the final image
FROM ubuntu:22.04 AS base

# Install ca-certificates to ensure TLS verification works
RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates && \
    update-ca-certificates && \
    rm -rf /var/lib/apt/lists/*

# Copy the binary from the builder stage
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
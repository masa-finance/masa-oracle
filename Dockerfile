# Build the Go binary in a separate stage utilizing Makefile
FROM golang:1.22 AS builder

# Install necessary packages for the final image - modified to be more robust
RUN apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends \
    curl \
    sudo \
    gpg \
    lsb-release \
    python3-software-properties \
    software-properties-common \
    git \
    apt-utils \
    && rm -rf /var/lib/apt/lists/* \
    && curl -fsSL https://deb.nodesource.com/setup_lts.x | bash - \
    && apt-get update \
    && DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends nodejs yarn \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Build with version from build arg
ARG VERSION
RUN VERSION=${VERSION:-$(date +%Y%m%d-%H%M%S)} make build || echo "Using fallback version: $(date +%Y%m%d-%H%M%S)" && make build

# Use the official Ubuntu 22.04 image as a base for the final image
FROM ubuntu:22.04 AS base

# Install ca-certificates to ensure TLS verification works
RUN apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends ca-certificates curl && \
    update-ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/bin/masa-node /usr/bin/masa-node
RUN chmod +x /usr/bin/masa-node

# Create the 'masa' user and set up the home directory
RUN useradd -m -s /bin/bash masa && mkdir -p /home/masa/.masa && chown -R masa:masa /home/masa

# Switch to user 'masa' for following commands
USER masa
WORKDIR /home/masa

# Expose necessary ports
EXPOSE 4001 8080

# Add healthcheck
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

# Set default command to start the Go application
ENTRYPOINT [ "/usr/bin/masa-node" ]
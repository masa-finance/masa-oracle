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
    && curl -sS https://dl.yarnpkg.com/debian/pubkey.gpg | gpg --dearmor -o /usr/share/keyrings/yarn-archive-keyring.gpg \
    && echo "deb [signed-by=/usr/share/keyrings/yarn-archive-keyring.gpg] https://dl.yarnpkg.com/debian/ stable main" | tee /etc/apt/sources.list.d/yarn.list \
    && apt-get update \
    && DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends yarn \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy only necessary files for building
COPY go.mod go.sum ./
COPY Makefile ./
COPY cmd/ ./cmd/
COPY pkg/ ./pkg/
COPY internal/ ./internal/
COPY node/ ./node/
COPY config/ ./config/
COPY contracts/ ./contracts/
COPY tools/ ./tools/

# Download dependencies
RUN go mod download

# Build with version from build arg
ARG VERSION
RUN VERSION=${VERSION:-$(date +%Y%m%d-%H%M%S)} make build

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
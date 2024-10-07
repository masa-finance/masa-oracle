# Build stage
FROM golang:1.22 AS builder

WORKDIR /app

# Install Node.js and Yarn
RUN apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    curl gpg \
    && curl -fsSL https://deb.nodesource.com/setup_lts.x | bash - \
    && apt-get install -y nodejs \
    && npm install -g yarn

# Copy the entire project
COPY . .

# Install contract dependencies
RUN cd contracts && yarn install

# Build the Go binary
RUN make build

# Final stage
FROM ubuntu:22.04

# Install necessary packages
RUN apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    nodejs npm \
    && npm install -g yarn

# Create the 'masa' user and set up the home directory
RUN useradd -m -s /bin/bash masa && \
    mkdir -p /home/masa/.masa && \
    chown -R masa:masa /home/masa

# Copy the built binary and set permissions
COPY --from=builder /app/bin/masa-node /usr/bin/masa-node
RUN chmod +x /usr/bin/masa-node

# Copy contracts directory including node_modules
COPY --from=builder --chown=masa:masa /app/contracts /home/masa/contracts

# Switch to user 'masa' for following commands
USER masa
WORKDIR /home/masa

# Copy the .env file
COPY --chown=masa:masa .env .

# Expose necessary ports
EXPOSE 4001 8080

# Set default command to start the MASA node
CMD /usr/bin/masa-node --bootnodes="$BOOTNODES" --env="$ENV" --validator="$VALIDATOR" --cachePath="$CACHE_PATH"
# Build the Go binary in a separate stage utilizing Makefile
FROM golang:1.22 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN make build

# Use the official Ubuntu 22.04 image as a base for the final image
FROM ubuntu:22.04 AS base

# Install necessary packages and N|Solid Runtime
RUN apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    curl sudo gpg lsb-release software-properties-common \
    && curl -fsSL https://deb.nodesource.com/setup_lts.x | bash - \
    && apt-get update && apt-get install -y git apt-utils nsolid -y

COPY --from=builder /app/bin/masa-node /usr/bin/masa-node
RUN chmod +x /usr/bin/masa-node

# Create the 'masa' user and set up the home directory
RUN useradd -m -s /bin/bash masa && mkdir -p /home/masa/.masa && chown -R masa:masa /home/masa

# Copy contracts directory
COPY --chown=masa:masa contracts /home/masa/contracts

# Switch to user 'masa' for following commands
USER masa
WORKDIR /home/masa

# Install contract dependencies
RUN cd /home/masa/contracts && npm install

# Copy the .env file into the container
COPY --chown=masa:masa .env .

# Expose necessary ports
EXPOSE 4001 8080

# Set default command to start the Go application
CMD /usr/bin/masa-node --bootnodes="$BOOTNODES" --env="$ENV" --validator="$VALIDATOR" --cachePath="$CACHE_PATH"
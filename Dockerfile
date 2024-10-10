# Build the Go binary in a separate stage utilizing Makefile
FROM golang:1.22 AS builder

WORKDIR /app

# Install necessary packages for the final image
RUN apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    curl gpg git \
    && curl -fsSL https://deb.nodesource.com/setup_lts.x | bash - \
    && curl -sS https://dl.yarnpkg.com/debian/pubkey.gpg | gpg --dearmor -o /usr/share/keyrings/yarn-archive-keyring.gpg \
    && echo "deb [signed-by=/usr/share/keyrings/yarn-archive-keyring.gpg] https://dl.yarnpkg.com/debian/ stable main" | tee /etc/apt/sources.list.d/yarn.list \
    && apt-get update && apt-get install -y yarn

# Copy the entire project
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Run the go:generate step and build the Go binary
RUN make build

# Use the official Ubuntu 22.04 image as a base for the final image
FROM ubuntu:22.04 AS base

# Copy the built binary and set permissions
COPY --from=builder /app/bin/masa-node /usr/bin/masa-node
RUN chmod +x /usr/bin/masa-node

# Create the 'masa' user and set up the home directory
RUN useradd -m -s /bin/bash masa && \
    mkdir -p /home/masa/.masa && \
    chown -R masa:masa /home/masa

# Switch to user 'masa' for following commands
USER masa
WORKDIR /home/masa

# Copy the .env file into the container
COPY --chown=masa:masa .env .

# Expose necessary ports
EXPOSE 4001 8080

# Set default command to start the Go application
CMD /usr/bin/masa-node --bootnodes="$BOOTNODES" --env="$ENV" --validator="$VALIDATOR" --cachePath="$CACHE_PATH"

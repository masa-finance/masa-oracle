# Use the official Ubuntu 22.04 image as a base for the final image
FROM ubuntu:22.04 AS base

# Install necessary packages for the final image
RUN apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    curl sudo gpg lsb-release software-properties-common \
    && curl -fsSL https://deb.nodesource.com/setup_lts.x | bash - \
    && apt-get install -y apt-utils \
    && apt-get install -y nodejs \
    && npm install -g npm@latest \
    && apt-get update && apt-get install -y git

# Create the 'masa' user and set up the home directory
RUN useradd -m -s /bin/bash masa && mkdir -p /home/masa/.masa && chown -R masa:masa /home/masa


# Switch to user 'masa' for following commands
USER masa
WORKDIR /home/masa

# Copy and install Node.js dependencies for the contracts
# Assuming your contracts directory is ready for copy at this stage
COPY --chown=masa:masa contracts/ ./contracts/
RUN cd contracts && npm install



# Switch back to root to install the Go binary
USER root

# Build the Go binary in a separate stage utilizing Makefile
FROM golang:1.22 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN make build

# Continue with the final image
FROM base

COPY --from=builder /app/bin/masa-node /usr/bin/masa-node

RUN chmod +x /usr/bin/masa-node

# Switch to 'masa' to run the application
USER masa
WORKDIR /home/masa

# Copy the .env file into the container
COPY --chown=masa:masa .env .

# Expose necessary ports
EXPOSE 4001 8080

# Set default command to start the Go application

CMD /usr/bin/masa-node --bootnodes="$BOOTNODES" --env="$ENV" --validator="$VALIDATOR"
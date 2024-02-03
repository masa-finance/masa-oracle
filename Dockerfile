# Use the official Ubuntu 22.04 image as a base for the final image
FROM ubuntu:22.04 as base

# Install necessary packages for the final image
RUN apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y curl sudo gpg lsb-release software-properties-common

# Node.js and Yarn setup for the final image
RUN curl -fsSL https://deb.nodesource.com/setup_lts.x | bash - && \
    curl -sS https://dl.yarnpkg.com/debian/pubkey.gpg | gpg --dearmor -o /usr/share/keyrings/yarn-archive-keyring.gpg && \
    echo "deb [signed-by=/usr/share/keyrings/yarn-archive-keyring.gpg] https://dl.yarnpkg.com/debian/ stable main" | tee /etc/apt/sources.list.d/yarn.list && \
    apt-get update && apt-get install -y nodejs yarn jq

# Install global npm to match version used in script for the final image
RUN npm install -g npm@10.4.0

# Create the 'masa' user and set up the home directory for the final image
RUN useradd -m masa

# Build the Go binary in a separate stage
FROM golang:1.21 as builder

# Set the Current Working Directory inside the container to something other than /go
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Install go dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -v -o masa-node ./cmd/masa-node

# Continue with the final image
FROM base

# Copy the pre-built binary file from the previous stage
COPY --from=builder /app/masa-node /usr/bin/masa-node

# Set execute permissions on the masa-node binary
RUN chmod +x /usr/bin/masa-node

# Set WORKDIR to /home/masa in the final image
WORKDIR /home/masa

# Switch to user 'masa' in the final image
USER masa

# Add the RPC_URL to masa user's .bash_profile (placeholder, update as necessary)
RUN echo "export RPC_URL=\${RPC_URL}" >> .bash_profile

# Expose necessary ports
EXPOSE 4001
EXPOSE 8080

# Set default command (adjust MASANODE_CMD based on your setup)
CMD /usr/bin/masa-node --port=4001 --udp=true --tcp=false --start --bootnodes=${BOOTNODES}


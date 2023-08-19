## DOCKER.md

### Running Masa Oracle Nodes in Docker

This guide explains how to run and connect two Masa Oracle nodes within Docker containers on the same machine.

#### Building the Docker Image

First, ensure you have a `Dockerfile` that sets up the environment for your Go application. Here's a basic example:

```Dockerfile
# Use the official Golang image as our base image
FROM golang:1.20

# Install git (required for fetching dependencies)
RUN apt-get update && apt-get install -y git && rm -rf /var/lib/apt/lists/*

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy everything from the current directory to the Working Directory inside the container
COPY . .

# Fetch dependencies
RUN go mod tidy

# Build the Go app
RUN go build -o main .

# Expose port 4001 (change if necessary)
EXPOSE 4001

# Command to run the executable
CMD ["./main"]
```

Build the Docker image:

```bash
docker build -t masa-node .
```

#### Running the First Node

When running the first node, give it a specific name for easier reference:

```bash
docker run --network=dev -p 4001:4001 --name masa-node1 -it masa-node
```

After running, note the libp2p host address from the logs, which will be required to connect the second node to it.

#### Getting the IP of the First Container

To get the IP of the first container within the Docker network, use:

```bash
docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' masa-node1
```

This might give you an IP like `172.19.0.2`.

#### Running the Second Node

To run the second node and connect it to the first, use the IP obtained above:

```bash
docker run --network=dev --name masa-node2 -it masa-node ./main /ip4/172.19.0.2/tcp/PORT_OF_FIRST_NODE/p2p/ID_OF_FIRST_NODE
```

Replace `PORT_OF_FIRST_NODE` and `ID_OF_FIRST_NODE` with the appropriate values from the first node's libp2p address.

By following these steps, the second container will connect to the first using the Docker network IP, allowing them to communicate. You can now interact with each node individually by attaching to their containers.
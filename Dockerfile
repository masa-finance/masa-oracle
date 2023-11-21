# Use the official Golang image as our base image
FROM golang:1.21

# Install git (required for fetching dependencies)
RUN apt-get update && apt-get install -y git && rm -rf /var/lib/apt/lists/*

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy everything from the current directory to the Working Directory inside the container
COPY . .

# Fetch dependencies
RUN go mod tidy

# Build the Go app
RUN go build -v -o masa-node ./cmd/masa-node

# Expose port 4001 (change if necessary)
EXPOSE 4001

# Command to run the executable with bootnode address
CMD ["./masa-node", "--bootnodes=/ip4/137.66.11.250/udp/4001/quic-v1/p2p/16Uiu2HAmJtYy4A8pzChDQQLrPsu1SQU5apzCftCVaAjFk539CLc9,/ip4/168.220.95.86/udp/4001/quic-v1/p2p/16Uiu2HAmNr4aSznspSGfqMtTBDauiyvRqArvEYRPUQArpCFtdy2P" "-port=4001", "--udp=true"]

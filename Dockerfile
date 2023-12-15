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
CMD ["./masa-node", "--bootnodes=/ip4/35.224.231.145/udp/4001/quic-v1/p2p/16Uiu2HAm47nBiewWLLzCREtY8vwPQtr5jTqyrEoUo6WnngwhsQuR,/ip4/104.198.43.138/udp/4001/quic-v1/p2p/16Uiu2HAkxiP8jjdHQWeCxTr7pD6BvoPkS8Z1skjCy9vdSRMACDcc,/ip4/107.223.13.174/udp/5001/quic-v1/p2p/16Uiu2HAmMkXJJpPAdEmp9QSqdcTPzvV2UxvZMEhYdVLFzbQHHczp,/ip4/35.202.227.74/udp/4001/quic-v1/p2p/16Uiu2HAmHuUejpUBFPCxy32QhGRAbv3tFwbzXmLkCoaNcZTyWWqN,/ip4/93.187.217.133/udp/4001/quic-v1/p2p/16Uiu2HAm5wvEfWGufJ1roGL6VhpFZ4scqPF1giLwES9jXfeEoeHs,/ip4/10.128.0.47/udp/4001/quic-v1/p2p/16Uiu2HAkxiP8jjdHQWeCxTr7pD6BvoPkS8Z1skjCy9vdSRMACDcc,/ip4/147.75.56.191/udp/4001/quic-v1/p2p/16Uiu2HAmVrXpTot74CFpdFNpTs26QminLwXT3HhXPSc1MFjnqqSR,/ip4/107.223.13.174/udp/4001/quic-v1/p2p/16Uiu2HAm2uQ5TGviRkqhYMpg7fjeoB4TfpSAhrbY87YZ4h9jYCNm,/ip4/34.171.201.124/udp/4001/quic-v1/p2p/16Uiu2HAmCKzfsynicpryPZTdcJsjmyzXn8tA13zMHHsoBxLdvVCE,/ip4/34.132.48.64/udp/4001/quic-v1/p2p/16Uiu2HAmNk4DDNiVu8ipN2cg5GLpGzN6ydd4EYps1NkiTDBRkctu", "--port=4001", "--udp=true", "--start=true"]
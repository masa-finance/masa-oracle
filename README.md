# Masa Oracle: Decentralized Data Protocol

Masa Oracle is a pioneering protocol designed to revolutionize the way data behavioral, and identity data is accessed, distributed, and incentivized in a decentralized manner. By leveraging the power of blockchain technology, the Masa Oracle ensures transparency, security, and fair rewards for nodes participating in the data distribution network.

## Getting Started

### Prerequisites

Ensure you have Go installed on your system. If not, you can download and install it from [here](https://golang.org/dl/).

### Running the Node

1. Clone the repository:

```bash
git clone https://github.com/masa-finance/masa-oracle.git
cd masa-oracle
```

2. Build the node and put the binary in the bin directory:

```bash
go build -v -o bin/masa-node ./cmd/masa-node
```

3. Run the node:

```bash
bin/masa-node   
```

When the node is started for the first time, it will generate a new key pair and store it in the `~/.masa-node/masa_oracle_node.env` file. The node will use this key pair for all future runs.

You may update the .env file to specify a specific port number if you wish to do so.
```
portNbr=4001
```
You should see the node's address printed on the console. This indicates that your node is up and running, ready to connect with other Masa nodes.

---

## Connecting Nodes

Once you have the Masa node set up, you can easily connect multiple nodes together. Here's a step-by-step guide on how to do this:


## Running additional Nodes with Docker
To run additional nodes on the same machine, you will need to use Docker
First make sure to have the multi address for your running node as described above.
You will need to update the "Dockerfile" with the multi address of your running node.
```
 # Command to run the executable with bootnode address
 CMD ["./masa-oracle", "/ip4/192.168.1.6/tcp/4001/ws/p2p/QmQq37unSom5Vv2dzEyiRqPc8V9JUAZXFmkyQfZtW4J1Bt"]
```

### Prerequisites

Ensure you have Docker installed on your system. If not, you can download and install it from [here](https://docs.docker.com/get-docker/).

### Building the Docker Image

1. Navigate to the project directory:
```bash
   cd path/to/masa-oracle
```

2. Build the Docker image:
```bash
   docker build -t masa-node .
```

This command builds a Docker image using the Dockerfile in the current directory and tags it as `masa-node`.

### Running the Docker Container

Run the Docker container with the following command:

```bash
docker run -p 4001:4001 masa-node
```

This address is the multiaddress of the node. It provides all the necessary information for another peer to locate and communicate with this node.

## Contribution

Contributions are always welcome. Please fork the repository and create a pull request with your changes. Ensure that your code follows Go best practices.

## License

This project is licensed under the terms of the [MIT license](LICENSE).

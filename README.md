# Masa Oracle: Decentralized Data Protocol

Masa Oracle is a pioneering protocol designed to revolutionize the way data behavioral, and identity data is accessed, distributed, and incentivized in a decentralized manner. By leveraging the power of blockchain technology, the Masa Oracle ensures transparency, security, and fair rewards for nodes participating in the data distribution network.

## Getting Started

### Prerequisites

Before you begin, make sure you have the following prerequisites installed on your system:

- **Go**: If not already installed, download and install it from [here](https://golang.org/dl/).
- **Docker**: If not already installed, download and install it from [here](https://docs.docker.com/get-docker/).


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
### Run the node:

```bash
bin/masa-node /ip4/34.133.16.77/udp/4001/quic-v1/p2p/16Uiu2HAmAEDCYv5RrbLhZRmHXGWXNuSFa7YDoC5BGeN3NtDmiZEb
```

This will start a new Masa Oracle node, and it will use the multi-address for connecting to the Masa Oracle node with the IP address 34.133.16.77.

You are now ready to connect your Masa node with the specified node in the network. Be sure to follow any additional configuration steps and best practices specific to your use case or network requirements.

Remember to check the Masa Oracle repository for any updates or additional information on using the protocol.

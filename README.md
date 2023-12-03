# Masa Oracle: Decentralized Data Protocol

The Masa Oracle redefines the way behavioral and identity data is managed, shared, and monetized in a decentralized manner. With a focus on transparency, security, and fair compensation, Masa Oracle empowers nodes within the zk-Data Network & Marketplace to operate with integrity and trust.

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
go build -v -o masa-node ./cmd/masa-node
```

3. Run the node:

```bash
bin/masa-node   
```

When the node is started for the first time, it will generate a new key pair and store it in the `~/.masa-node/masa_oracle_node.env` file. The node will use this key pair for all future runs.

## Running the Program

You can run the program with various flags to customize its behavior. Here's how you can specify the flags:

```go run main.go --bootnodes=node1,node2,node3 --port=8080 --udp=true --tcp=false```

In this command:

- `--bootnodes=node1,node2,node3` sets the `bootnodes` argument to `"node1,node2,node3"`.
- `--port=8080` sets the `port` argument to `8080`.
- `--udp=true` sets the `udp` argument to `true`.
- `--tcp=false` sets the `tcp` argument to `false`.

If an argument is not specified in the command line, its default value will be used. The default values are:

- `bootnodes`: The value of the `BOOTNODES` environment variable. If the environment variable is not set, the default value is an empty string.
- `port`: The value of the `portNbr` environment variable. If the environment variable is not set or is not a valid integer, the default value is `0`.
- `udp`: The value of the `UDP` environment variable. If the environment variable is not set or is not a valid boolean, the default value is `false`.
- `tcp`: The value of the `TCP` environment variable. If the environment variable is not set or is not a valid boolean, the default value is `false`.

If neither `udp` nor `tcp` are set, `udp` will default to `true`.
`bootnodes` is a comma-separated list of multi-addresses of the nodes to connect to. If `bootnodes` is not set, the node will not connect to any other nodes and behave as the main bootnode.
a multiaddress looks like: `/ip4/10.0.0.18/tcp/4001/p2p/16Uiu2HAm2uQ5TGviRkqhYMpg7fjeoB4TfpSAhrbY87YZ4h9jYCNm`

You may update the .env file to specify a specific port number if you wish to do so.
this is located in:

`$HOME/.masa-node/masa_oracle_node.env`

```
portNbr=4001
bootnodes=/ip4/10.0.0.18/tcp/4001/p2p/16Uiu2HAm2uQ5TGviRkqhYMpg7fjeoB4TfpSAhrbY87YZ4h9jYCNm,
udp=true
tcp=false
```
You should see the node's address printed on the console. This indicates that your node is up and running, ready to connect with other Masa nodes.

---

## Connecting Nodes
### Run the node:

```bash
bin/masa-node --bootnodes=/ip4/34.133.16.77/udp/4001/quic-v1/p2p/16Uiu2HAmAEDCYv5RrbLhZRmHXGWXNuSFa7YDoC5BGeN3NtDmiZEb --port=4001 --udp=true --tcp=false
```

This will start a new Masa Oracle node, and it will use the multi-address for connecting to the Masa Oracle node with the IP address 34.133.16.77.

You are now ready to connect your Masa node with the specified node in the network. Be sure to follow any additional configuration steps and best practices specific to your use case or network requirements.

Remember to check the Masa Oracle repository for any updates or additional information on using the protocol.

# Masa Oracle: Decentralized Data Protocol 🌐

The Masa Oracle defines how private behavioral and identity data is accessed, shared, and rewarded in a decentralized and private way. The Masa Oracle guarantees transparency, security, and equitable compensation for nodes that particiapte in the Masa zk-Data Network & Marketplace.

## Contents
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
- [Staking Tokens](#staking-tokens)
- [Running the Node](#running-the-node)
- [Command-Line Interface (CLI)](#command-line-interface-cli)
- [Configuration](#configuration)
- [Connecting Nodes](#connecting-nodes)
- [Ad Network Proof of Concept](#ad-network-proof-of-concept)
- [Use Cases](#use-cases)
- [Updates & Additional Information](#updates--additional-information)

## Getting Started

### Prerequisites

Before diving in, ensure these prerequisites are installed:
- **Go**: Grab it from [Go's official site](https://golang.org/dl/).
- **Docker**: Install from [Docker's official docs](https://docs.docker.com/get-docker/).

### Installation

1. Clone the repository:
```bash
git clone https://github.com/masa-finance/masa-oracle.git
cd masa-oracle
```

2. Build the node executable:
```bash
go build -v -o masa-node ./cmd/masa-node
```

## Staking Tokens

🔐 To participate in the network and earn rewards, you must first stake your tokens:
```bash
./masa-node --stake 100
```
This command initiates the staking process, allowing the staking contract to spend tokens on your behalf and then stakes the specified amount.

## Running the Node

🚀 To start your node and join the Masa network:
```bash
./masa-node --start
```
This command will start the node using the list of bootnodes specified in the configuration file, if available.

## Command-Line Interface (CLI)

Customize your node's behavior with various flags:
```bash
./masa-node --bootnodes=node1,node2,node3 --port=8080 --udp=true --tcp=false
```
- `--bootnodes=node1,node2,node3`: Connect to specified bootnodes.
- `--port=8080`: Listen on port `8080`.
- `--udp=true`: Enable UDP protocol.
- `--tcp=false`: Disable TCP protocol.

Defaults are used if flags are not set:

- `bootnodes`: Falls back to `BOOTNODES` env variable or an empty string.
- `port`: Defaults to `portNbr` env variable or `0`.
- `udp`: Defaults to `UDP` env variable or `false`.
- `tcp`: Defaults to `TCP` env variable or `false`.

## Configuration

🔧 To use a custom configuration file:

```bash
./masa-node --config=path/to/config.json
```

The configuration file is a JSON format that includes an array of bootnodes.

## Connecting Nodes

🔗 To connect to a specific node in the network:
```bash
./masa-node --bootnodes=/ip4/34.133.16.77/udp/4001/quic-v1/p2p/16Uiu2HAmAEDCYv5RrbLhZRmHXGWXNuSFa7YDoC5BGeN3NtDmiZEb --port=4001 --udp=true --tcp=false
```
This will connect your Masa Oracle node to the specified node using the provided multi-address.

## Ad Network Proof of Concept

📢 Masa Oracle introduces a decentralized ad network as a proof of concept within the decentralized data protocol. This network allows publishers who are staked in the system to publish advertisements to a dedicated topic.

### How It Works

- Publishers must first stake tokens to participate in the ad network, ensuring a commitment to the network's integrity.
- Once staked, publishers can publish ads to the `ad-topic` using the `PublishAd` method in the `OracleNode`.
- Ads are structured with content and metadata, as defined in the `Ad` struct in `pkg/ad/ad.go`.
- Only staked nodes can publish to the ad topic, as enforced by the `PublishAd` method in `OracleNode`, which checks the `IsStaked` flag before allowing publication.

## Use Cases

The Masa Oracle ad network can be utilized in various scenarios, such as:

- **Targeted Advertising**: Leveraging the decentralized identity data to deliver personalized ads without compromising user privacy.
- **Content Monetization**: Content creators can receive compensation directly through the protocol for hosting ads.
- **Community Governance**: Staked nodes can vote on ad policies, ensuring that the network remains aligned with the community's values.

## Updates & Additional Information

📢 Stay tuned to the Masa Oracle repository for updates and more details on how to use the protocol effectively.

---

After setting up, your node's address will be displayed, indicating it's ready to connect with other Masa nodes. Follow any additional configuration steps and best practices as per your use case or network requirements.
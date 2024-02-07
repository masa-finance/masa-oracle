# Masa Oracle: Decentralized Data Protocol 🌐

The Masa Oracle governs the access, sharing, and rewarding of private behavioral and identity data in a decentralized and private manner. The Masa Oracle Network ensures transparency and security of data sharing, while  enabling equitable compensation for nodes that participate in the Masa zk-Data Network and Marketplace.

## Contents
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
  - [Docker Setup](#docker-setup)
- [Staking Tokens](#staking-tokens)
- [Running the Node](#running-the-node)
- [Updates & Additional Information](#updates--additional-information)

## Getting Started

### Prerequisites

Ensure these prerequisites are installed for a local setup:
- **Go**: Grab it from [Go's official site](https://golang.org/dl/).
- **Yarn**: Install it via [Yarn's official site](https://classic.yarnpkg.com/en/docs/install/).
- **Git**: Required for cloning the repository.

### Installation


#### Docker Setup

For complete instructions on building, staking, and running a node with Docker, please see [here](./DOCKER.md) 

#### Local Setup

1. Clone the repository
```
git clone https://github.com/masa-finance/masa-oracle.git
```
2. Build the go code into the masa-node binary:
```
go build -v -o masa-node ./cmd/masa-node
```
3. Go into the contracts directory and build the contract npm modules that the go binary uses:
```
cd contracts/ 
npm install
cd ../
```
4. Start up masa-node. Later you'll want to set masa-node up as a service and export the RPC_URL and BOOTNODES you want to use in the environment your service runs in, but for now, you can set them in the command line to start the service up:

RPC_URL=https://ethereum-sepolia.publicnode.com masa-node masa-node --start
   ```

## Funding the Node (in order to Stake)

Find the public key of your node in the logs. 

Send 1000 MASA and .01 sepoliaETH to the node's public key / wallet address.

When the transactions have settled, you can stake

### Staking Tokens

- For local setup, stake tokens with:
  ```bash
  ./masa-node --stake 1000
  ```
- For Docker setup, stake tokens with:
  ```bash
  docker-compose run --rm masa-node /usr/bin/masa-node --stake 1000
  ```

### Running the Node

- **Local Setup**: Connect your node to the Masa network:
  ```bash
  ./masa-node --bootnodes=/ip4/35.223.224.220/udp/4001/quic-v1/p2p/16Uiu2HAmPxXXjR1XJEwckh6q1UStheMmGaGe8fyXdeRs3SejadSa --port=4001 --udp=true --tcp=false --start=true
  ```
- **Docker Setup**: Your node will start automatically with `docker-compose up -d`. Verify it's running correctly:
  ```bash
  docker-compose logs -f masa-node
  ```

After setting up your node, its address will be displayed, indicating it's ready to connect with other Masa nodes. Follow any additional configuration steps and best practices as per your use case or network requirements.

## Updates & Additional Information

Stay tuned to the Masa Oracle repository for updates and additional details on effectively using the protocol. For Docker users, update your node by pulling the latest changes from the Git repository, then rebuild and restart your Docker containers.


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

- For complete instructions on building, staking, and running a node with Docker, please see [here](./DOCKER.md)

### Installation

#### Docker Setup

For complete instructions on building, staking, and running a node with Docker, please see [here](./DOCKER.md)

#### Local Setup

##### 1. Clone the repository

```shell
git clone https://github.com/masa-finance/masa-oracle.git
```

##### 2. Build the go code into the masa-node binary

```shell
go build -v -o masa-node ./cmd/masa-node
```

##### 3. Go into the contracts directory and build the contract npm modules that the go binary uses

```shell
cd contracts/ 
npm install
cd ../
```

##### 4. Set env vars

Ensure your environment has the required env var, ENV, exported. It should be set to test to join the testnet.
Optionally, you can set RPC_URL to change it from the default (<https://ethereum-sepolia.publicnode.com>)

```shell
export ENV=test
export RPC_URL=https://1rpc.io/sepolia # This is optional and will be set to the default https://ethereum-sepolia.publicnode.com without it.
```

##### 5. Start up masa-node. Be sure to include your bootnodes list with the --bootnodes flag

```shell
/ip4/35.223.224.220/udp/4001/quic-v1/p2p/16Uiu2HAmPxXXjR1XJEwckh6q1UStheMmGaGe8fyXdeRs3SejadSa
```

```shell
./masa-node --bootnodes=/ip4/35.223.224.220/udp/4001/quic-v1/p2p/16Uiu2HAmPxXXjR1XJEwckh6q1UStheMmGaGe8fyXdeRs3SejadSa
```

## Makefile Commands

The Makefile provides several commands to build, install, run, test, and clean the Masa Node project. Here's a description of each command:

### make build

The build command compiles the Masa Node binary and places it in the ./bin directory. It uses the go build command with the following flags:

-v: Enables verbose output to show the packages being compiled.
-o ./bin/masa-node: Specifies the output binary name and location.
./cmd/masa-node: Specifies the package to build (the main package).
make install
The install command runs the node_install.sh script to install any necessary dependencies or perform additional setup steps required by the Masa Node.

### make run

The run command first builds the Masa Node binary using the build command and then executes the binary located at ./bin/masa-node. This command allows you to compile and run the Masa Node in a single step.

### make test

The test command runs all the tests in the project using the go test command. It recursively searches for test files in all subdirectories and runs them.

### make clean

The clean command performs cleanup tasks for the project. It removes the bin directory, which contains the compiled binary, and deletes the masa_oracle_node.log file, which may contain log output from previous runs.

To execute any of these commands, simply run make in your terminal from the project's root directory. For example, make build will compile the Masa Node binary, make test will run the tests, and make clean will remove the binary and log file.

## Funding the Node (in order to Stake)

Find the public key of your node in the logs.

Send 1000 MASA and .01 sepoliaETH to the node's public key / wallet address.

When the transactions have settled, you can stake

### Staking Tokens

- For local setup, stake tokens with:

  ```shell
  ./masa-node --stake 1000
  ```

- For Docker setup, stake tokens with:
  
  ```shell
  docker-compose run --rm masa-node /usr/bin/masa-node --stake 1000
  ```

### Running the Node

- **Local Setup**: Connect your node to the Masa network:
  
  ```shell
  ./masa-node --bootnodes=/ip4/35.223.224.220/udp/4001/quic-v1/p2p/16Uiu2HAmPxXXjR1XJEwckh6q1UStheMmGaGe8fyXdeRs3SejadSa --port=4001 --udp=true --tcp=false --start=true --env=test
  ```

- **Docker Setup**: Your node will start automatically with `docker-compose up -d`. Verify it's running correctly:
  
  ```shell
  docker-compose logs -f masa-node
  ```

After setting up your node, its address will be displayed, indicating it's ready to connect with other Masa nodes. Follow any additional configuration steps and best practices as per your use case or network requirements.

## Updates & Additional Information

Stay tuned to the Masa Oracle repository for updates and additional details on effectively using the protocol. For Docker users, update your node by pulling the latest changes from the Git repository, then rebuild and restart your Docker containers.

## Masa Node CLI

For more detailed documentation, please refer to the [CLI.md](docs/CLI.md) file.

## Masa Node Twitter Sentiment Analysis

For more detailed documentation, please refer to the [LLM.md](docs/LLM.md) file.

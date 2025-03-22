# MASA Node Docker Setup Guide

Welcome to the MASA Node Docker setup guide. This document will walk you through the process of setting up and running your own MASA node in a Docker environment. Follow these steps to get up and running quickly.

## Prerequisites

Before you begin, ensure you have the following installed on your system:

- **Docker**: You'll need Docker to build and run containers. Download and install Docker for your operating system from [Docker's official website](https://www.docker.com/products/docker-desktop).
- **Docker Compose**: This project uses Docker Compose to manage multi-container Docker applications. Docker Desktop for Windows and Mac includes Docker Compose. On Linux, you may need to install it separately following the instructions [here](https://docs.docker.com/compose/install/).
- **Git**

## Getting Started

### 1. Clone the Repository

Start by cloning the masa-node repository to your local machine. Open a terminal and run:

```bash
git clone git@github.com:masa-finance/masa-oracle.git
cd masa-oracle
```

### 2. Environment Configuration

Create a `.env` file in the root of your project directory. This file will store environment variables required by the MASA node, such as `BOOTNODES` and `RPC_URL`. You can obtain these values from the project maintainers or documentation.

Example `.env` file content:

```env
BOOTNODES=/dns4/boot-1.test.miners.masa.ai/udp/4001/quic-v1/p2p/16Uiu2HAm9Nkz9kEMnL1YqPTtXZHQZ1E9rhquwSqKNsUViqTojLZt,/dns4/boot-2.test.miners.masa.ai/udp/4001/quic-v1/p2p/16Uiu2HAm7KfNcv3QBPRjANctYjcDnUvcog26QeJnhDN9nazHz9Wi,/dns4/boot-3.test.miners.masa.ai/udp/4001/quic-v1/p2p/16Uiu2HAmBcNRvvXMxyj45fCMAmTKD4bkXu92Wtv4hpzRiTQNLTsL
RPC_URL=https://ethereum-sepolia.publicnode.com
ENV=test
TWITTER_USERNAME="your_username"
TWITTER_PASSWORD="your_password"
TWITTER_2FA_CODE="your_2fa_code"
TWITTER_SCRAPER=True

```

*be sure to use ENV=test to join the masa oracle testnet

### 3. Building the Docker Image

With Docker and Docker Compose installed and your `.env` file configured, build the Docker image using the following command:

```bash
docker-compose build
```

This command builds the Docker image based on the instructions in the provided `Dockerfile` and `docker-compose.yaml`.

### 4. Running the MASA Node

To start the MASA node, use Docker Compose:

```bash
docker-compose up -d
```

This command starts the MASA node in a detached mode, allowing it to run in the background.

### 5. Verifying the Node

After starting the node, you can verify it's running correctly by checking the logs:

```bash
docker-compose logs -f masa-node
```

This command displays the logs of the MASA node container. Look for any error messages or confirmations that the node is running properly.

## Accessing Generated Keys

The MASA node generates keys that are stored in the `.masa-keys/` directory in your project directory.
This directory is mapped from `/home/masa/.masa/` inside the Docker container, ensuring that your keys are safely stored on your host machine.

## Funding the Node (in order to Stake)

### Step 1: Find your Node's public key

The public key of your new node is shown in the output at the beginning of the logs when it starts up:

```bash
docker-compose logs -f masa-node
```

### Step 2: Send sepolia ETH and MASA to your node's public key address

You can obtain sepolia ETH from a faucet such as:

Like so:

You can obtain test (sepolia) MASA from our faucet:

Like so:

### Step 3: Stake the Node

Once the transactions settle, you can stake your node

```bash
docker-compose run --build --rm masa-node /usr/bin/masa-node --stake 1000
```

### Step 4: Restart your Node

Stop your running daemonized node:

```bash
docker compose down
```

Start it up again with the -d flag: (If you have changed settings you may wish to `--force-recreate`)

```bash
docker compose up --build --force-recreate -d
```

## Updating the Node

To update your node, pull the latest changes from the Git repository (if applicable), then rebuild and restart your Docker containers:

```bash
git pull
docker-compose build
docker-compose down
docker-compose up -d
```

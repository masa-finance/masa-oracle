Below is a comprehensive README guide designed to help new users set up and run their own MASA node using Docker. This guide includes step-by-step instructions and covers any prerequisites needed.

---

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
BOOTNODES=<bootnodes-value>
RPC_URL=<rpc-url-value>
```

Replace `<bootnodes-value>` and `<rpc-url-value>` with the actual values.

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

The MASA node generates keys that are stored in the `.masa-keys/` directory in your project directory. This directory is mapped from `/home/masa/.masa/` inside the Docker container, ensuring that your keys are safely stored on your host machine.

## Updating the Node

To update your node, pull the latest changes from the Git repository (if applicable), then rebuild and restart your Docker containers:

```bash
git pull
docker-compose build
docker-compose down
docker-compose up -d
```

## Troubleshooting

If you encounter any issues during setup, ensure you have followed all steps correctly and check the Docker and Docker Compose documentation for additional help. For specific issues related to the MASA node, consult the project's support channels or documentation.

---

This README guide provides a comprehensive overview for new users to set up and run their MASA node using Docker. It covers prerequisites, environment configuration, building, running, and updating the node, as well as accessing generated keys and troubleshooting common issues.

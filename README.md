# [Welcome to Masa](https://docs.masa.ai/docs/welcome-to-masa)

## [Introduction](https://docs.masa.ai/docs/welcome-to-masa#introduction)

Welcome to Masa, the network powering Fair AI. We are on a mission to revolutionize the AI landscape by providing an open, permissionless marketplace for specialized AI training data and compute resources. Our vision is to enable any builder, anywhere in the world, to access the tools they need to create innovative and specialized AI applications.

- **[Masa Protocol](https://docs.masa.ai/docs/protocol/welcome):** Learn about the Oracle Nodes, Worker Nodes, and other components that make up the Masa Protocol.
- **[Masa Bittensor Subnet](https://docs.masa.ai/docs/masa-subnet/welcome):** Discover how the Masa Bittensor Subnet operates, focusing on Twitter data and validator-miner interactions.

## [The Need for Fair AI](https://docs.masa.ai/docs/welcome-to-masa#the-need-for-fair-ai)

Nodes can run on any hardware for which you can build a golang application.

- LLMs are often unfair to people, failing to properly attribute outputs to data sources. This has resulted in lawsuits and risks for companies using these models.
- As general purpose LLMs become commoditized, the most valuable AI applications will be the most specialized ones. However, building specialized AI requires access to specialized training data, which is not readily available.
- Leading compute and inference providers charge high prices and can arbitrarily deny service to builders, stifling innovation and limiting access.

While a Masa Protocol node itself requires few resources to run on testnet, if you wish to create a worker node that performs a useful task, such as running an LLM model, your hardware choices should be dictated by the requirements of that task.

## [The Masa Solution](https://docs.masa.ai/docs/welcome-to-masa#the-masa-solution)

Masa has built the leading marketplace for data and compute, connecting data and compute contributors with developers. Our platform incentivizes people to contribute specialized data sets and sell compute resources in an open, permissionless manner.

Key features of the Masa network include:

- port 8080 is only required to provide access to the API, and can be changed with environment configuration.
- Only 4001 is required to be open publicly for participation in the p2p Masa Protocol network as a Worker node.
- A basic node will still find the bootnodes and register itself as part of the network without any specific inbound ports open.

- **Specialized Data & Open Source LLMs**: Masa enables the contribution of specialized data sets to train AI models and access to open source LLMs provisioned by workers. This empowers builders to create the most valuable and specialized AI applications by leveraging unique data and powerful open source LLMs, all within a decentralized, permissionless ecosystem.

- **Open Marketplace**: Our open, permissionless marketplace democratizes access to AI training data and compute resources. Contributors can earn rewards by contributing data and selling compute, while developers can access the resources they need to build innovative applications.

## [AI Worker Nodes Introduction](https://docs.masa.ai/docs/worker-node/introduction)

Masa empowers contributors to monetize their data and compute resources by becoming Worker Nodes on the network. As a Worker Node, you can:

- Stake tokens to provide work to the network and earn rewards for your nodes availablity and uptime.
- Receive and process data requests from Oracle Nodes, servicing data and LLM requests.
- Earn network emissions and fees for the work you provide, generating a new revenue stream.
- Access the network's vast dataset and LLM resources to efficiently process requests and deliver accurate results.

By contributing to Masa as a Worker Node, you play a vital role in powering the decentralized AI ecosystem while being rewarded for your efforts.

## [Masa for Developers: AI Oracle Nodes](https://docs.masa.ai/docs/oracle-node/introduction)

Masa offers developers a decentralized platform to access diverse data sources and powerful LLM services through Oracle Nodes. By running an Oracle Node, you can:

- Stake tokens to access the network's rich data and powerful open source LLM services, submitting data or LLM requests to Worker Nodes.
- Submit data or LLM requests to Worker Nodes to power your AI applications.
- Tap into a wide range of data sources and LLM models to fulfill various data and processing requirements, for example:
  - **Crypto Sentiment Analysis**: Combine data from our Twitter Scraper and Web Scraper to gather real-time information about cryptocurrency trends, news, and public sentiment.
  - **Crypto Community Insights**: Leverage our Discord Profile scraper to extract comprehensive data from prominent crypto users on Discord.
  - **Crypto News Aggregation and Summarization**: Utilize our Web Scraper to collect real-time data from leading crypto news websites and blogs.

By leveraging Masa as a Oracle Node, developers can build innovative AI applications with the power of decentralized data and compute at their fingertips.

## [Join the Fair AI Revolution](https://docs.masa.ai/docs/welcome-to-masa#join-the-fair-ai-revolution)

Masa is more than just a technology platform - it's a movement to make AI more accessible, equitable, and beneficial for all. By contributing data or compute resources to the Masa network, you can help power the next generation of Fair AI applications and be rewarded for your contributions.

We invite you to join us in building a decentralized future for AI. Explore our documentation to learn more about how Masa works and how you can get involved. Together, let's unlock the true potential of AI - powered by the people, for the people.

```shell
go build -v -o masa-node ./cmd/masa-node
```

##### 3. Go into the contracts directory and build the contract npm modules that the go binary uses

```shell
cd contracts/ 
yarn install
cd ../
```

##### 4. Set env vars using the following template

```plaintext
# Node Configuration

BOOTNODES=/ip4/35.223.224.220/udp/4001/quic-v1/p2p/16Uiu2HAmPxXXjR1XJEwckh6q1UStheMmGaGe8fyXdeRs3SejadSa

API_KEY=
ENV=test
FILE_PATH=.
PORT=8080
RPC_URL=https://ethereum-sepolia.publicnode.com
VALIDATOR=false

# AI LLM
CLAUDE_API_KEY=
CLAUDE_API_URL=https://api.anthropic.com/v1/messages
CLAUDE_API_VERSION=2023-06-01
ELAB_URL=https://api.elevenlabs.io/v1/text-to-speech/ErXwobaYiN019PkySvjV/stream
ELAB_KEY=
OPENAI_API_KEY=
PROMPT="You are a helpful assistant."

# Bring your own Twitter credentials
TWITTER_USER="yourusername"
TWITTER_PASS="yourpassword"
TWITTER_2FA_CODE="your2fa"

# Worker participation
TWITTER_SCRAPER=true
DISCORD_SCRAPER=true
WEB_SCRAPER=true
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

```shell
make build
```

The build command compiles the Masa Node binary and places it in the ./bin directory. It uses the go build command with the following flags:

-v: Enables verbose output to show the packages being compiled.
-o ./bin/masa-node: Specifies the output binary name and location.
./cmd/masa-node: Specifies the package to build (the main package).
make install
The install command runs the node_install.sh script to install any necessary dependencies or perform additional setup steps required by the Masa Node.

```shell
make run
```

The run command first builds the Masa Node binary using the build command and then executes the binary located at ./bin/masa-node. This command allows you to compile and run the Masa Node in a single step.

```shell
make test
```

The test command runs all the tests in the project using the go test command. It recursively searches for test files in all subdirectories and runs them.

```shell
make clean
```

The clean command performs cleanup tasks for the project. It removes the bin directory, which contains the compiled binary, and deletes the masa_node.log file, which may contain log output from previous runs.

To execute any of these commands, simply run make in your terminal from the project's root directory. For example, make build will compile the Masa Node binary, make test will run the tests, and make clean will remove the binary and log file.

## Funding the Node (in order to Stake)

```shell
  make faucet
```

>OR

Find the public key of your node in the logs.

Send 1000 MASA and .01 sepoliaETH to the node's public key / wallet address.

When the transactions have settled, you can stake

### Staking Tokens

- For local setup, stake tokens with:

  ```shell
  ./bin/masa-node --stake 1000
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

For more detailed documentation, please refer to the [CLI.md](md/CLI.md) file.

## Masa Node Twitter Sentiment Analysis

For more detailed documentation, please refer to the [LLM.md](md/LLM.md) file.

## API Swagger Docs

```shell
http://<masa-node-ip>:8080/swagger/index.html
```

## LLM Endpoint examples

ollama

```shell
curl http://localhost:8080/api/chat -d '{"model": "llama2","messages": [{"role": "user", "content": "why is the sky blue?" }], "stream": false}'
```

:::info
**Join our Community on Discord!**  
Ready to dive deeper into the Masa ecosystem? Connect with our vibrant community on Discord for the latest updates, discussions, and support. [Join us here](https://discord.gg/masafinance)
:::

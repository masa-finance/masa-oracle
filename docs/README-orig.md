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

2. Build the node:

```bash
go build -o masa-node
```

3. Run the node:

```bash
./masa-node   
```

You should see the node's address printed on the console. This indicates that your node is up and running, ready to connect with other Masa nodes.

---

## Connecting Nodes

Once you have the Masa node set up, you can easily connect multiple nodes together. Here's a step-by-step guide on how to do this:

### Starting a Listening Node

In one terminal window, start a node in listening mode:

```bash
$ ./masa-node
```

You should see an output similar to:

```bash
libp2p host address: /ip4/127.0.0.1/tcp/64924/p2p/12D3KooWSWGcJjMW75PL9LPpSFJGxf5wapXHe9V9auZD6hK1Tf26
```

This address is the multiaddress of the node. It provides all the necessary information for another peer to locate and communicate with this node.

### Connecting a Second Node

In another terminal window, run a second node and pass the multiaddress of the first node as a command-line argument:

```bash
./masa-node /ip4/127.0.0.1/tcp/49175/p2p/12D3KooWSWGcJjMW75PL9LPpSFJGxf5wapXHe9V9auZD6hK1Tf26
```

Once you run the command, the second node should attempt to connect to the first node. If successful, you might see some form of acknowledgment or interaction between the nodes, such as ping responses (depending on your implementation details).

## Introduction to Masa's Data Sources and Behavioral Tracking

Masa Oracle emerges as a groundbreaking solution in the web3 space, addressing the pressing need for a unified data layer that encapsulates a user's holistic behavior and identity data. In the decentralized realm, while the promise of privacy and control over one's data is paramount, the absence of a comprehensive data layer has led to fragmented experiences and inefficiencies.

### Masa's Data Sources

Masa Oracle leverages three core data sources, meticulously crafted to respect and enhance user privacy:

1. **Offchain Behavioral Data**: Using our advanced cookieless tracking mechanism, Masa captures detailed behavioral data without using traditional cookies. This unique approach ensures user behaviors are captured while maintaining utmost privacy.
2. **User Permissioned Offchain Data**: We harness data from platforms like Discord, Twitter, and also through processes such as Identity Verification and Sanctions Checks. Importantly, this data is accessed only with explicit user permission, emphasizing user control.
3. **Onchain Data**: Insights into asset ownership, historical balances, transactions, and DID credentials are accessed on-chain.

Using our proprietary cookieless tracking, Masa provides a deeper understanding of a web3 user, going beyond mere address identifications. By associating all user addresses through device sessions and formulating a unique Masa Identity, a comprehensive behavioral perspective of a user is realized. This identity ties together on-chain and off-chain behaviors and interactions.

### Use Cases

1. **Behavioral Analytics**: Projects can tap into the power of behavioral data to understand their users better, tailor experiences, and drive growth. With Masa, this doesn't come at the cost of privacy.
  
2. **Data-Driven Decentralized AI**: By providing a unified data layer, Masa fuels the next generation of AI models in the decentralized space. Data scientists can leverage this rich data, training AI models that are both powerful and privacy-preserving.
  
3. **Governance and Community Building**: With the comprehensive view that Masa provides, platforms can foster stronger communities. They can understand user needs better, drive engagement, and even facilitate governance mechanisms that are truly representative of the community's desires.

4. **Decentralized Identity Verification**: With Masa Oracle, platforms can seamlessly verify a user's identity without compromising on their privacy. From simple sign-ins to complex identity checks, Masa streamlines the process.


## Node Incentivization through Masa Tokens

### Earning Masa Tokens

Nodes participating in Masa Oracle are rewarded with the native Masa tokens. These are a token of appreciation for their contributions and to ensure the health and integrity of the decentralized data network.

### Revenue from Data Requests

Nodes have an additional revenue stream by servicing data requests. Revenue earned from these data requests, in various native currencies, is seamlessly converted into Masa tokens via a DEX (Decentralized Exchange).

## Protocol Governance & Voting

Community is at the heart of Masa Oracle. Node operators staking Masa tokens are empowered to participate in protocol governance, having a voice in proposals and pivotal decisions, ensuring the community drives the protocol's evolution.

## Technical Overview

Masa's Oracle ensures a robust decentralized system, with nodes processing and writing transactions from offchain data sources to the blockchain. With a focus on resilience, efficiency, and transparency, Masa Oracle stands out in the decentralized data landscape.

## Features

### Domain

- **Node**: Represents a participant in the network.
- **Transaction**: Represents the data that needs to be written to the blockchain.
- **Stake**: Represents the commitment a node makes to be a part of the network.
- **Webhook**: Represents the external data that triggers a transaction.

### Application

- **Node Service**: Manages node-related functionalities, including joining the network and peer handling.
- **Transaction Service**: Handles transaction processing and writing to the blockchain.
- **Stake Service**: Manages node staking, rewards, and penalties.

### Infrastructure

- **Libp2p**: Sets up decentralized node communication.
- **DHT**: Manages decentralized node registration and discovery.
- **Ethereum**: Handles interactions with Ethereum smart contracts.
- **DB**: Manages data persistence using the DHT.
- **Webhook**: Accepts and processes incoming webhook data.
- **Security**: Ensures secure communication and data handling.

### Utility Functions

Utility functions and common helpers are available for general operations, including rate limiting and error handling.

### Folder Structure
```
/masa-oracle
├── /domain              # Core business logic and entities
│   ├── /node
│   │   ├── node.go              # Entity
│   │   └── node_registered.go   # Domain Event
│   ├── /transaction
│   │   ├── transaction.go              # Entity
│   │   └── transaction_processed.go   # Domain Event
│   ├── /stake
│   │   ├── stake.go               # Entity
│   │   ├── stake_increased.go    # Domain Event
│   │   └── stake_decreased.go    # Domain Event
│   └── /webhook
│       └── webhook_data.go        # Value Object
├── /application         # Application's use-case-specific logic
│   ├── node_service.go
│   ├── transaction_service.go
│   └── stake_service.go
├── /infrastructure      # External tools, libraries, and modules
│   ├── /libp2p
│   │   ├── node_config.go
│   │   ├── peer_discovery.go
│   │   └── transport.go
│   ├── /dht
│   │   ├── dht_config.go
│   │   ├── node_registration.go
│   │   └── node_discovery.go
│   ├── /ethereum
│   │   ├── /contracts
│   │   │   ├── MasaToken.sol
│   │   │   └── StakingContract.sol
│   │   ├── staking.go
│   │   ├── rewards.go
│   │   └── truffle_config.go
│   ├── /db
│   │   ├── access_control.go
│   │   └── operations.go
│   ├── /webhook
│   │   ├── api_server.go
│   │   └── data_propagation.go
│   └── /security
│       ├── authentication.go
│       └── encryption.go
├── /utils                # Utility functions and common helpers
├── /tests                # Tests for the system
├── LICENSE
└── README.md
```

## Contribution

Contributions are always welcome. Please fork the repository and create a pull request with your changes. Ensure that your code follows Go best practices.

## License

This project is licensed under the terms of the [MIT license](LICENSE).

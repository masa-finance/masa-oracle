# Masa Oracle: Decentralized Data Protocol 🌐

The Masa Oracle defines how private behavioral and identity data is accessed, shared, and rewarded in a decentralized and private way. The Masa Oracle guarantees transparency, security, and equitable compensation for nodes that particiapte in the Masa zk-Data Network & Marketplace.

## Contents
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
- [Staking Tokens](#staking-tokens)
- [Running the Node](#running-the-node)
- [Updates & Additional Information](#updates--additional-information)

## Getting Started

### Prerequisites

Before diving in, ensure these prerequisites are installed:
- **Go**: Grab it from [Go's official site](https://golang.org/dl/).
- **Yarn**: Install it via [Yarn's official site](https://classic.yarnpkg.com/en/docs/install/).

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
3. Install node_modules in contracts directory

```bash
cd contracts/ && yarn install
```
4. Export RPC_URL to environment variable, can do 1 of 2 ways.

```bash
nano /Users/{USER}/.masa/masa_oracle_node.env and set RPC_URL=https://ethereum-sepolia.publicnode.com
```

```bash
export RPC_URL=https://ethereum-sepolia.publicnode.com
```
## Staking Tokens

🔐 To participate in the network and earn rewards, you must first stake your tokens:
```bash
./masa-node --stake 100
```
This command initiates the staking process, allowing the staking contract to spend tokens on your behalf and then stakes the specified amount.

## Running the Node

🚀 To start your node and join the Masa network you must connect to a bootnode: We have two bootnodes available for you to connect to. 
```
/ip4/34.121.111.128/udp/4001/quic-v1/p2p/16Uiu2HAmKULCxKgiQn1EcfKnq1Qam6psYLDTM99XsZFhr57wLadF

/ip4/35.223.224.220/udp/4001/quic-v1/p2p/16Uiu2HAmPxXXjR1XJEwckh6q1UStheMmGaGe8fyXdeRs3SejadSa
```

The command line parameters are as follows:
The command line parameters are as follows:  
`--bootnodes`: The multiaddress of the bootnode you want to connect to.  
`--port`: The port number you want to listen on.  
`--udp`: Enable or disable the UDP protocol. Right now set this to true as the bootnodes are using UDP.  
`--tcp`: Enable or disable the TCP protocol. Right now set this to false as the bootnodes are using UDP.  
`--start` Enable connection on the network.

```bash
./masa-node --bootnodes=/ip4/35.223.224.220/udp/4001/quic-v1/p2p/16Uiu2HAmPxXXjR1XJEwckh6q1UStheMmGaGe8fyXdeRs3SejadSa --port=4001 --udp=true --tcp=false --start=true
```
This will connect your Masa Oracle node to the specified node using the provided multi-address.


## Updates & Additional Information

📢 Stay tuned to the Masa Oracle repository for updates and more details on how to use the protocol effectively.

---

After setting up, your node's address will be displayed, indicating it's ready to connect with other Masa nodes. Follow any additional configuration steps and best practices as per your use case or network requirements.

# Basic libp2p Node

This repository contains a basic implementation of a libp2p node with resource management, DHT initialization, and the ping protocol. It serves as a foundational structure upon which more complex P2P applications can be built.

## Functionality

1. **Resource Management**:
    - Utilizes the Resource Manager from the `rcmgr` package to manage the node's resources.
    - Employs default scaling limits that adjust automatically based on the system's resources.

2. **libp2p Host Creation**:
    - Initiates a libp2p host that acts as an entity in the libp2p network, allowing other nodes to connect.
    - Listens on the IP address `127.0.0.1` (localhost) on a dynamically chosen port.

3. **DHT Initialization**:
    - Incorporates a Distributed Hash Table (DHT) using the Kademlia DHT implementation from `go-libp2p-kad-dht`.
    - Uses DHT for peer discovery, enabling the finding of other nodes in the network.
    - Bootstraps the DHT to populate its routing table by connecting to a few known nodes.

4. **Ping Protocol**:
    - Configured to handle ping requests using the `ping` protocol from libp2p.
    - Allows nodes to check the liveness of other nodes.
    - Serves as a fundamental debugging and network health-check tool.

5. **Display Host Address**:
    - Prints its own address to the console, which can be shared with other nodes for direct connection.

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

Once you have the libp2p node set up, you can easily connect multiple nodes together. Here's a step-by-step guide on how to do this:

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

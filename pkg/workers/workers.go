package workers

import (
	"fmt"

	"github.com/multiformats/go-multiaddr"

	"github.com/anthdm/hollywood/actor"
	masa "github.com/masa-finance/masa-oracle/pkg"
	msg "github.com/masa-finance/masa-oracle/pkg/proto/msg"
)

type Worker struct{}

// NewWorker creates a new instance of the Worker actor.
// It implements the actor.Receiver interface, allowing it to receive and handle messages.
//
// Returns:
//   - An instance of the Worker struct that implements the actor.Receiver interface.
func NewWorker() actor.Receiver {
	return &Worker{}
}

// Receive is the message handling method for the Worker actor.
// It receives messages through the actor context and processes them based on their type.
func (w *Worker) Receive(ctx *actor.Context) {
	switch m := ctx.Message().(type) {
	case *msg.Message:
		fmt.Println("actor received work", m.Data)
		// @todo the work

		// assumptions :
		// ✓ node must be staked
		// node must have n number of staked tokens?
		// do we want to offer scaled rewards based on how many tokens were staked?
		// how are the rewards distributed?

		// @todo consensus
		// 	- let un-staked / staked participate and infer the quality of their requests
		// 	- node uptime ie epoch/period
		// 	- staked
		// 	- how much staked
		// 	- participation rate
		// 	- let staked nodes rate each other
		// 	- let un-staked nodes rate each other

	default:
		break
	}
}

// SendWorkToPeers sends work data to peer nodes in the network.
// It subscribes to the local actor engine and the actor engines of peer nodes.
// The work data is then broadcast as an event to all subscribed nodes.
//
// Parameters:
//   - node: A pointer to the OracleNode instance.
//   - data: The work data to be sent, as a byte slice.
// Examples:
//	d, _ := json.Marshal(map[string]string{"request": "web", "url": "https://en.wikipedia.org/wiki/Maize", "depth": "2"})
//	d, _ := json.Marshal(map[string]string{"request": "twitter", "query": "$MASA", "count": "5", "model": "gpt-4"})
//	go workers.SendWorkToPeers(node, d)

func SendWorkToPeers(node *masa.OracleNode, data []byte) {
	node.ActorEngine.Subscribe(actor.NewPID("0.0.0.0:4001", fmt.Sprintf("%s/%s", "peer_worker", "peer")))
	peers := node.Host.Network().Peers()
	for _, peer := range peers {
		conns := node.Host.Network().ConnsToPeer(peer)
		for _, conn := range conns {
			addr := conn.RemoteMultiaddr()
			ipAddr, _ := addr.ValueForProtocol(multiaddr.P_IP4)
			peerPID := actor.NewPID(fmt.Sprintf("%s:4001", ipAddr), fmt.Sprintf("%s/%s", "peer_worker", "peer"))
			node.ActorEngine.Subscribe(peerPID)
		}
		// testing another peer on my local lan this will be removed
		// node.ActorEngine.Subscribe(actor.NewPID(fmt.Sprintf("%s:4001", "192.168.4.164"), fmt.Sprintf("%s/%s", "peer_worker", "peer")))
		// testing another peer on my local lan this will be removed
	}
	node.ActorEngine.BroadcastEvent(&msg.Message{Data: string(data)})
}

package workers

import (
	"fmt"

	"github.com/anthdm/hollywood/actor"
	masa "github.com/masa-finance/masa-oracle/pkg"
	msg "github.com/masa-finance/masa-oracle/pkg/proto/msg"
	"github.com/multiformats/go-multiaddr"
)

type Worker struct{}

func NewWorker() actor.Receiver {
	return &Worker{}
}

func (w *Worker) Receive(ctx *actor.Context) {
	switch m := ctx.Message().(type) {
	case *msg.Message:
		fmt.Println("actor received work", m.Data)
	default:
		break
	}
}

func SendWorkToPeers(node *masa.OracleNode, data string) {
	peers := node.Host.Network().Peers()
	for _, peer := range peers {
		conns := node.Host.Network().ConnsToPeer(peer)
		node.ActorEngine.Subscribe(actor.NewPID("0.0.0.0:4001", fmt.Sprintf("%s/%s", "peer_worker", "peer")))
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
	node.ActorEngine.BroadcastEvent(&msg.Message{Data: data})
}

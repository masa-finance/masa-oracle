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

func (f *Worker) Receive(ctx *actor.Context) {
	switch m := ctx.Message().(type) {
	case actor.Started:
		fmt.Println("actor started")
	case actor.Stopped:
		fmt.Println("actor stopped")
	case *msg.Message:
		fmt.Println("actor received work", m.Data)
	}
}

func SendWorkToPeers(node *masa.OracleNode, data string) {
	peers := node.Host.Network().Peers()
	for _, peer := range peers {
		conns := node.Host.Network().ConnsToPeer(peer)

		peerPIDLocal := actor.NewPID("0.0.0.0:4001", fmt.Sprintf("%s/%s", "peer_worker", "peer"))
		node.ActorEngine.Subscribe(peerPIDLocal)
		// test on my local lan this will be removed
		peerPID2 := actor.NewPID(fmt.Sprintf("%s:4001", "192.168.4.164"), fmt.Sprintf("%s/%s", "peer_worker", "peer"))
		node.ActorEngine.Subscribe(peerPID2)
		// test on my local lan this will be removed
		for _, conn := range conns {
			addr := conn.RemoteMultiaddr()
			ipAddr, _ := addr.ValueForProtocol(multiaddr.P_IP4)
			peerPID := actor.NewPID(fmt.Sprintf("%s:4001", ipAddr), fmt.Sprintf("%s/%s", "peer_worker", "peer"))
			node.ActorEngine.Subscribe(peerPID)
		}
	}
	node.ActorEngine.BroadcastEvent(&msg.Message{Data: data})
}

// @notes
// getPid := node.ActorEngine.Registry.GetPID("peer_worker", "peer")
// fmt.Println(getPid)

// use this where we want to stop an actor listener
// node.ActorEngine.Poison(pid).Wait()

package workers

import (
	"encoding/json"

	"github.com/multiformats/go-multiaddr"
	"github.com/sirupsen/logrus"

	masa "github.com/masa-finance/masa-oracle/pkg"
	"github.com/masa-finance/masa-oracle/pkg/pubsub"
	"github.com/masa-finance/masa-oracle/pkg/workers/messages"
)

// GetEligibleWorkers Uses the new NodeTracker method to get the eligible workers for a given message type
func GetEligibleWorkers(node *masa.OracleNode, message *messages.Work) []Worker {
	var workers []Worker
	// right now the message has the Twitter work type hard coded so we have to get it from the message data
	// TODO: can we get this fixed in the protobuf code?
	var workData map[string]string
	err := json.Unmarshal([]byte(message.Data), &workData)
	if err != nil {
		logrus.Errorf("[-] Error parsing work data: %v", err)
		return workers
	}

	workType, err := StringToWorkerType(workData["request"])
	if err != nil {
		logrus.Errorf("[-] Error parsing work type: %v", err)
		return workers
	}
	category := WorkerTypeToCategory(workType)
	for _, eligible := range node.NodeTracker.GetEligibleWorkerNodes(category) {
		if eligible.PeerId.String() == node.Host.ID().String() {
			workers = append(workers, Worker{IsLocal: true, NodeData: pubsub.NodeData{PeerId: node.Host.ID()}})
			continue
		}
		for _, addr := range eligible.Multiaddrs {
			ipAddr, _ := addr.ValueForProtocol(multiaddr.P_IP4)
			workers = append(workers, Worker{IsLocal: false, NodeData: eligible, IPAddr: ipAddr})
			break
		}
	}
	return workers
}

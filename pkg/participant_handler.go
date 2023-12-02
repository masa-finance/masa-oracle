package masa

import (
	"encoding/json"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/sirupsen/logrus"
)

type ParticipantHandler struct {
	Node *OracleNode
}

func NewParticipantHandler(node *OracleNode) *ParticipantHandler {
	return &ParticipantHandler{
		Node: node,
	}
}

func (handler *ParticipantHandler) HandleMessage(msg *pubsub.Message) {
	var nodeData NodeData
	if err := json.Unmarshal(msg.Data, &nodeData); err != nil {
		logrus.Errorf("failed to unmarshal node data: %v", err)
	}
	// Handle the nodeData by calling NodeEventTracker.HandleIncomingData
	handler.Node.NodeTracker.HandleIncomingData(&nodeData)
}

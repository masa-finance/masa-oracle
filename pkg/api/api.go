package api

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	node "github.com/masa-finance/masa-oracle/node"
	"github.com/masa-finance/masa-oracle/pkg/event"
)

type API struct {
	Node         *node.OracleNode
	EventTracker *event.EventTracker
}

// NewAPI creates a new API instance with the given OracleNode.
func NewAPI(node *node.OracleNode) *API {
	eventTracker := event.NewEventTracker(nil)
	if eventTracker == nil {
		logrus.Error("Failed to create EventTracker")
	} else {
		logrus.Debug("EventTracker created successfully")
	}

	api := &API{
		Node:         node,
		EventTracker: eventTracker,
	}

	logrus.Debugf("Created API instance with EventTracker: %v", api.EventTracker)
	return api
}

// GetPathInt converts the path parameter with name to an int.
// It returns the int value and nil error if the path parameter is present and a valid integer.
// It returns 0 and a formatted error if the path parameter is missing or not a valid integer.
func GetPathInt(ctx *gin.Context, name string) (int, error) {
	val, ok := ctx.GetQuery(name)
	if !ok {
		return 0, fmt.Errorf("the value for path parameter %s empty or not specified", name)
	}
	return strconv.Atoi(val)
}

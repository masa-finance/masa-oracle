package tests

import (
	"testing"

	"github.com/masa-finance/masa-oracle/node"
	"github.com/masa-finance/masa-oracle/pkg/api"
	"github.com/masa-finance/masa-oracle/pkg/pubsub"
	"github.com/masa-finance/masa-oracle/pkg/workers"
)

func TestAPI(t *testing.T) {
	// Create a new OracleNode instance
	n := &node.OracleNode{}
	whm := &workers.WorkHandlerManager{}
	pubKeySub := &pubsub.PublicKeySubscriptionHandler{}

	// Initialize the API
	api := api.NewAPI(n, whm, pubKeySub)

	// Test API initialization
	if api == nil {
		t.Fatal("Failed to initialize API")
	}
}

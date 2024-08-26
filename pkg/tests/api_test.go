package tests

import (
	"testing"

	"github.com/masa-finance/masa-oracle/node"
	"github.com/masa-finance/masa-oracle/pkg/api"
)

func TestAPI(t *testing.T) {
	// Create a new OracleNode instance
	n := &node.OracleNode{}

	// Initialize the API
	api := api.NewAPI(n)

	// Test API initialization
	if api == nil {
		t.Fatal("Failed to initialize API")
	}

}

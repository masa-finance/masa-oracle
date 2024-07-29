package tests

import (
	"testing"

	masa "github.com/masa-finance/masa-oracle/pkg"
	"github.com/masa-finance/masa-oracle/pkg/api"
)

func TestAPI(t *testing.T) {
	// Create a new OracleNode instance
	node := &masa.OracleNode{}

	// Initialize the API
	api := api.NewAPI(node)

	// Test API initialization
	if api == nil {
		t.Fatal("Failed to initialize API")
	}

}

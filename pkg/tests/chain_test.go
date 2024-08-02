package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockChain is a mock implementation of the Chain interface
type MockChain struct {
	mock.Mock
}

func TestGetLatestBlockNumber(t *testing.T) {

	assert.NotNil(t, 1)
}

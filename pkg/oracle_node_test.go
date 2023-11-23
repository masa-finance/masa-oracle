package masa

import "testing"

func TestIsPublisher(t *testing.T) {
	node := &OracleNode{
		Stake:      101,
		Reputation: 101,
	}

	if !node.IsPublisher() {
		t.Errorf("Expected node to be a publisher, but it's not")
	}

	node.Stake = 50
	if node.IsPublisher() {
		t.Errorf("Expected node not to be a publisher, but it is")
	}
}

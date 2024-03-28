package network

import (
	"testing"

	"github.com/multiformats/go-multiaddr"
	"github.com/stretchr/testify/assert"
)

func TestGetPriorityAddress(t *testing.T) {

	// Create a list of multiaddresses with different IP types.
	loopbackAddr, _ := multiaddr.NewMultiaddr("/ip4/127.0.0.1")
	privateAddr, _ := multiaddr.NewMultiaddr("/ip4/192.168.1.1")
	publicAddr, _ := multiaddr.NewMultiaddr("/ip4/93.184.216.34") // Example public IP
	addrs := []multiaddr.Multiaddr{loopbackAddr, privateAddr, publicAddr}

	// Call the function under test.
	selectedAddr := GetPriorityAddress(addrs)

	// Convert the selected multiaddress to a net.Addr to check the IP.
	netAddr, err := selectedAddr.ValueForProtocol(multiaddr.P_IP4)
	if err != nil {
		t.Fatalf("Failed to extract IP from multiaddress: %v", err)
	}

	assert.Equal(t, "93.184.216.34", netAddr, "Expected the public IP to be selected.")

	// Now test with no public IP available.
	addrs = []multiaddr.Multiaddr{loopbackAddr, privateAddr}
	selectedAddr = GetPriorityAddress(addrs)
	netAddr, err = selectedAddr.ValueForProtocol(multiaddr.P_IP4)
	if err != nil {
		t.Fatalf("Failed to extract IP from multiaddress: %v", err)
	}

	// Assert that the selected IP is the private IP since no public IP was available.
	assert.Equal(t, "192.168.1.1", netAddr, "Expected the private IP to be selected when no public IP is available.")
}

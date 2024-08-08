package network

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/chyeh/pubip"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	"github.com/sirupsen/logrus"
)

// getOutboundIP returns the outbound IP address of the current machine 172.17.0.2 10.0.0.2 etc
func getOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		logrus.Warn("Error getting outbound IP")
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().String()
	idx := strings.LastIndex(localAddr, ":")
	return localAddr[0:idx]
}

// GetMultiAddressesForHost returns the multiaddresses for the host
func GetMultiAddressesForHost(host host.Host) ([]multiaddr.Multiaddr, error) {
	peerInfo := peer.AddrInfo{
		ID:    host.ID(),
		Addrs: host.Addrs(),
	}
	multiaddrs, err := peer.AddrInfoToP2pAddrs(&peerInfo)
	if err != nil {
		return nil, err
	}
	logrus.Debug("Multiaddresses from AddrInfoToP2pAddrs: ", multiaddrs)
	addresses := make([]multiaddr.Multiaddr, 0)
	for _, addr := range multiaddrs {
		logrus.Debug(addr.String())
		// skip using localhost since it provides no value
		if !strings.Contains(addr.String(), "127.0.0.1") {
			addresses = append(addresses, addr)
		}
	}
	return addresses, nil
}

// GetMultiAddressesForHostQuiet returns the multiaddresses for the host without logging
func GetMultiAddressesForHostQuiet(host host.Host) []multiaddr.Multiaddr {
	ma, err := GetMultiAddressesForHost(host)
	if err != nil {
		logrus.Fatal(err)
	}
	return ma
}

// GetPriorityAddress returns the best public or private IP address
func GetPriorityAddress(addrs []multiaddr.Multiaddr) multiaddr.Multiaddr {
	bestAddr := getBestPublicAddress(addrs)
	if bestAddr != nil {
		return bestAddr
	}

	var bestPrivateAddr multiaddr.Multiaddr

	for _, addr := range addrs {
		ipComponent, err := addr.ValueForProtocol(multiaddr.P_IP4)
		if err != nil {
			continue // Not an IP address
		}

		ip := net.ParseIP(ipComponent)
		if ip == nil || ip.IsLoopback() {
			continue // Skip invalid or loopback addresses
		}

		if ip.IsPrivate() {
			if bestPrivateAddr == nil || isPreferredAddress(addr) {
				bestPrivateAddr = addr
			}
		}
	}

	if bestPrivateAddr != nil {
		return bestPrivateAddr
	}

	if len(addrs) > 0 {
		logrus.Warn("No suitable address found, returning the first entry")
		return addrs[0]
	}

	return nil
}

// getBestPublicAddress returns the best public IP address
func getBestPublicAddress(addrs []multiaddr.Multiaddr) multiaddr.Multiaddr {
	var externalIP net.IP
	var err error

	if os.Getenv("ENV") == "local" {
		externalIP = net.ParseIP(getOutboundIP())
	} else {
		externalIP, err = pubip.Get()
		if err != nil {
			logrus.Warnf("Failed to get public IP: %v", err)
			return nil
		}
	}

	if externalIP == nil || externalIP.IsPrivate() {
		return nil
	}

	// Create a new multiaddr with the public IP
	publicAddr, err := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/%s", externalIP.String()))
	if err != nil {
		logrus.Warnf("Failed to create multiaddr with public IP: %v", err)
		return nil
	}

	// Find a suitable port from existing addresses
	for _, addr := range addrs {
		if strings.HasPrefix(addr.String(), "/ip4/") {
			port, err := addr.ValueForProtocol(multiaddr.P_TCP)
			if err == nil {
				return publicAddr.Encapsulate(multiaddr.StringCast("/tcp/" + port))
			}
		}
	}

	return publicAddr
}

// isPreferredAddress checks if the multiaddress contains the UDP protocol
func isPreferredAddress(addr multiaddr.Multiaddr) bool {
	// Check if the multiaddress contains the UDP protocol
	for _, p := range addr.Protocols() {
		if p.Code == multiaddr.P_UDP {
			return true
		}
	}
	return false
}

// GetBootNodesMultiAddress returns the multiaddresses for the bootstrap nodes
func GetBootNodesMultiAddress(bootstrapNodes []string) ([]multiaddr.Multiaddr, error) {
	addrs := make([]multiaddr.Multiaddr, 0)
	for _, peerAddr := range bootstrapNodes {
		if peerAddr == "" {
			continue
		}
		addr, err := multiaddr.NewMultiaddr(peerAddr)
		if err != nil {
			return nil, err
		}
		addrs = append(addrs, addr)
	}
	return addrs, nil
}

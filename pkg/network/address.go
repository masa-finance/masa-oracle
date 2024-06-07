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

func GetMultiAddressesForHostQuiet(host host.Host) []multiaddr.Multiaddr {
	ma, err := GetMultiAddressesForHost(host)
	if err != nil {
		logrus.Fatal(err)
	}
	return ma
}

func GetPriorityAddress(addrs []multiaddr.Multiaddr) multiaddr.Multiaddr {
	var bestPublicAddr multiaddr.Multiaddr
	var bestPrivateAddr multiaddr.Multiaddr

	for _, addr := range addrs {

		ipComponent, err := addr.ValueForProtocol(multiaddr.P_IP4)
		if err != nil {
			ipComponent, err = addr.ValueForProtocol(multiaddr.P_IP6)
			if err != nil {
				continue // Not an IP address
			}
		}

		ip := net.ParseIP(ipComponent)
		if ip == nil || ip.IsLoopback() {
			continue // Skip invalid or loopback addresses
		}

		if ip.IsPrivate() {
			// If it's the first private address found or if it's preferred over the current best, keep it
			if bestPrivateAddr == nil || isPreferredAddress(addr) {
				bestPrivateAddr = addr
			}
		} else {
			// If it's the first public address found or if it's preferred over the current best, keep it
			if bestPublicAddr == nil || isPreferredAddress(addr) {
				bestPublicAddr = addr
			}
		}
	}
	var baseAddr multiaddr.Multiaddr
	// Prefer public addresses over private ones
	if bestPublicAddr != nil {
		baseAddr = bestPublicAddr
	} else if bestPrivateAddr != nil {
		baseAddr = bestPrivateAddr
	} else {
		logrus.Warn("No address matches the priority criteria, returning the first entry")
		baseAddr = addrs[0]
	}
	logrus.Debug("Best public address: ", bestPublicAddr)
	logrus.Debug("Best private address: ", bestPrivateAddr)
	logrus.Infof("Base address: %s", baseAddr)
	addr := replaceAddress(baseAddr)
	if addr != nil {
		baseAddr = addr
	}
	return baseAddr
}

func isPreferredAddress(addr multiaddr.Multiaddr) bool {
	// Check if the multiaddress contains the UDP protocol
	for _, p := range addr.Protocols() {
		if p.Code == multiaddr.P_UDP {
			return true
		}
	}
	return false
}

func getOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		fmt.Println("Error getting outbound IP")
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().String()
	idx := strings.LastIndex(localAddr, ":")
	return localAddr[0:idx]
}

func replaceAddress(addr multiaddr.Multiaddr) multiaddr.Multiaddr {
	var err error
	var bestAddr multiaddr.Multiaddr

	var externalIP net.IP
	if os.Getenv("ENV") == "local" {
		externalIP = net.ParseIP(getOutboundIP())
	} else {
		externalIP, _ = pubip.Get()
	}

	if externalIP.String() != "" {
		bestAddr, err = replaceIPComponent(addr, externalIP.String())
		if err != nil {
			logrus.Warnf("Failed to replace IP component: %s", err)
			return nil
		}
	}
	logrus.Debug("Address after replacing IP component: ", bestAddr)
	return bestAddr
}

func replaceIPComponent(maddr multiaddr.Multiaddr, newIP string) (multiaddr.Multiaddr, error) {
	var components []multiaddr.Multiaddr
	for _, component := range multiaddr.Split(maddr) {
		if component.Protocols()[0].Code == multiaddr.P_IP4 || component.Protocols()[0].Code == multiaddr.P_IP6 {
			// Create a new IP component
			newIPComponent, err := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/%s", newIP))
			if err != nil {
				return nil, err
			}
			components = append(components, newIPComponent)
		} else {
			components = append(components, component)
		}
	}
	return multiaddr.Join(components...), nil
}

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

package network

import (
	"net"
	"strings"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr/net"
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
	addresses := make([]multiaddr.Multiaddr, 0)
	for _, addr := range multiaddrs {
		logrus.Debug(addr.String())
		// skip using localhost since it provides no value
		if strings.Contains(addr.String(), "127.0.0.1") {
			addresses = append(addresses, addr)
		}
	}
	return addresses, nil
}

func GetMultiAddressesForHostQuiet(host host.Host) []multiaddr.Multiaddr {
	multiaddr, err := GetMultiAddressesForHost(host)
	if err != nil {
		logrus.Fatal(err)
	}
	return multiaddr
}

func GetPriorityAddress(addrs []multiaddr.Multiaddr) multiaddr.Multiaddr {
	var udpQUIC, tcp, public, private, nonLocal []multiaddr.Multiaddr

	for _, addr := range addrs {
		if strings.Contains(addr.String(), "/udp/") || strings.Contains(addr.String(), "/quic/") {
			udpQUIC = append(udpQUIC, addr)
		} else if strings.Contains(addr.String(), "/tcp/") {
			tcp = append(tcp, addr)
		}
		netAddr, err := manet.ToNetAddr(addr)
		if err != nil {
			continue
		}

		tcpAddr, ok := netAddr.(*net.TCPAddr)
		if !ok {
			continue
		}

		ip := tcpAddr.IP
		if !ip.IsLoopback() {
			nonLocal = append(nonLocal, addr)
			if !ip.IsPrivate() {
				public = append(public, addr)
			} else {
				private = append(private, addr)
			}
		}
	}

	// Prioritize UDP/QUIC over TCP
	if len(udpQUIC) > 0 {
		// Prioritize public over private and non-local
		for _, addr := range udpQUIC {
			if contains(public, addr) {
				return addr
			}
		}
		for _, addr := range udpQUIC {
			if contains(nonLocal, addr) {
				return addr
			}
		}
	} else if len(tcp) > 0 {
		// Prioritize public over private and non-local
		for _, addr := range tcp {
			if contains(public, addr) {
				return addr
			}
		}
		for _, addr := range tcp {
			if contains(nonLocal, addr) {
				return addr
			}
		}
	}
	// If no address matches the criteria, return the first entry
	logrus.Warn("No address matches the priority criteria, returning the first entry")
	return addrs[0]
}

func contains(slice []multiaddr.Multiaddr, item multiaddr.Multiaddr) bool {
	for _, a := range slice {
		if a.Equal(item) {
			return true
		}
	}
	return false
}

func GetBootNodesMultiAddress(input string) ([]multiaddr.Multiaddr, error) {
	logrus.Infof("Getting bootnodes from %s", input)
	bootstrapPeers := strings.Split(input, ",")
	addrs := make([]multiaddr.Multiaddr, 0)
	for _, peerAddr := range bootstrapPeers {
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

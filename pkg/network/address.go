package network

import (
	"strings"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	"github.com/sirupsen/logrus"
)

func GetMultiAddressForHost(host host.Host) (multiaddr.Multiaddr, error) {
	peerInfo := peer.AddrInfo{
		ID:    host.ID(),
		Addrs: host.Addrs(),
	}
	multiaddrs, err := peer.AddrInfoToP2pAddrs(&peerInfo)
	if err != nil {
		return nil, err
	}
	return multiaddrs[0], nil
}

func GetMultiAddressForHostQuiet(host host.Host) multiaddr.Multiaddr {
	multiaddr, err := GetMultiAddressForHost(host)
	if err != nil {
		logrus.Fatal(err)
	}
	return multiaddr
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

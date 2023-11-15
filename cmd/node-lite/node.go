package main

import (
	"bufio"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
)

func NewNodeLite() error {
	//ctx := context.Background()

	// libp2p.New constructs a new libp2p Host.
	// Other options can be added here.
	host, err := libp2p.New()
	if err != nil {
		return err
	}
	host.SetStreamHandler("/chat/1.1.0", handleStream)

	n := &discoveryNotifee{}
	service := mdns.NewMdnsService(host, "/chat/1.0.0", n)
	if err != nil {
		return err
	}
	if err := service.Start(); err != nil {
		return err
	}
	return nil
}

func handleStream(stream network.Stream) {

	// Create a buffer stream for non-blocking read and write.
	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

	go readData(rw)
	go writeData(rw)

	// 'stream' will stay open until you close it (or the other side closes it).
}

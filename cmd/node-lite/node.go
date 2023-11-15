package main

import (
	"bufio"
	"context"
	"fmt"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/multiformats/go-multiaddr"
	"github.com/sirupsen/logrus"

	myNetwork "github.com/masa-finance/masa-oracle/pkg/network"
)

const nodeLiteProtocol = "masa_node_lite_protocol/1.0.0"

type NodeLite struct {
	Host       host.Host
	PrivKey    crypto.PrivKey
	Protocol   protocol.ID
	multiAddrs multiaddr.Multiaddr
	ctx        context.Context
}

func NewNodeLite(privKey crypto.PrivKey, ctx context.Context) *NodeLite {

	host, err := libp2p.New(
		libp2p.Identity(privKey),
	)
	if err != nil {
		logrus.Fatal(err)
	}
	return &NodeLite{
		Host:       host,
		PrivKey:    privKey,
		Protocol:   nodeLiteProtocol,
		multiAddrs: myNetwork.GetMultiAddressForHostQuiet(host),
		ctx:        ctx,
	}
	return nil
}

func (node *NodeLite) Start() (err error) {
	node.StartMDNSDiscovery(string(node.Protocol))
	return nil
}

func (node *NodeLite) StartMDNSDiscovery(rendezvous string) {
	peerChan := myNetwork.StartMDNS(node.Host, rendezvous)
	go func() {
		for {
			select {
			case peer := <-peerChan: // will block until we discover a peer
				logrus.Info("Found peer:", peer, ", connecting")

				if err := node.Host.Connect(node.ctx, peer); err != nil {
					logrus.Error("Connection failed:", err)
					continue
				}

				// open a stream, this stream will be handled by handleStream other end
				stream, err := node.Host.NewStream(node.ctx, peer.ID, node.Protocol)

				if err != nil {
					logrus.Error("Stream open failed", err)
				} else {
					rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

					go node.writeData(rw)
					go node.readData(rw)
					logrus.Info("Connected to:", peer)
				}
			case <-node.ctx.Done():
				return
			}
		}
	}()
}

func (node *NodeLite) readData(rw *bufio.ReadWriter) {
	for {
		str, err := rw.ReadString('\n')
		if err != nil {
			logrus.Error("Error reading from buffer:", err)
			return
		}

		if str != "" && str != "\n" {
			logrus.Infof("MDNS Received message: %s from %s", str, node.multiAddrs.String())
		}
	}
}

func (node *NodeLite) writeData(rw *bufio.ReadWriter) {
	for {
		// Generate a message including the multiaddress of the sender
		sendData := fmt.Sprintf("MDNS Hello from %s\n", node.multiAddrs.String())

		_, err := rw.WriteString(sendData)
		if err != nil {
			logrus.Error("Error writing to buffer:", err)
			return
		}
		err = rw.Flush()
		if err != nil {
			logrus.Error("Error flushing buffer:", err)
			return
		}
		// Sleep for a while before sending the next message
		time.Sleep(time.Second * 5)
	}
}

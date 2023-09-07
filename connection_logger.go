package main

import (
	"github.com/libp2p/go-libp2p/core/network"
	ma "github.com/multiformats/go-multiaddr"
	log "github.com/sirupsen/logrus"
)

type ConnectionLogger struct{}

func (cl *ConnectionLogger) Listen(n network.Network, a ma.Multiaddr) {
	log.WithFields(log.Fields{
		"network": n,
		"address": a,
	}).Info("Started listening")
}

func (cl *ConnectionLogger) ListenClose(n network.Network, a ma.Multiaddr) {
	log.WithFields(log.Fields{
		"network": n,
		"address": a,
	}).Info("Stopped listening")
}

func (cl *ConnectionLogger) Connected(n network.Network, c network.Conn) {
	log.WithFields(log.Fields{
		"network": n,
		"conn":    c,
	}).Info("Connected")
}

func (cl *ConnectionLogger) Disconnected(n network.Network, c network.Conn) {
	log.WithFields(log.Fields{
		"network": n,
		"conn":    c,
	}).Info("Disconnected")
}

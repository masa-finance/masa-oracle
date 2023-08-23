package tests

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/transport/websocket"
	"github.com/multiformats/go-multiaddr"
	"github.com/sirupsen/logrus"
)

func TestWebsocketListener(t *testing.T) {
	nl, err := NewNodeListener("/ip4/0.0.0.0/tcp/3000/ws")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	err = nl.Start()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	messageSender(nl.Address, nl.ServerId)
	//connectAndExchangeData(multiaddr.StringCast("/ip4/127.0.0.1/tcp/0/wss"), true)
}

func TestWebsocketConnection(t *testing.T) {
	//t.Run("unencrypted", func(t *testing.T) {
	//	connectAndExchangeData(multiaddr.StringCast("/ip4/127.0.0.1/tcp/0/ws"), false)
	//})
	t.Run("encrypted", func(t *testing.T) {
		connectAndExchangeData(multiaddr.StringCast("/ip4/127.0.0.1/tcp/0/wss"), true)
	})
}

func messageSender(addr multiaddr.Multiaddr, id peer.ID) {
	fmt.Println("listening on", addr)
	var opts []websocket.Option
	opts = append(opts, websocket.WithTLSClientConfig(&tls.Config{InsecureSkipVerify: true}))
	privKey, _, err := crypto.GenerateKeyPair(crypto.RSA, 2048)
	if err != nil {
		logrus.Error(err)
	}

	_, u, err := newUpgrader(privKey)
	if err != nil {
		logrus.Error(err)
	}

	tpt, err := websocket.New(u, &network.NullResourceManager{}, opts...)
	if err != nil {
		logrus.Error(err)
	}
	conn, err := tpt.Dial(context.Background(), addr, id)
	if err != nil {
		logrus.Error(err)
	}
	str, err := conn.OpenStream(context.Background())
	if err != nil {
		logrus.Error(err)
	}
	defer str.Close()
	for {
		fmt.Printf("sending message to %s\n", addr)
		_, err = str.Write([]byte("test message"))
		if err != nil {
			logrus.Error(err)
		}
		time.Sleep(5 * time.Second)
	}
}

func connectAndExchangeData(laddr multiaddr.Multiaddr, secure bool) error {
	var opts []websocket.Option
	var tlsConf *tls.Config
	var err error
	if secure {
		tlsConf, err = generateTLSConfig()
		opts = append(opts, websocket.WithTLSConfig(tlsConf))
	}
	privKey, _, err := crypto.GenerateKeyPair(crypto.RSA, 2048)
	if err != nil {
		return err
	}
	server, u, err := newUpgrader(privKey)
	messageSender(laddr, server)

	tpt, err := websocket.New(u, &network.NullResourceManager{}, opts...)
	if err != nil {
		return err
	}

	l, err := tpt.Listen(laddr)
	if err != nil {
		return err
	}
	defer l.Close()

	c, err := l.Accept()
	defer c.Close()
	str, err := c.AcceptStream()
	defer str.Close()

	out, err := io.ReadAll(str)
	if err != nil {
		return err
	}
	fmt.Printf("received message: %s\n", string(out))
	return nil
}

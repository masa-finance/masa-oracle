package main

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/libp2p/go-libp2p/core/host"
)

func TestOracleNodeCommunication(t *testing.T) {
	// Create two OracleNodes
	// Get or create the private key
	privKey1, err := getOrCreatePrivateKey(filepath.Join("tests/node1/private.key"))
	if err != nil {
		t.Fatal(err)
	}
	node1, err := NewOracleNode(privKey1)
	if err != nil {
		t.Fatal(err)
	}
	err = node1.Start()
	if err != nil {
		t.Fatal(err)
	}

	privKey2, err := getOrCreatePrivateKey(filepath.Join("tests/node2/private.key"))
	if err != nil {
		t.Fatal(err)
	}
	node2, err := NewOracleNode(privKey2)
	if err != nil {
		t.Fatal(err)
	}
	err = node2.Start()
	if err != nil {
		t.Fatal(err)
	}

	node1.Connect(node2)
	if err != nil {
		t.Fatal(err)
	}
	node2Info := host.InfoFromHost(node2.Host)
	// Send 5 messages from node1 to node2
	for i := 0; i < 5; i++ {
		// Open a stream to node2
		stream, err := node1.Host.NewStream(context.Background(), node2Info.ID, "/masa_oracle_protocol/1.0.0")
		if err != nil {
			t.Fatal(err)
		}

		// Write a message to the stream
		_, err = stream.Write([]byte(fmt.Sprintf("Message %d from node1", i+1)))
		if err != nil {
			t.Fatal(err)
		}

		// Close the stream
		ack := make([]byte, 3)
		_, err = stream.Read(ack)
		if err != nil {
			t.Fatal(err)
		}
		if string(ack) != "ACK" {
			t.Fatalf("Did not receive expected acknowledgement for message %d", i+1)
		}

		// Close the stream
		err = stream.Close()
		if err != nil {
			t.Fatal(err)
		}

		// Wait for 2 seconds before sending the next message
		time.Sleep(2 * time.Second)
	}
}

package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	masa "github.com/masa-finance/masa-oracle/pkg"
	"github.com/masa-finance/masa-oracle/pkg/crypto"
)

func init() {
	f, err := os.OpenFile("masa_node_lite.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		logrus.Fatal(err)
	}
	mw := io.MultiWriter(os.Stdout, f)
	logrus.SetOutput(mw)
	logrus.SetLevel(logrus.InfoLevel)

	usr, err := user.Current()
	if err != nil {
		logrus.Fatal("could not find user.home directory")
	}
	envFilePath := filepath.Join(usr.HomeDir, ".masa", "masa_oracle_node.env")
	keyFilePath := filepath.Join(usr.HomeDir, ".masa", "masa_oracle_key")

	// Create the directories if they don't already exist
	if _, err := os.Stat(filepath.Dir(envFilePath)); os.IsNotExist(err) {
		err = os.MkdirAll(filepath.Dir(envFilePath), 0755)
		if err != nil {
			logrus.Fatal("could not create directory:", err)
		}
	}
	// Check if the .env file exists
	if _, err := os.Stat(envFilePath); os.IsNotExist(err) {
		// If not, create it with default values
		builder := strings.Builder{}
		builder.WriteString(fmt.Sprintf("%s=%s\n", masa.KeyFileKey, keyFilePath))
		err = os.WriteFile(envFilePath, []byte(builder.String()), 0644)
		if err != nil {
			logrus.Fatal("could not write to .env file:", err)
		}
	}
	err = godotenv.Load(envFilePath)
	if err != nil {
		logrus.Error("Error loading .env file")
	}
}

func main() {
	logrus.Infof("arg size is %d", len(os.Args))
	if len(os.Args) > 1 {
		logrus.Infof("found arg: %s", os.Args[1])
		err := os.Setenv(masa.Peers, os.Args[1])
		if err != nil {
			logrus.Error(err)
		}
		if len(os.Args) == 3 {
			err := os.Setenv(masa.PortNbr, os.Args[2])
			if err != nil {
				logrus.Error(err)
			}
		}
	}
	// Create a cancellable context
	ctx, cancel := context.WithCancel(context.Background())

	// Listen for SIGINT (CTRL+C)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Cancel the context when SIGINT is received
	go func() {
		<-c
		cancel()
	}()

	privKey, err := crypto.GetOrCreatePrivateKey(os.Getenv(masa.KeyFileKey))
	if err != nil {
		logrus.Fatal(err)
	}
	node, err := NewNodeLite(privKey, ctx)
	if err != nil {
		logrus.Fatal(err)
	}

	node.Start()

	//
	//// Set a function as stream handler.
	//// This function is called when a peer initiates a connection and starts a stream with this peer.
	//host.SetStreamHandler("/masa_chat/1.0.0", handleStream)
	//
	//fmt.Printf("\n[*] Your Multiaddress Is: /ip4/%s/tcp/%v/p2p/%s\n", "0.0.0.0", 4001, host.ID())
	//
	//peerChan := network.StartMDNS(host, "masa-chat")
	//for { // allows multiple peers to join
	//	peer := <-peerChan // will block until we discover a peer
	//	fmt.Println("Found peer:", peer, ", connecting")
	//
	//	if err := host.Connect(ctx, peer); err != nil {
	//		fmt.Println("Connection failed:", err)
	//		continue
	//	}
	//
	//	// open a stream, this stream will be handled by handleStream other end
	//	stream, err := host.NewStream(ctx, peer.ID, "/masa_chat/1.0.0")
	//
	//	if err != nil {
	//		fmt.Println("Stream open failed", err)
	//	} else {
	//		rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))
	//
	//		go writeData(rw)
	//		go readData(rw)
	//		fmt.Println("Connected to:", peer)
	//	}
	//}
	<-ctx.Done()
}

func readData(rw *bufio.ReadWriter) {
	for {
		str, err := rw.ReadString('\n')
		if err != nil {
			logrus.Error("Error reading from buffer:", err)
		}

		if str == "" {
			return
		}
		if str != "\n" {
			// Green console colour: 	\x1b[32m
			// Reset console colour: 	\x1b[0m
			fmt.Printf("\x1b[32m%s\x1b[0m> ", str)
		}

	}
}

func writeData(rw *bufio.ReadWriter) {
	stdReader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		sendData, err := stdReader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from stdin")
			panic(err)
		}

		_, err = rw.WriteString(fmt.Sprintf("%s\n", sendData))
		if err != nil {
			fmt.Println("Error writing to buffer")
			panic(err)
		}
		err = rw.Flush()
		if err != nil {
			fmt.Println("Error flushing buffer")
			panic(err)
		}
	}
}

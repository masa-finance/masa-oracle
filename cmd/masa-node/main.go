package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	masa "github.com/masa-finance/masa-oracle/pkg"
	api "github.com/masa-finance/masa-oracle/pkg/api"
	"github.com/masa-finance/masa-oracle/pkg/crypto"
)

var (
	bootnodes string
	portNbr   int
	udp       bool
	tcp       bool
)

func init() {
	f, err := os.OpenFile("masa_oracle_node.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(err)
	}
	mw := io.MultiWriter(os.Stdout, f)
	logrus.SetOutput(mw)
	if os.Getenv("debug") == "true" {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}

	usr, err := user.Current()
	if err != nil {
		log.Fatal("could not find user.home directory")
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

	// Define flags
	flag.StringVar(&bootnodes, "bootnodes", os.Getenv("BOOTNODES"), "A comma-separated list of multiAddress strings")
	flag.IntVar(&portNbr, "port", getPort("portNbr"), "The port number")
	flag.BoolVar(&udp, "udp", getEnvAsBool("UDP", false), "UDP flag")
	flag.BoolVar(&tcp, "tcp", getEnvAsBool("TCP", false), "TCP flag")
	flag.Parse()

	err = os.Setenv(masa.Peers, bootnodes)
	if err != nil {
		logrus.Error(err)
	}
	//if neither udp nor tcp are set, default to udp
	if !udp && !tcp {
		udp = true
	}
}

func main() {
	// log the flags
	bootnodesList := strings.Split(bootnodes, ",")
	logrus.Infof("Bootnodes: %v", bootnodesList)
	logrus.Infof("Port number: %s", portNbr)
	logrus.Infof("UDP: %v", udp)
	logrus.Infof("TCP: %v", tcp)

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
	//crypto.VerifyEthereumCompatibility(privKey)

	node, err := masa.NewOracleNode(ctx, privKey, portNbr, udp, tcp)
	if err != nil {
		logrus.Fatal(err)
	}
	err = node.Start()
	if err != nil {
		logrus.Fatal(err)
	}

	// BP: Add gin router to get peers (multiaddress) and get peer addresses @Bob - I am not sure if this is the right place for this to live if we end up building out more endpoints

	// Create an API instance
	api := api.NewAPI(node)

	// Set up your Gin router
	router := gin.Default()

	// Add the /peers endpoint
	router.GET("/peers", api.GetPeersHandler())

	// Add the /peerAddresses endpoint
	router.GET("/peerAddresses", api.GetPeerAddresses())

	// Start the server
	go router.Run() // By default it serves on :8080

	<-ctx.Done()
}

func getPort(name string) int {
	valueStr := os.Getenv(name)
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return 0
}

// getEnvAsBool will return the environment variable as a boolean or the default value
func getEnvAsBool(name string, defaultVal bool) bool {
	valueStr := os.Getenv(name)
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultVal
}

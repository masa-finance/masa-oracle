// cli.go
package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"math/big"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"

	masa "github.com/masa-finance/masa-oracle/pkg"
	masaCrypto "github.com/masa-finance/masa-oracle/pkg/crypto"
	masaStaking "github.com/masa-finance/masa-oracle/pkg/staking"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Bootnodes []string `json:"bootnodes"`
}

const (
	// Define the Ethereum node endpoint as a constant
	ethNodeEndpoint = "https://rpc.sepolia.org" // Sepolia endpoint
)

var (
	configFile    string
	start         bool
	portNbr       int
	udp           bool
	tcp           bool
	signature     string
	bootnodes     string
	flagBootnodes string
	data          string
	stakeAmount   string
)

func init() {
	// Define flags
	flag.StringVar(&configFile, "config", "config.json", "Path to the config file")
	flag.BoolVar(&start, "start", false, "Start flag")
	flag.IntVar(&portNbr, "port", getPort("portNbr"), "The port number")
	flag.BoolVar(&udp, "udp", getEnvAsBool("UDP", false), "UDP flag")
	flag.BoolVar(&tcp, "tcp", getEnvAsBool("TCP", false), "TCP flag")
	flag.StringVar(&signature, "signature", "", "The signature from the staking contract")
	flag.StringVar(&flagBootnodes, "bootnodes", "", "Comma-separated list of bootnodes")
	flag.StringVar(&data, "data", "", "The data to verify the signature against")
	flag.StringVar(&stakeAmount, "stake", "", "Amount of tokens to stake")
	flag.Parse()

	if start {
		if flagBootnodes != "" {
			bootnodes = flagBootnodes
		} else {
			config, err := loadConfig(configFile)
			if err != nil {
				logrus.Fatal(err)
			}
			bootnodes = strings.Join(config.Bootnodes, ",")
		}
	}

	err := os.Setenv(masa.Peers, bootnodes)
	if err != nil {
		logrus.Error(err)
	}
	//if neither udp nor tcp are set, default to udp
	if !udp && !tcp {
		udp = true
	}

	// New code to handle staking
	if stakeAmount != "" {
		amount, ok := new(big.Int).SetString(stakeAmount, 10)
		if !ok {
			logrus.Fatal("Invalid stake amount")
		}

		// Retrieve the current user's home directory
		usr, err := user.Current()
		if err != nil {
			logrus.Fatal("Failed to get user's home directory:", err)
		}

		// Construct the path to the private key file within the .masa directory
		keyFilePath := filepath.Join(usr.HomeDir, ".masa", "masa_oracle_key")

		// Retrieve or create the private key using the GetOrCreatePrivateKey function
		_, ecdsaPrivateKey, err := masaCrypto.GetOrCreatePrivateKey(keyFilePath)
		if err != nil {
			logrus.Fatal(err)
		}

		// Use the constant ethNodeEndpoint for the Ethereum node endpoint
		stakingClient, err := masaStaking.NewStakingClient(ethNodeEndpoint, ecdsaPrivateKey)
		if err != nil {
			logrus.Fatal(err)
		}

		receipt, err := stakingClient.Stake(amount)
		if err != nil {
			logrus.Fatal(err)
		}

		logrus.Infof("Stake transaction receipt: %v", receipt)
	}
}

func loadConfig(file string) (*Config, error) {
	var config Config
	configFile, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(configFile, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
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

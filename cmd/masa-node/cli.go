// cli.go
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	masa "github.com/masa-finance/masa-oracle/pkg"
	masaCrypto "github.com/masa-finance/masa-oracle/pkg/crypto"
	masaStaking "github.com/masa-finance/masa-oracle/pkg/staking"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Bootnodes []string `json:"bootnodes"`
}

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
	flag.BoolVar(&udp, "udp", false, "UDP flag") // Default value set to false
	flag.BoolVar(&tcp, "tcp", false, "TCP flag") // Default value set to false
	flag.StringVar(&signature, "signature", "", "The signature from the staking contract")
	flag.StringVar(&flagBootnodes, "bootnodes", "", "Comma-separated list of bootnodes")
	flag.StringVar(&data, "data", "", "The data to verify the signature against")
	flag.StringVar(&stakeAmount, "stake", "", "Amount of tokens to stake")
	flag.Parse()

	// Staking logic
	if stakeAmount != "" {
		amount, ok := new(big.Int).SetString(stakeAmount, 10)
		if !ok {
			logrus.Fatal("Invalid stake amount")
		}

		usr, err := user.Current()
		if err != nil {
			logrus.Fatal("Failed to get user's home directory:", err)
		}

		keyFilePath := filepath.Join(usr.HomeDir, ".masa", "masa_oracle_key")

		_, ecdsaPrivateKey, err := masaCrypto.GetOrCreatePrivateKey(keyFilePath)
		if err != nil {
			logrus.Fatal(err)
		}

		stakingClient, err := masaStaking.NewStakingClient(ecdsaPrivateKey)
		if err != nil {
			logrus.Fatal(err)
		}

		// Function to start and stop a spinner with a message
		startSpinner := func(msg string, txHash *string, done chan bool) {
			go func() {
				spinner := []string{"|", "/", "-", "\\"}
				i := 0
				for {
					select {
					case <-done:
						if txHash != nil {
							fmt.Printf("\r%s %s\n", msg, *txHash) // Print final message with txHash when done
						} else {
							fmt.Printf("\r%s\n", msg) // Print final message without txHash when done
						}
						return
					default:
						if txHash != nil {
							fmt.Printf("\r%s %s %s", spinner[i], msg, *txHash)
						} else {
							fmt.Printf("\r%s %s", spinner[i], msg)
						}
						i = (i + 1) % len(spinner)
						time.Sleep(100 * time.Millisecond)
					}
				}
			}()
		}

		// Approve the staking contract to spend tokens on behalf of the user
		var approveTxHash string
		done := make(chan bool)
		startSpinner("Approving staking contract to spend tokens... TxHash:", &approveTxHash, done)
		approveTxHash, err = stakingClient.Approve(amount)
		done <- true // Stop the spinner
		if err != nil {
			logrus.Fatal("Failed to approve tokens for staking:", err)
		}
		color.Green("Approve transaction hash: %s", approveTxHash)

		// Stake the tokens after approval
		var stakeTxHash string
		done = make(chan bool)
		startSpinner("Staking tokens... TxHash:", &stakeTxHash, done)
		stakeTxHash, err = stakingClient.Stake(amount)
		done <- true // Stop the spinner
		if err != nil {
			logrus.Fatal("Failed to stake tokens:", err)
		}
		color.Green("Stake transaction hash: %s", stakeTxHash)

		// Exit after staking, do not proceed to start the node
		os.Exit(0)
	}

	// Node startup logic
	if start {
		// Set the UDP and TCP flags based on environment variables if not already set
		if !udp {
			udp = getEnvAsBool("UDP", false)
		}
		if !tcp {
			tcp = getEnvAsBool("TCP", false)
		}

		if flagBootnodes != "" {
			bootnodes = flagBootnodes
		} else {
			config, err := loadConfig(configFile)
			if err != nil {
				logrus.Fatal(err)
			}
			bootnodes = strings.Join(config.Bootnodes, ",")
		}

		err := os.Setenv(masa.Peers, bootnodes)
		if err != nil {
			logrus.Error(err)
		}

		// Additional node startup logic here...
		// This is where you would start the node
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

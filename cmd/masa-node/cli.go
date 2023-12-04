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
	"github.com/sirupsen/logrus"

	masa "github.com/masa-finance/masa-oracle/pkg"
	masaCrypto "github.com/masa-finance/masa-oracle/pkg/crypto"
	masaStaking "github.com/masa-finance/masa-oracle/pkg/staking"
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
	debug         bool
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
	flag.BoolVar(&debug, "debug", false, "Override some protections for debugging (temporary)")
	flag.Parse()

	// Staking logic
	if stakeAmount != "" {
		// Convert the stake amount to the smallest unit, assuming 18 decimal places
		amountBigInt, ok := new(big.Int).SetString(stakeAmount, 10)
		if !ok {
			logrus.Fatal("Invalid stake amount")
		}
		amountInSmallestUnit := new(big.Int).Mul(amountBigInt, big.NewInt(1e18))

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
		startSpinner := func(msg string, txHashChan <-chan string, done chan bool) {
			spinner := []string{"|", "/", "-", "\\"}
			i := 0
			var txHash string
			for {
				select {
				case txHash = <-txHashChan: // Receive the transaction hash
					// Do not print anything here, just update the txHash variable
				case <-done:
					fmt.Printf("\r%s\n", msg) // Print final message when done
					if txHash != "" {
						fmt.Println(txHash) // Print the transaction hash on a new line
					}
					return
				default:
					// Use carriage return `\r` to overwrite the spinner animation on the same line
					// Remove the newline character `\n` from the print statement
					if txHash != "" {
						fmt.Printf("\r%s %s - %s", spinner[i], msg, txHash)
					} else {
						fmt.Printf("\r%s %s", spinner[i], msg)
					}
					i = (i + 1) % len(spinner)
					time.Sleep(100 * time.Millisecond)
				}
			}
		}

		// Approve the staking contract to spend tokens on behalf of the user
		var approveTxHash string
		done := make(chan bool)
		txHashChan := make(chan string, 1) // Buffer of 1 to prevent blocking
		go startSpinner("Approving staking contract to spend tokens...", txHashChan, done)
		approveTxHash, err = stakingClient.Approve(amountInSmallestUnit)
		if err != nil {
			logrus.Fatal("Failed to approve tokens for staking:", err)
		}
		txHashChan <- approveTxHash // Send the transaction hash to the spinner
		done <- true                // Stop the spinner
		color.Green("Approve transaction hash: %s", approveTxHash)

		// Stake the tokens after approval
		var stakeTxHash string
		done = make(chan bool)
		txHashChan = make(chan string, 1) // Buffer of 1 to prevent blocking
		go startSpinner("Staking tokens...", txHashChan, done)
		stakeTxHash, err = stakingClient.Stake(amountInSmallestUnit)
		if err != nil {
			logrus.Fatal("Failed to stake tokens:", err)
		}
		txHashChan <- stakeTxHash // Send the transaction hash to the spinner
		done <- true              // Stop the spinner
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

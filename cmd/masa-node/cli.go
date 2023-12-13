// cli.go
package main

import (
	"encoding/json"
	"flag"
	"os"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"

	masa "github.com/masa-finance/masa-oracle/pkg"
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
	}
}

func loadConfig(file string) (*Config, error) {
	var config Config
	configFile, err := os.ReadFile(file)
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

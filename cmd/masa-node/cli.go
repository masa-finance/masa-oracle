// cli.go
package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	masa "github.com/masa-finance/masa-oracle/pkg"
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

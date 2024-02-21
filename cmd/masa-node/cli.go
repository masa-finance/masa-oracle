// cli.go
package main

import (
	"flag"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	masa "github.com/masa-finance/masa-oracle/pkg"
)

var (
	portNbr     int
	udp         bool
	tcp         bool
	signature   string
	bootnodes   string
	data        string
	stakeAmount string
	debug       bool
	env         string
)

func init() {
	// Define flags
	flag.IntVar(&portNbr, "port", viper.GetInt("PORT_NBR"), "The port number")
	flag.BoolVar(&udp, "udp", viper.GetBool("UDP"), "UDP flag") // Default value set to false
	flag.BoolVar(&tcp, "tcp", viper.GetBool("TCP"), "TCP flag") // Default value set to false
	flag.StringVar(&signature, "signature", "", "The signature from the staking contract")
	flag.StringVar(&bootnodes, "bootnodes", viper.GetString("BOOTNODES"), "Comma-separated list of bootnodes")
	flag.StringVar(&data, "data", "", "The data to verify the signature against")
	flag.StringVar(&stakeAmount, "stake", viper.GetString("STAKE_AMOUNT"), "Amount of tokens to stake")
	flag.BoolVar(&debug, "debug", viper.GetBool("LOG_LEVEL"), "Override some protections for debugging (temporary)")
	flag.StringVar(&env, "env", viper.GetString("ENV"), "Environment to connect to")
	flag.Parse()

	err := os.Setenv(masa.Environment, env)
	if err != nil {
		logrus.Error(err)
	}
	err = os.Setenv(masa.Peers, bootnodes)
	if err != nil {
		logrus.Error(err)
	}

}

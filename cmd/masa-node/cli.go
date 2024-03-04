// cli.go
package main

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
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
	allowedPeer bool
)

func init() {
	// Define flags
	pflag.IntVar(&portNbr, "port", viper.GetInt("PORT_NBR"), "The port number")
	pflag.BoolVar(&udp, "udp", viper.GetBool("UDP"), "UDP flag") // Default value set to false
	pflag.BoolVar(&tcp, "tcp", viper.GetBool("TCP"), "TCP flag") // Default value set to false
	pflag.StringVar(&signature, "signature", "", "The signature from the staking contract")
	pflag.StringVar(&bootnodes, "bootnodes", viper.GetString("BOOTNODES"), "Comma-separated list of bootnodes")
	pflag.StringVar(&data, "data", "", "The data to verify the signature against")
	pflag.StringVar(&stakeAmount, "stake", viper.GetString("STAKE_AMOUNT"), "Amount of tokens to stake")
	pflag.BoolVar(&debug, "debug", viper.GetBool("LOG_LEVEL"), "Override some protections for debugging (temporary)")
	pflag.StringVar(&env, "env", viper.GetString("ENV"), "Environment to connect to")
	pflag.BoolVar(&allowedPeer, "allowedPeer", viper.GetBool("allowedPeer"), "Set to true to allow setting this node as the allowed peer")
	pflag.Parse()
	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		logrus.Error(err)
		return
	}
	if debug {
		logrus.Debugf("Direct flag value: %s", env)
		for key, val := range viper.AllSettings() {
			logrus.Debugf("%s: %v", key, val)
		}
		logrus.Debugf("ENV from Viper: %s", viper.Get("ENV"))
	}
}

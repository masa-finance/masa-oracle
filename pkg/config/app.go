package config

import (
	"log"
	"os/user"
	"path/filepath"
	"reflect"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// AppConfig is designed as a singleton to ensure that there is only one instance
// of the configuration throughout the application. This design pattern is useful
// for managing global application settings, allowing various parts of the application
// to access configuration settings consistently without the need to pass the configuration
// object around or risk having multiple, potentially conflicting instances of the configuration.
//
// The singleton pattern is implemented using a combination of a private instance variable
// (`instance`) and a public `GetInstance` method. The `instance` variable holds the single
// instance of AppConfig, while the `GetInstance` method provides a global access point to that instance.
// Additionally, the `sync.Once` mechanism ensures that the AppConfig instance is initialized only once,
// making the initialization thread-safe.
//
// Usage:
// To access the AppConfig instance, call the GetInstance method from anywhere in your application:
//
//     config := config.GetInstance()
//
// This call will return the singleton AppConfig instance. If the instance has not been initialized yet,
// `GetInstance` will initialize it by setting default values, reading configuration from files,
// environment variables, and command-line flags, and then return the instance. Subsequent calls to
// `GetInstance` will return the same instance without reinitializing it.
//
// It's important to note that since AppConfig is a singleton, any modifications to the configuration
// settings through the AppConfig instance will be reflected across the entire application.

var (
	instance *AppConfig
	once     sync.Once
)

type AppConfig struct {
	PortNbr              int      `mapstructure:"portNbr"`
	UDP                  bool     `mapstructure:"udp"`
	TCP                  bool     `mapstructure:"tcp"`
	PrivateKey           string   `mapstructure:"privateKey"`
	Signature            string   `mapstructure:"signature"`
	Bootnodes            []string `mapstructure:"bootnodes"`
	Data                 string   `mapstructure:"data"`
	StakeAmount          string   `mapstructure:"stakeAmount"`
	Debug                bool     `mapstructure:"debug"`
	Environment          string   `mapstructure:"env"`
	PrivateKeyFile       string   `mapstructure:"privateKeyFile"`
	MasaDir              string   `mapstructure:"masaDir"`
	RpcUrl               string   `mapstructure:"rpcUrl"`
	AllowedPeer          bool     `mapstructure:"allowedPeer"`
	AllowedPeerId        string   `mapstructure:"allowedPeerId"`
	AllowedPeerPublicKey string   `mapstructure:"allowedPeerPublicKey"`
	LogLevel             string   `mapstructure:"logLevel"`
	LogFilePath          string   `mapstructure:"logFilePath"`
	DbPath               string   `mapstructure:"dbPath"`
}

func GetInstance() *AppConfig {
	once.Do(func() {
		instance = &AppConfig{}

		instance.setDefaultConfig()
		instance.setEnvVariableConfig()
		instance.setFileConfig()
		err := instance.setCommandLineConfig()
		if err != nil {
			logrus.Fatal(err)
		}

		err = viper.Unmarshal(instance)
		if err != nil {
			logrus.Errorf("Unable to unmarshal config into struct, %v", err)
			instance = nil // Ensure instance is nil if unmarshalling fails
		}
	})
	return instance
}

func (c *AppConfig) setDefaultConfig() {
	usr, err := user.Current()
	if err != nil {
		log.Fatal("could not find user.home directory")
	}
	viper.SetDefault(MasaDir, filepath.Join(usr.HomeDir, ".masa"))
	viper.SetDefault(RpcUrl, "https://ethereum-sepolia.publicnode.com")
	viper.SetDefault(PortNbr, "4001")
	viper.SetDefault(UDP, true)
	viper.SetDefault(TCP, false)
	viper.SetDefault(StakeAmount, "")
	viper.SetDefault(AllowedPeer, false)
	viper.SetDefault(LogLevel, "info")
	viper.SetDefault(LogFilePath, "masa_oracle_node.log")
	viper.SetDefault(PrivKeyFile, filepath.Join(viper.GetString(MasaDir), "masa_oracle_key"))
	viper.SetDefault(DbPath, filepath.Join(viper.GetString(MasaDir), "masa-node-db"))

}

// TODO: add a variable to allow for the config file location to be set
func (c *AppConfig) setFileConfig() {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".") // Optionally: add other paths, e.g., home directory or etc

	// Attempt to read the config file
	if err := viper.ReadInConfig(); err != nil {
		logrus.Warnf("Error reading config file: %s", err)
	}
}

func (c *AppConfig) setEnvVariableConfig() {
	viper.AutomaticEnv()
}

func (c *AppConfig) setCommandLineConfig() error {
	var bootnodes string
	pflag.IntVar(&c.PortNbr, "port", viper.GetInt(PortNbr), "The port number")
	pflag.BoolVar(&c.UDP, "udp", viper.GetBool(UDP), "UDP flag")
	pflag.BoolVar(&c.TCP, "tcp", viper.GetBool(TCP), "TCP flag")
	pflag.StringVar(&bootnodes, "bootnodes", viper.GetString(BootNodes), "Comma-separated list of bootnodes")
	pflag.StringVar(&c.StakeAmount, "stake", viper.GetString(StakeAmount), "Amount of tokens to stake")
	pflag.BoolVar(&c.Debug, "debug", viper.GetBool(Debug), "Override some protections for debugging (temporary)")
	pflag.StringVar(&c.Environment, "env", viper.GetString(Environment), "Environment to connect to")
	pflag.BoolVar(&c.AllowedPeer, "allowedPeer", viper.GetBool(AllowedPeer), "Set to true to allow setting this node as the allowed peer")
	pflag.StringVar(&c.PrivateKey, "privateKey", viper.GetString(PrivateKey), "The private key")
	pflag.StringVar(&c.PrivateKeyFile, PrivKeyFile, viper.GetString(PrivKeyFile), "The private key file")
	pflag.StringVar(&c.MasaDir, MasaDir, viper.GetString(MasaDir), "The masa directory")
	pflag.StringVar(&c.RpcUrl, RpcUrl, viper.GetString(RpcUrl), "The RPC URL")
	pflag.StringVar(&c.Signature, Signature, viper.GetString(Signature), "The signature from the staking contract")
	pflag.StringVar(&c.Data, "data", viper.GetString("data"), "The data to verify the signature against")
	pflag.StringVar(&c.LogLevel, LogLevel, viper.GetString(LogLevel), "The log level")
	pflag.StringVar(&c.LogFilePath, LogFilePath, viper.GetString(LogFilePath), "The log file path")
	pflag.StringVar(&c.DbPath, DbPath, viper.GetString(DbPath), "The badger database path")
	pflag.Parse()

	// Bind command line flags to Viper (optional, if you want to use Viper for additional configuration)
	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		return err
	}
	c.Bootnodes = strings.Split(bootnodes, ",")
	return nil
}

// LogConfig logs the non-sensitive parts of the AppConfig.
// Adjust the fields being logged according to your application's needs.
func (c *AppConfig) LogConfig() {
	val := reflect.ValueOf(*c)
	typeOfStruct := val.Type()

	logrus.Info("Current AppConfig values:")
	for i := 0; i < val.NumField(); i++ {
		field := typeOfStruct.Field(i)
		value := val.Field(i).Interface()

		// Example of skipping sensitive fields
		if field.Name == "PrivateKeyFile" || field.Name == "Signature" {
			continue
		}
		logrus.Infof("%s: %v", field.Name, value)
	}
}

func (c *AppConfig) HasBootnodes() bool {
	return c.Bootnodes[0] != ""
}
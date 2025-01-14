package config

import (
	"fmt"
	"log"
	"os/user"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/masa-finance/masa-oracle/internal/versioning"
	"github.com/masa-finance/masa-oracle/pkg/masacrypto"

	"github.com/gotd/contrib/bg"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// AppConfig represents the configuration settings for the application.
// It holds various parameters and settings that control the behavior and runtime environment of the application.
// Most of the fields in this struct are tagged with `mapstructure` to facilitate configuration loading from various
// sources such as configuration files, environment variables, and command-line flags using the `viper` library.
type AppConfig struct {
	Version              string   `mapstructure:"version"`
	PortNbr              int      `mapstructure:"portNbr"`
	APIListenAddress     string   `mapstructure:"api_listen_address"`
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
	FilePath             string   `mapstructure:"filePath"`
	Validator            bool     `mapstructure:"validator"`
	CachePath            string   `mapstructure:"cachePath"`
	Faucet               bool     `mapstructure:"faucet"`

	// These may be moved to a separate struct
	TwitterCookiesPath string `mapstructure:"twitterCookiesPath"`
	TwitterUsername    string `mapstructure:"twitterUsername"`
	TwitterPassword    string `mapstructure:"twitterPassword"`
	Twitter2FaCode     string `mapstructure:"twitter2FaCode"`
	DiscordBotToken    string `mapstructure:"discordBotToken"`
	TwitterScraper     bool   `mapstructure:"twitterScraper"`
	DiscordScraper     bool   `mapstructure:"discordScraper"`
	TelegramScraper    bool   `mapstructure:"telegramScraper"`
	WebScraper         bool   `mapstructure:"webScraper"`
	APIEnabled         bool   `mapstructure:"api_enabled"`
	TunnelEnabled      bool   `mapstructure:"tunnel_enabled"`
	TunnelListenAddr   string `mapstructure:"tunnelListenAddr"`
	TunnelTargetPort   uint16 `mapstructure:"tunnelTargetPort"`

	// Key: port number, value: peerId
	// Note that the TUNNEL_PORTS env var has to be JSON (thanks, Viper)
	TunnelPorts map[string]string `mapstructure:"tunnelPorts"`

	KeyManager   *masacrypto.KeyManager
	TelegramStop bg.StopFunc
}

// GetConfig parses and fills in the AppConfig. The configuration values are generated
// by the following steps in order:
//
// 1. Set default configuration values.
// 2. Override with values from environment variables.
// 3. Override with values from the configuration file.
// 4. Override with command-line flags.
//
// In case of any errors it returns nil with an error object
func GetConfig() (*AppConfig, error) {
	instance := &AppConfig{}
	instance.setDefaultConfig()
	instance.setEnvVariableConfig()
	instance.setFileConfig(viper.GetString("FILE_PATH"))

	if err := instance.setCommandLineConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(instance); err != nil {
		return nil, fmt.Errorf("unable to unmarshal config into struct, %v", err)
	}

	if instance.PrivateKeyFile == DefaultPrivKeyFile {
		instance.PrivateKeyFile = filepath.Join(instance.MasaDir, DefaultPrivKeyFile)
	}
	instance.APIEnabled = viper.GetBool("api_enabled")

	keyManager, err := masacrypto.NewKeyManager(instance.PrivateKey, instance.PrivateKeyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize keys: %v", err)
	}
	instance.KeyManager = keyManager

	return instance, nil
}

// setDefaultConfig sets the default configuration values for the AppConfig instance.
// It retrieves the user's home directory and sets default values for various configuration options
// such as the MasaDir, Bootnodes, RpcUrl, Environment, FilePath, Validator, and CachePath.
func (c *AppConfig) setDefaultConfig() {

	usr, err := user.Current()
	if err != nil {
		log.Fatal("could not find user.home directory")
	}

	// Set defaults
	viper.SetDefault(MasaDir, filepath.Join(usr.HomeDir, ".masa"))

	// Set defaults
	viper.SetDefault("version", versioning.ProtocolVersion)
	viper.SetDefault(PortNbr, "4001")
	viper.SetDefault(APIListenAddress, "127.0.0.1:8080")
	viper.SetDefault(UDP, true)
	viper.SetDefault(TCP, false)
	viper.SetDefault(StakeAmount, "")
	viper.SetDefault(Faucet, false)
	viper.SetDefault(AllowedPeer, true)
	viper.SetDefault(LogLevel, "info")
	viper.SetDefault(LogFilePath, "masa_node.log")
	// It should be ${MASA_DIR}/masa_oracle_key. However, we really don't know the value of MASA_DIR at this point.
	// We set it up with a marker value which we will have to fix once we load the full config.
	viper.SetDefault(PrivKeyFile, DefaultPrivKeyFile)

	viper.SetDefault(APIEnabled, false)

	viper.SetDefault(TunnelEnabled, true)
	viper.SetDefault(TunnelListenAddr, "127.0.0.1")
	viper.SetDefault(TunnelTargetPort, "26656")

	viper.SetDefault(TunnelPorts, "")
}

// setFileConfig loads configuration from a YAML file.
// It takes the file path as a parameter and sets up Viper to read the config file.
// If the config file exists, it will be read into the AppConfig struct.
func (c *AppConfig) setFileConfig(path string) {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(path) // Optionally: add other paths, e.g., home directory or etc

	// Attempt to read the config file if exists
	_ = viper.ReadInConfig()
}

// setEnvVariableConfig loads environment variables into the AppConfig struct.
// It reads the .env file using the godotenv package and automatically binds
// environment variables to Viper for configuration management.
func (c *AppConfig) setEnvVariableConfig() {
	err := godotenv.Load()
	if err != nil {
		logrus.Error("[-] Error loading .env file")
	}
	viper.AutomaticEnv()
}

// setCommandLineConfig parses command line flags and binds them to the AppConfig struct.
// It takes no parameters and returns an error if there is an issue binding the flags.
// The function sets up command line flags for various configuration options using the pflag package.
// After parsing the flags, it binds them to Viper for additional configuration management.
// Finally, it splits the 'bootnodes' flag value into a slice and assigns it to the Bootnodes field of the AppConfig struct.
func (c *AppConfig) setCommandLineConfig() error {
	var bootnodes string
	pflag.IntVar(&c.PortNbr, "port", viper.GetInt(PortNbr), "The port number")
	pflag.BoolVar(&c.UDP, "udp", viper.GetBool(UDP), "UDP flag")
	pflag.BoolVar(&c.TCP, "tcp", viper.GetBool(TCP), "TCP flag")
	pflag.StringVar(&bootnodes, "bootnodes", viper.GetString(BootNodes), "Comma-separated list of bootnodes")
	pflag.StringVar(&c.StakeAmount, "stake", viper.GetString(StakeAmount), "Amount of tokens to stake")
	pflag.BoolVar(&c.Debug, "debug", viper.GetBool(Debug), "Override some protections for debugging (temporary)")
	pflag.StringVar(&c.Environment, "env", viper.GetString(Environment), "Environment to connect to")
	pflag.StringVar(&c.Version, "version", viper.GetString("version"), "Application version")
	pflag.BoolVar(&c.AllowedPeer, "allowedPeer", viper.GetBool(AllowedPeer), "Set to true to allow setting this node as the allowed peer")
	pflag.StringVar(&c.PrivateKey, "privateKey", viper.GetString(PrivateKey), "The private key")
	pflag.StringVar(&c.PrivateKeyFile, "privKeyFile", viper.GetString(PrivKeyFile), "The private key file")
	pflag.StringVar(&c.MasaDir, "masaDir", viper.GetString(MasaDir), "The masa directory")
	pflag.StringVar(&c.RpcUrl, "rpcUrl", viper.GetString(RpcUrl), "The RPC URL")
	pflag.StringVar(&c.Signature, "signature", viper.GetString(Signature), "The signature from the staking contract")
	pflag.StringVar(&c.Data, "data", viper.GetString("data"), "The data to verify the signature against")
	pflag.StringVar(&c.LogLevel, "logLevel", viper.GetString(LogLevel), "The log level")
	pflag.StringVar(&c.LogFilePath, "logFilePath", viper.GetString(LogFilePath), "The log file path")
	pflag.StringVar(&c.FilePath, "filePath", viper.GetString(FilePath), "The node file path")
	pflag.BoolVar(&c.Validator, "validator", viper.GetBool(Validator), "Approved validator node boolean")
	pflag.StringVar(&c.CachePath, "cachePath", viper.GetString(CachePath), "The cache path")
	pflag.StringVar(&c.TwitterUsername, "twitterUsername", viper.GetString(TwitterUsername), "Twitter Username")
	pflag.StringVar(&c.TwitterPassword, "twitterPassword", viper.GetString(TwitterPassword), "Twitter Password")
	pflag.StringVar(&c.Twitter2FaCode, "twitter2FaCode", viper.GetString(Twitter2FaCode), "Twitter 2FA Code")
	pflag.StringVar(&c.DiscordBotToken, "discordBotToken", viper.GetString(DiscordBotToken), "Discord Bot Token")
	pflag.BoolVar(&c.TwitterScraper, "twitterScraper", viper.GetBool(TwitterScraper), "Twitter Scraper")
	pflag.BoolVar(&c.DiscordScraper, "discordScraper", viper.GetBool(DiscordScraper), "Discord Scraper")
	pflag.BoolVar(&c.TelegramScraper, "telegramScraper", viper.GetBool(TelegramScraper), "Telegram Scraper")
	pflag.BoolVar(&c.WebScraper, "webScraper", viper.GetBool(WebScraper), "Web Scraper")
	pflag.BoolVar(&c.Faucet, "faucet", viper.GetBool(Faucet), "Faucet")
	pflag.BoolVar(&c.APIEnabled, "api-enabled", viper.GetBool(APIEnabled), "Enable API server")
	pflag.StringVar(&c.APIListenAddress, "api-port", viper.GetString(APIListenAddress), "API Listening address")
	pflag.BoolVar(&c.TunnelEnabled, "tunnel-enabled", viper.GetBool(TunnelEnabled), "Enable tunnel")
	pflag.StringVar(&c.TunnelListenAddr, "tunnelListenAddr", viper.GetString(TunnelListenAddr), "Tunnel listen address")
	pflag.Uint16Var(&c.TunnelTargetPort, "tunnelTargetPort", viper.GetUint16(TunnelTargetPort), "Tunnel target port")
	pflag.StringToStringVarP(&c.TunnelPorts, "tunnelPorts", "t", viper.GetStringMapString(TunnelPorts), "Tunnel ports, use `port=peerId`, can specify multiple times")

	pflag.Parse()

	// Bind command line flags to Viper
	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		return err
	}

	// Add this line after binding flags
	viper.Set("api_enabled", c.APIEnabled)

	c.Bootnodes = strings.Split(bootnodes, ",")

	return nil
}

// LogConfig logs the non-sensitive parts of the AppConfig.
// Adjust the fields being logged according to your application's needs.
func (c *AppConfig) LogConfig() {
	val := reflect.ValueOf(*c)
	typeOfStruct := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := typeOfStruct.Field(i)
		value := val.Field(i).Interface()

		// Skipping sensitive fields
		if field.Name == "PrivateKey" || field.Name == "Signature" || field.Name == "PrivateKeyFile" {
			continue
		}
		logrus.Infof("%s: %v", field.Name, value)
	}
}

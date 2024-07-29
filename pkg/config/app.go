package config

import (
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"reflect"
	"strings"
	"sync"

	"github.com/joho/godotenv"
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

// AppConfig represents the configuration settings for the application.
// It holds various parameters and settings that control the behavior and runtime environment of the application.
// The fields in this struct are tagged with `mapstructure` to facilitate configuration loading from various sources
// such as configuration files, environment variables, and command-line flags using the `viper` library.
type AppConfig struct {
	Version              string   `mapstructure:"version"`
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
	ClaudeApiKey       string `mapstructure:"claudeApiKey"`
	ClaudeApiURL       string `mapstructure:"claudeApiURL"`
	ClaudeApiVersion   string `mapstructure:"claudeApiVersion"`
	GPTApiKey          string `mapstructure:"gptApiKey"`
	TwitterScraper     bool   `mapstructure:"twitterScraper"`
	DiscordScraper     bool   `mapstructure:"discordScraper"`
	WebScraper         bool   `mapstructure:"webScraper"`
	LlmServer          bool   `mapstructure:"llmServer"`
	LLMChatUrl         string `mapstructure:"llmChatUrl"`
	LLMCfUrl           string `mapstructure:"llmCfUrl"`
}

// GetInstance returns the singleton instance of AppConfig.
//
// If the instance has not been initialized yet, GetInstance will initialize it by:
// 1. Creating a new AppConfig instance.
// 2. Setting default configuration values.
// 3. Overriding defaults with values from environment variables.
// 4. Overriding defaults and environment variables with values from the configuration file.
// 5. Overriding all previous values with command-line flags.
// 6. Unmarshalling the configuration into the AppConfig instance.
//
// If the unmarshalling fails, the instance is set to nil.
//
// Subsequent calls to GetInstance will return the same initialized instance.
func GetInstance() *AppConfig {
	once.Do(func() {
		instance = &AppConfig{}

		cmd := exec.Command("git", "describe", "--tags", "--abbrev=0")
		output, e := cmd.Output()
		if e != nil {
			logrus.Errorf("Error running git describe: %v\n", e)
			os.Exit(1)
		}
		logrus.Infof("output: %s", output)

		instance.setDefaultConfig()
		instance.Version = strings.TrimSpace(string(output))
		instance.setEnvVariableConfig()
		instance.setFileConfig(viper.GetString("FILE_PATH"))
		err := instance.setCommandLineConfig()
		if err != nil {
			logrus.Fatal(err)
		}

		err = viper.Unmarshal(instance)
		if err != nil {
			logrus.Errorf("[-] Unable to unmarshal config into struct, %v", err)
			instance = nil // Ensure instance is nil if unmarshalling fails
		}

	})
	return instance
}

// setDefaultConfig sets the default configuration values for the AppConfig instance.
// It retrieves the user's home directory and sets default values for various configuration options
// such as the MasaDir, Bootnodes, RpcUrl, Environment, FilePath, Validator, and CachePath.
// It also fetches bootnode information from a remote URL based on the environment (dev, test, or main).
func (c *AppConfig) setDefaultConfig() {

	usr, err := user.Current()
	if err != nil {
		log.Fatal("could not find user.home directory")
	}

	// Set defaults
	viper.SetDefault(MasaDir, filepath.Join(usr.HomeDir, ".masa"))

	// Set defaults
	viper.SetDefault("Version", Version)
	viper.SetDefault(PortNbr, "4001")
	viper.SetDefault(UDP, true)
	viper.SetDefault(TCP, false)
	viper.SetDefault(StakeAmount, "")
	viper.SetDefault(Faucet, false)
	viper.SetDefault(AllowedPeer, true)
	viper.SetDefault(LogLevel, "info")
	viper.SetDefault(LogFilePath, "masa_node.log")
	viper.SetDefault(PrivKeyFile, filepath.Join(viper.GetString(MasaDir), "masa_oracle_key"))
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
	pflag.StringVar(&c.ClaudeApiKey, "claudeApiKey", viper.GetString(ClaudeApiKey), "Claude API Key")
	pflag.StringVar(&c.ClaudeApiURL, "claudeApiUrl", viper.GetString(ClaudeApiURL), "Claude API Url")
	pflag.StringVar(&c.ClaudeApiVersion, "claudeApiVersion", viper.GetString(ClaudeApiVersion), "Claude API Version")
	pflag.StringVar(&c.GPTApiKey, "gptApiKey", viper.GetString(GPTApiKey), "OpenAI API Key")
	pflag.StringVar(&c.LLMChatUrl, "llmChatUrl", viper.GetString(LlmChatUrl), "URL for support LLM Chat calls")
	pflag.StringVar(&c.LLMCfUrl, "llmCfUrl", viper.GetString(LlmCfUrl), "URL for support LLM Cloudflare calls")
	pflag.BoolVar(&c.TwitterScraper, "twitterScraper", viper.GetBool(TwitterScraper), "TwitterScraper")
	pflag.BoolVar(&c.DiscordScraper, "discordScraper", viper.GetBool(DiscordScraper), "DiscordScraper")
	pflag.BoolVar(&c.WebScraper, "webScraper", viper.GetBool(WebScraper), "WebScraper")
	pflag.BoolVar(&c.LlmServer, "llmServer", viper.GetBool(LlmServer), "Can service LLM requests")
	pflag.BoolVar(&c.Faucet, "faucet", viper.GetBool(Faucet), "Faucet")
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

	// logrus.Info("Current AppConfig values:")
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

// HasBootnodes checks if the AppConfig has any bootnodes configured.
// It returns true if there is at least one bootnode in the Bootnodes slice and it is not an empty string.
// Otherwise, it returns false, indicating that no bootnodes are configured.
func (c *AppConfig) HasBootnodes() bool {
	return c.Bootnodes[0] != ""
}

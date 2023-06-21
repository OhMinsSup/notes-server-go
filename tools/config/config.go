package config

import (
	"github.com/fatih/color"
	"github.com/spf13/viper"
)

const (
	DefaultServerRoot     = "localhost:8000"
	DefaultPort           = 8000
	ConfigFilePath        = "./config.json"
	RecordAuthTokenSecret = "1a2b3c4d5e6f7g8h9i1a2b3c4d5e6f7g8h9i1a2b3c4d5e6f7g8h9i"
)

type Configuration struct {
	ServerRoot            string `json:"serverRoot" mapstructure:"serverRoot"`
	Port                  int    `json:"port" mapstructure:"port"`
	IsDebug               bool   `json:"isDebug" mapstructure:"isDebug"`
	RecordAuthTokenSecret string `json:"recordAuthTokenSecret" mapstructure:"recordAuthTokenSecret"`
}

// ReadConfigFile reads the configuration file and returns the configuration.
func ReadConfigFile(configFilePath string) (*Configuration, error) {
	if configFilePath == "" {
		viper.SetConfigFile(ConfigFilePath)
	} else {
		viper.SetConfigFile(configFilePath)
	}

	viper.SetDefault("ServerRoot", DefaultServerRoot)
	viper.SetDefault("Port", DefaultPort)
	viper.SetDefault("IsDebug", true)
	viper.SetDefault("RecordAuthTokenSecret", RecordAuthTokenSecret)

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		return nil, err
	}

	configuration := Configuration{}

	err = viper.Unmarshal(&configuration)
	if err != nil {
		return nil, err
	}

	regular := color.New()
	regular.Printf(" ➜ ReadConfigFile: %+v\n", removeSecurityData(configuration))

	return &configuration, nil
}

func removeSecurityData(config Configuration) Configuration {
	clean := config
	return clean
}

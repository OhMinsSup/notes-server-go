package config

import (
	"github.com/fatih/color"
	"github.com/spf13/viper"
)

const (
	DefaultServerRoot = "localhost:8000"
	DefaultPort       = 8000
)

type Configuration struct {
	ServerRoot     string `json:"serverRoot" mapstructure:"serverRoot"`
	Port           int    `json:"port" mapstructure:"port"`
	DBConfigString string `json:"dbconfig" mapstructure:"dbconfig"`
	IsDebug          bool   `json:"isDebug" mapstructure:"isDebug"`
}

// ReadConfigFile reads the configuration file and returns the configuration.
func ReadConfigFile(configFilePath string) (*Configuration, error) {
	if configFilePath == "" {
		viper.SetConfigFile("./config.json")
	} else {
		viper.SetConfigFile(configFilePath)
	}

	viper.SetDefault("ServerRoot", DefaultServerRoot)
	viper.SetDefault("Port", DefaultPort)
	viper.SetDefault("DBConfigString", "./database.sqlite")

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
	regular.Printf(" ➜ RreadConfigFile: %+v", removeSecurityData(configuration))

	return &configuration, nil
}

func removeSecurityData(config Configuration) Configuration {
	clean := config
	return clean
}
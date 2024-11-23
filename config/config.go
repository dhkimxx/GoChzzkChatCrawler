package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type AppConfig struct {
	AccessToken string `yaml:"access_token"`
}

var Config AppConfig

func init() {
	_, err := os.Getwd()
	if err != nil {
		panic(fmt.Errorf("failed to get current working directory: %s", err))
	}
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("error reading config file, %s", err))
	}

	if err := viper.Unmarshal(&Config); err != nil {
		panic(fmt.Errorf("error unmarshaling config into struct, %s", err))
	}
}

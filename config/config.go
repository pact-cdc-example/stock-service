package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

func New() Manager {
	directory, _ := os.Getwd()
	configPath := fmt.Sprintf(".%s", path)
	if strings.HasSuffix(directory, "pact-cdc-test") {
		configPath = fmt.Sprintf("./stock-service%s", path)
	}

	viper.AddConfigPath(configPath)
	viper.SetConfigName(local)
	viper.SetConfigType(yaml)

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Sprintf("error while reading config file: %s", err))
	}

	var config config
	if err := viper.Unmarshal(&config); err != nil {
		panic(fmt.Sprintf("error while unmarshalling config file: %s", err))
	}

	global = &manager{config: &config}
	return global
}

const (
	path  = "/.config"
	local = "local"
	yaml  = "yaml"
)

func Global() Manager {
	if global == nil {
		return New()
	}

	return global
}

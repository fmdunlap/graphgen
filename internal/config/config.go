package config

import (
	"fmt"
	"github.com/caarlos0/env/v11"
)

var cachedConfig *Config

type Config struct {
	Auth     AuthConfig
	Database DatabaseConfig
	Server   ServerConfig
}

// ParseConfig parses the environment variables and returns a Config struct
func ParseConfig() *Config {
	if cachedConfig != nil {
		return cachedConfig
	}

	conf := Config{}
	err := env.Parse(&conf)
	if err != nil {
		fmt.Println("Error parsing initial config!", err)
		return nil
	}

	cachedConfig = &conf
	return cachedConfig
}

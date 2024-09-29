package config

import "time"

type ServerConfig struct {
	Address                string `env:"SERVER_ADDRESS" envDefault:"0.0.0.0"`
	Environment            string `env:"SERVER_ENVIRONMENT" envDefault:"dev"`
	Port                   int    `env:"SERVER_PORT" envDefault:"8080"`
	ShutdownTimeoutSeconds int    `env:"SERVER_SHUTDOWN_TIMEOUT" envDefault:"15"`
}

func (sc ServerConfig) GetShutdownTimeoutDuration() time.Duration {
	return time.Second * time.Duration(sc.ShutdownTimeoutSeconds)
}

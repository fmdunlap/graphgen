package config

type ServerConfig struct {
	Address     string `env:"SERVER_ADDRESS" envDefault:"0.0.0.0"`
	Environment string `env:"SERVER_ENVIRONMENT" envDefault:"dev"`
	Port        int    `env:"SERVER_PORT" envDefault:"8080"`
}

package config

type DatabaseConfig struct {
	Database string `env:"DB_DATABASE" envDefault:"postgres"`
	Host     string `env:"DB_HOST,required"`
	Password string `env:"DB_PASSWORD,required"`
	Port     int    `env:"DB_PORT" envDefault:"5432"`
	Schema   string `env:"DB_SCHEMA" envDefault:"public"`
	Username string `env:"DB_USERNAME,required"`
}

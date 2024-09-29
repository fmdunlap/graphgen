package config

type AuthConfig struct {
	PublicKey  string `env:"AUTH_PUBLIC_KEY,required"`
	PrivateKey string `env:"AUTH_PRIVATE_KEY,required"`
}

package config

type Config struct {
	ClientID string `env:"CLIENT_ID"`
	AuthURL  string `env:"AUTH_URL"`
	Address  string `mapstructure:"address"`
}

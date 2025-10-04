package config

import "time"

// TODO: validate method

type Config struct {
	Env        string
	HTTPServer HTTPServer `mapstructure:"http_server"`
	JWT        Jwt
}

type HTTPServer struct {
	Address  string
	Timeout  time.Duration
	APIToken string `env:"API_TOKEN"`
}

type Jwt struct {
	SecretKey  string        `env:"JWT_SECRET_KEY"`
	AccessTTL  time.Duration `mapstructure:"access_ttl"`
	RefreshTTL time.Duration `mapstructure:"refresh_ttl"`
}

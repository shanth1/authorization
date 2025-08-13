package authcfg

import "time"

// TODO: validate method

type Config struct {
	Env        string
	HTTPServer HTTPServer `mapstructure:"http_server"`
	Telegram   Telegram
	JWT        Jwt
}

type HTTPServer struct {
	Address  string
	Timeout  time.Duration
	APIToken string `env:"HTTP_SERVER_API_TOKEN"`
}

type Telegram struct {
	BotName string `env:"TELEGRAM_BOT_NAME"`
}

type Jwt struct {
	SecretKey  string        `env:"JWT_SECRET_KEY"`
	AccessTTL  time.Duration `mapstructure:"access_ttl"`
	RefreshTTL time.Duration `mapstructure:"refresh_ttl"`
}

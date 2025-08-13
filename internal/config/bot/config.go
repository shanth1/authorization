package botcfg

type Config struct {
	Env            string
	APIToken       string `env:"API_TOKEN"`
	BotToken       string `env:"TELEGRAM_BOT_TOKEN"`
	AuthServiceURL string `env:"AUTH_SERVICE_URL"`
}

// TODO: validate method

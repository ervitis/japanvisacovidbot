package telegram

import "github.com/kelseyhightower/envconfig"

var (
	Config TelegramConfigParameters
)

func LoadTelegramConfig() {
	envconfig.MustProcess("", &Config)
}

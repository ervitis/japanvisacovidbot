package telegram

import "github.com/kelseyhightower/envconfig"

var (
	TelegramConfig TelegramConfigParameters
)

func LoadTelegramConfig() {
	envconfig.MustProcess("", &TelegramConfig)
}

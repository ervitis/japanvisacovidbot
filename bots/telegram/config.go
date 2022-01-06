package telegram

import "github.com/kelseyhightower/envconfig"

var (
	Config ConfigParameters
)

func LoadTelegramConfig() {
	envconfig.MustProcess("", &Config)
}

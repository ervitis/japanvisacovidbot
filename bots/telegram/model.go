package telegram

type (
	TelegramConfigParameters struct {
		ID    int    `envconfig:"TELEGRAM_USERID"`
		Token string `envconfig:"TOKEN"`
	}
)

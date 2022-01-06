package telegram

type (
	ConfigParameters struct {
		ID    int    `envconfig:"TELEGRAM_USERID"`
		Token string `envconfig:"TOKEN"`
	}
)

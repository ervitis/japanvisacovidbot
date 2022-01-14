package email

import (
	"github.com/kelseyhightower/envconfig"
	"os"
)

type (
	Parameters struct {
		SmtpServer string `envconfig:"SMTP_SERVER"`
		Username   string `envconfig:"USERNAME"`
		Password   string `envconfig:"PASSWORD"`
		SmtpPort   string `envconfig:"SMTP_PORT"`
		To         string `envconfig:"TO"`
		From       string `envconfig:"FROM"`
		Subject    string `envconfig:"SUBJECT"`
		TestEmail  bool   `envconfig:"TEST" default:"true"`
	}
)

var (
	Config Parameters
)

func LoadConfig() {
	envconfig.MustProcess("EMAIL", &Config)

	if Config.TestEmail {
		Config.To = os.Getenv("EMAIL_FOR_TESTING")
		if Config.To == "" {
			panic("First set var EMAIL_FOR_TESTING")
		}
	}
}

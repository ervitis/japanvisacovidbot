package repo

import "github.com/kelseyhightower/envconfig"

type (
	DBConfigParameters struct {
		Password string `envconfig:"POSTGRES_PASSWORD"`
		User     string `envconfig:"POSTGRES_USER" default:"covid"`
		DB       string `envconfig:"POSTGRES_DB" default:"japancovid"`
		Host     string `envconfig:"POSTGRES_HOST" default:"localhost"`
		Port     int    `envconfig:"POSTGRES_PORT" default:"5432"`
		Options  string `envconfig:"POSTGRES_CONN_OPTS"`
	}
)

var (
	DBConfig DBConfigParameters
)

func LoadDBConfig() {
	envconfig.MustProcess("", &DBConfig)
}

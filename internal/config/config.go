package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Environment struct {
	Address         string `envconfig:"ADDRESS"`
	SecretPhrase    string `envconfig:"SECRET_PHRASE"`
	AdminEmail      string `envconfig:"ADMIN_EMAIL"`
	SlackWebhookUrl string `envconfig:"SLACK_WEBHOOK_URL"`
	DBHost          string `envconfig:"DB_HOST"`
	DBPort          string `envconfig:"DB_PORT" default:"3306"`
	DBUser          string `envconfig:"DB_USER"`
	DBPassword      string `envconfig:"DB_PASSWORD"`
	DBName          string `envconfig:"DB_NAME"`
}

var Env Environment

const (
	EnvLocal = "local"
	EnvProd  = "prod"
)

func init() {
	err := envconfig.Process("", &Env)
	if err != nil {
		panic(err)
	}
}

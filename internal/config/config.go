package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Environment struct {
	Address         string `envconfig:"ADDRESS"`
	SecretPhrase    string `envconfig:"SECRET_PHRASE"`
	AdminEmail      string `envconfig:"ADMIN_EMAIL"`
	SlackWebhookUrl string `envconfig:"SLACK_WEBHOOK_URL"`
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

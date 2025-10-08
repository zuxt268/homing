package config

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Environment struct {
	Address                string `envconfig:"ADDRESS"`
	SecretPhrase           string `envconfig:"SECRET_PHRASE"`
	AdminEmail             string `envconfig:"ADMIN_EMAIL"`
	NoticeWebAppChannelUrl string `envconfig:"NOTICE_WEB_APP_CHANNEL_URL"`
	PrjARootChannelUrl     string `envconfig:"PRJ_AROOT_CHANNEL_URL"`
	NoticeRpaChannelUrl    string `envconfig:"NOTICE_RPA_CHANNEL_URL"`
	DBHost                 string `envconfig:"DB_HOST"`
	DBPort                 string `envconfig:"DB_PORT" default:"3306"`
	DBUser                 string `envconfig:"DB_USER"`
	DBPassword             string `envconfig:"DB_PASSWORD"`
	DBName                 string `envconfig:"DB_NAME"`
	ClientID               string `envconfig:"CLIENT_ID"`
	ClientSecret           string `envconfig:"CLIENT_SECRET"`
}

var Env Environment

const (
	EnvLocal = "local"
	EnvProd  = "prod"
)

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}
	err := envconfig.Process("", &Env)
	if err != nil {
		panic(err)
	}
}

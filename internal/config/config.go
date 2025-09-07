package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Environment struct {
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

package conf

import (
	"fmt"
	"log"

	"github.com/kelseyhightower/envconfig"
)

type LogConfig struct {
	EnableStacktrace bool `envconfig:"AUTH_LOG_ENABLE_STACKTRACE" default:"false"`
}

type Config struct {
	Debug bool `envconfig:"AUTH_DEBUG" default:"false"`

	Port int `envconfig:"AUTH_PORT" default:"8080"`

	Log LogConfig
}

func New() *Config {
	cfg := new(Config)

	if err := envconfig.Process("", cfg); err != nil {
		err = fmt.Errorf("error while parse env config | %w", err)

		log.Fatal(err)
	}

	return cfg
}

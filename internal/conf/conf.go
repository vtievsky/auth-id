package conf

import (
	"fmt"
	"log"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type DBConfig struct {
	URL string `envconfig:"AUTH_DB_URL" required:"true"`
}

type SessionConfig struct {
	URL             string        `envconfig:"AUTH_SESSION_URL" required:"true"`
	SigningKey      string        `envconfig:"AUTH_SESSION_SIGNING_KEY" required:"true"`
	SessionTTL      time.Duration `envconfig:"AUTH_SESSION_TTL" default:"24h"`
	AccessTokenTTL  time.Duration `envconfig:"AUTH_ACCESS_TOKEN_TTL" default:"15m"`
	RefreshTokenTTL time.Duration `envconfig:"AUTH_REFRESH_TOKEN_TTL" default:"24h"`
}

type LogConfig struct {
	EnableStacktrace bool `envconfig:"AUTH_LOG_ENABLE_STACKTRACE" default:"false"`
}

// type TracingConfig struct {
// 	Enabled bool    `envconfig:"AUTH_TRACING_ENABLED" default:"true"`
// 	Address string  `envconfig:"AUTH_TRACING_URL" required:"true"`
// 	Ratio   float64 `envconfig:"AUTH_TRACING_RATIO" default:"1.0"`
// }

type MetricsConfig struct {
	URL     string    `envconfig:"AUTH_METRICS_URL" required:"true"`
	Buckets []float64 `envconfig:"AUTH_METRICS_BUCKETS"`
}

type Config struct {
	Port        int    `envconfig:"AUTH_PORT" default:"8080"`
	Debug       bool   `envconfig:"AUTH_DEBUG" default:"false"`
	ServiceName string `envconfig:"AUTH_SERVER_SERVICE_NAME" required:"true"`
	Environment string `envconfig:"AUTH_ENVIRONMENT" required:"true"`
	DB          DBConfig
	Log         LogConfig
	Session     SessionConfig
	// Tracing     TracingConfig
	Metrics MetricsConfig
}

func New() *Config {
	cfg := new(Config)

	if err := envconfig.Process("", cfg); err != nil {
		err = fmt.Errorf("error while parse env config | %w", err)

		log.Fatal(err)
	}

	return cfg
}

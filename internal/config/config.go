package config

import (
	"time"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Port            string        `env:"PORT" envDefault:":8080"`
	Env             string        `env:"APP_ENV" envDefault:"local"`
	IdentityURL     string        `env:"IDENTITY_URL" envDefault:"http://localhost:8080"`
	AvailabilityURL string        `env:"AVAILABILITY_URL" envDefault:"http://localhost:8080"`
	S3BaseURL       string        `env:"S3_BASE_URL" envDefault:"https://s3.eu-west-1.amazon.com/bytestreamfake"`
	HTTPTimeout     time.Duration `env:"HTTP_TIMEOUT" envDefault:"10s"`
}

func Load() (Config, error) {
	return env.ParseAs[Config]()
}

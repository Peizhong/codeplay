package config

import (
	"os"

	"github.com/caarlos0/env/v10"
	"github.com/tidwall/gjson"
)

type Config struct {
	Hostname     string `env:"CODEPLAY_HOSTNAME"`
	HttpPort     int    `env:"CODEPLAY_HTTP_PORT" envDefault:"3000"`
	FeatureGates string `env:"CODEPLAY_FEATURE_GATES" envDefault:"{}"`

	RedisAddr     string `env:"REDIS_ADDR" envDefault:"10.10.10.1:6379"`
	RedisPassword string `env:"REDIS_PASSWORD" envDefault:"app_user"`
}

var C Config

func init() {
	env.Parse(&C)
	if C.Hostname == "" {
		C.Hostname, _ = os.Hostname()
	}
}

func (c Config) GetFeature(name string) gjson.Result {
	return gjson.Parse(c.FeatureGates).Get(name)
}

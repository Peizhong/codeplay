package config

import (
	"os"

	"github.com/caarlos0/env/v10"
	"github.com/tidwall/gjson"
)

type Config struct {
	HttpPort     int    `env:"CODEPLAY_HTTP_PORT" envDefault:"3000"`
	FeatureGates string `env:"CODEPLAY_FEATURE_GATES" envDefault:"{}"`
}

var C Config
var HostName string

func Init() {
	HostName, _ = os.Hostname()
	env.Parse(&C)
}

func (c Config) GetFeature(name string) gjson.Result {
	return gjson.Parse(c.FeatureGates).Get(name)
}

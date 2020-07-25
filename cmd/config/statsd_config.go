package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type StatsDConfig struct {
	host      string
	port      string
	namespace string
}

func newStatsDConfig() StatsDConfig {
	return StatsDConfig{
		host:      viper.GetString("STATSD_HOST"),
		port:      viper.GetString("STATSD_PORT"),
		namespace: viper.GetString("STATSD_NAMESPACE"),
	}
}

func (sdc StatsDConfig) Namespace() string {
	return sdc.namespace
}

func (sdc StatsDConfig) Address() string {
	return fmt.Sprintf("%s:%s", sdc.host, sdc.port)
}

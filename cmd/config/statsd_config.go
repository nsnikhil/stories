package config

import (
	"fmt"
)

type StatsDConfig struct {
	host      string
	port      string
	namespace string
}

func newStatsDConfig() StatsDConfig {
	return StatsDConfig{
		host:      getString("STATSD_HOST"),
		port:      getString("STATSD_PORT"),
		namespace: getString("STATSD_NAMESPACE"),
	}
}

func (sdc StatsDConfig) Namespace() string {
	return sdc.namespace
}

func (sdc StatsDConfig) Address() string {
	return fmt.Sprintf("%s:%s", sdc.host, sdc.port)
}

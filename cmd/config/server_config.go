package config

import (
	"fmt"
)

type ServerConfig struct {
	host                           string
	port                           string
	idleConnectionTimeoutInMinutes int
}

func newServerConfig() ServerConfig {
	return ServerConfig{
		host:                           getString("APP_HOST"),
		port:                           getString("APP_PORT"),
		idleConnectionTimeoutInMinutes: getInt("IDLE_CONNECTION_TIMEOUT_IN_MINUTES"),
	}
}

func (sc ServerConfig) Address() string {
	return fmt.Sprintf(":%s", sc.port)
}

func (sc ServerConfig) IdleConnectionTimeoutInMinutes() int {
	return sc.idleConnectionTimeoutInMinutes
}

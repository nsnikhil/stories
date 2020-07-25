package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type ServerConfig struct {
	host                           string
	port                           string
	idleConnectionTimeoutInMinutes int
}

func newServerConfig() ServerConfig {
	return ServerConfig{
		host:                           viper.GetString("APP_HOST"),
		port:                           viper.GetString("APP_PORT"),
		idleConnectionTimeoutInMinutes: viper.GetInt("IDLE_CONNECTION_TIMEOUT_IN_MINUTES"),
	}
}

func (sc ServerConfig) Address() string {
	return fmt.Sprintf(":%s", sc.port)
}

func (sc ServerConfig) IdleConnectionTimeoutInMinutes() int {
	return sc.idleConnectionTimeoutInMinutes
}

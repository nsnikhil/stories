package config

import (
	"fmt"
)

type GRPCServerConfig struct {
	host                           string
	port                           string
	idleConnectionTimeoutInMinutes int
}

func newGRPCServerConfig() GRPCServerConfig {
	return GRPCServerConfig{
		host:                           getString("GRPC_SERVER_HOST"),
		port:                           getString("GRPC_SERVER_PORT"),
		idleConnectionTimeoutInMinutes: getInt("GRPC_SERVER_IDLE_CONNECTION_TIMEOUT_IN_MINUTES"),
	}
}

func (sc GRPCServerConfig) Address() string {
	return fmt.Sprintf(":%s", sc.port)
}

func (sc GRPCServerConfig) IdleConnectionTimeoutInMinutes() int {
	return sc.idleConnectionTimeoutInMinutes
}

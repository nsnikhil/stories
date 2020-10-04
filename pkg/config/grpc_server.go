package config

import (
	"fmt"
)

type GRPCServerConfig struct {
	host string
	port int
	// TODO: REMOVE THIS FROM HERE
	prometheusHTTPServerPort       int
	idleConnectionTimeoutInMinutes int
}

func newGRPCServerConfig() GRPCServerConfig {
	return GRPCServerConfig{
		host:                           getString("GRPC_SERVER_HOST"),
		port:                           getInt("GRPC_SERVER_PORT"),
		prometheusHTTPServerPort:       getInt("GRPC_PROMETHEUS_HTTP_PORT"),
		idleConnectionTimeoutInMinutes: getInt("GRPC_SERVER_IDLE_CONNECTION_TIMEOUT_IN_MINUTES"),
	}
}

func (sc GRPCServerConfig) Address() string {
	return fmt.Sprintf(":%d", sc.port)
}

func (sc GRPCServerConfig) PrometheusHTTPAddress() string {
	return fmt.Sprintf(":%d", sc.prometheusHTTPServerPort)
}

func (sc GRPCServerConfig) IdleConnectionTimeoutInMinutes() int {
	return sc.idleConnectionTimeoutInMinutes
}

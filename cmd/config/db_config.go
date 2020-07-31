package config

import "fmt"

type DatabaseConfig struct {
	host                  string
	port                  int
	username              string
	password              string
	name                  string
	maxIdleConnections    int
	maxOpenConnections    int
	connectionMaxLifetime int
}

func newDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		host:                  getString("DB_HOST"),
		port:                  getInt("DB_PORT"),
		name:                  getString("DB_NAME"),
		username:              getString("DB_USER"),
		password:              getString("DB_PASSWORD"),
		maxIdleConnections:    getInt("DB_MAX_IDLE_CONNECTIONS"),
		maxOpenConnections:    getInt("DB_MAX_OPEN_CONNECTIONS"),
		connectionMaxLifetime: getInt("DB_CONNECTION_MAX_LIFETIME_IN_MIN"),
	}
}

func (dc DatabaseConfig) Source() string {
	return fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", dc.username, dc.password, dc.host, dc.port, dc.name)
}

func (dc DatabaseConfig) GetIdleConnections() int {
	return dc.maxIdleConnections
}

func (dc DatabaseConfig) GetMaxOpenConnections() int {
	return dc.maxOpenConnections
}

func (dc DatabaseConfig) GetConnectionMaxLifetime() int {
	return dc.connectionMaxLifetime
}

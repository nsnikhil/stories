package config

import "fmt"

type DatabaseConfig struct {
	host     string
	port     int
	username string
	password string
	name     string
}

func newDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		host:     getString("DB_HOST"),
		port:     getInt("DB_PORT"),
		name:     getString("DB_NAME"),
		username: getString("DB_USER"),
		password: getString("DB_PASSWORD"),
	}
}

func (dc *DatabaseConfig) Source() string {
	return fmt.Sprintf("dbname=%s user=%s password=%s host=%s port=%d sslmode=disable", dc.name, dc.username, dc.password, dc.host, dc.port)
}

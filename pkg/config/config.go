package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	env              string
	migrationPath    string
	grpcServerConfig GRPCServerConfig
	httpServerConfig HTTPServerConfig
	newRelicConfig   NewRelicConfig
	databaseConfig   DatabaseConfig
	storyConfig      StoryConfig
	logConfig        LogConfig
	logFileConfig    LogFileConfig
}

func (c Config) GRPCServerConfig() GRPCServerConfig {
	return c.grpcServerConfig
}

func (c Config) HTTPServerConfig() HTTPServerConfig {
	return c.httpServerConfig
}

func (c Config) NewRelicConfig() NewRelicConfig {
	return c.newRelicConfig
}

func (c Config) DatabaseConfig() DatabaseConfig {
	return c.databaseConfig
}

func (c Config) LogConfig() LogConfig {
	return c.logConfig
}

func (c Config) LogFileConfig() LogFileConfig {
	return c.logFileConfig
}

func (c Config) Env() string {
	return c.env
}

func (c Config) MigrationPath() string {
	return c.migrationPath
}

func (c Config) StoryConfig() StoryConfig {
	return c.storyConfig
}

func NewConfig(configFile string) Config {
	viper.AutomaticEnv()
	viper.SetConfigFile(configFile)

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}

	return Config{
		env:              getString("ENV"),
		migrationPath:    getString("MIGRATION_PATH"),
		grpcServerConfig: newGRPCServerConfig(),
		httpServerConfig: newHTTPServerConfig(),
		newRelicConfig:   newNewRelicConfig(),
		databaseConfig:   newDatabaseConfig(),
		storyConfig:      newStoryConfig(),
		logConfig:        newLogConfig(),
		logFileConfig:    newLogFileConfig(),
	}
}

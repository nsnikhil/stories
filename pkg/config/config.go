package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	env            string
	migrationPath  string
	serverConfig   ServerConfig
	newRelicConfig NewRelicConfig
	databaseConfig DatabaseConfig
	storyConfig    StoryConfig
	logConfig      LogConfig
	logFileConfig  LogFileConfig
}

func (c Config) ServerConfig() ServerConfig {
	return c.serverConfig
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

func NewConfig() Config {
	viper.AutomaticEnv()
	viper.SetConfigName("local")
	viper.SetConfigType("env")
	viper.AddConfigPath("./")
	viper.AddConfigPath("./../")
	viper.AddConfigPath("./../../")
	viper.AddConfigPath("./../../../")
	viper.AddConfigPath("./../../../../")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}

	return Config{
		env:            getString("ENV"),
		migrationPath:  getString("MIGRATION_PATH"),
		serverConfig:   newServerConfig(),
		newRelicConfig: newNewRelicConfig(),
		databaseConfig: newDatabaseConfig(),
		storyConfig:    newStoryConfig(),
		logConfig:      newLogConfig(),
		logFileConfig:  newLogFileConfig(),
	}

}

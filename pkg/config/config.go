package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	env            string
	serverConfig   ServerConfig
	newRelicConfig NewRelicConfig
	databaseConfig DatabaseConfig
	storyConfig    StoryConfig
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

func (c Config) Env() string {
	return c.env
}

func (c Config) BlogConfig() StoryConfig {
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

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}

	return Config{
		env:            getString("ENV"),
		serverConfig:   newServerConfig(),
		newRelicConfig: newNewRelicConfig(),
		databaseConfig: newDatabaseConfig(),
		storyConfig:    newStoryConfig(),
	}

}

package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	env string
	sc  ServerConfig
	nr  NewRelicConfig
	sdc StatsDConfig
	dbc DatabaseConfig
	bc  BlogConfig
}

func (c Config) GetServerConfig() ServerConfig {
	return c.sc
}

func (c Config) GetNewRelicConfig() NewRelicConfig {
	return c.nr
}

func (c Config) GetStatsDConfig() StatsDConfig {
	return c.sdc
}

func (c Config) GetDatabaseConfig() DatabaseConfig {
	return c.dbc
}

func (c Config) GetEnv() string {
	return c.env
}

func (c Config) GetBlogConfig() BlogConfig {
	return c.bc
}

func LoadConfigs() Config {
	viper.AutomaticEnv()
	viper.SetConfigName("env")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./")
	viper.AddConfigPath("./../../../")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}

	return Config{
		env: getString("ENV"),
		sc:  newServerConfig(),
		nr:  newNewRelicConfig(),
		sdc: newStatsDConfig(),
		dbc: newDatabaseConfig(),
		bc:  newBlogConfig(),
	}

}

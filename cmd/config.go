package main

import (
	"fmt"
	"github.com/spf13/viper"
)

type serverConfig struct {
	host string
	port string
}

type newRelicConfig struct {
	appName    string
	licenseKey string
}

type statsDConfig struct {
	host      string
	port      string
	namespace string
}

func (sdc statsDConfig) address() string {
	return fmt.Sprintf("%s:%s", sdc.host, sdc.port)
}

func (sc serverConfig) address() string {
	return fmt.Sprintf(":%s", sc.port)
}

type config struct {
	env string
	sc  serverConfig
	nr  newRelicConfig
	sdc statsDConfig
}

var cfg config

func initConfigs() {
	viper.AutomaticEnv()
	viper.SetConfigName("env")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	cfg = config{
		env: viper.GetString("ENV"),
		sc: serverConfig{
			host: viper.GetString("APP_HOST"),
			port: viper.GetString("APP_PORT"),
		},
		nr: newRelicConfig{
			appName:    viper.GetString("NEW_RELIC_APP_NAME"),
			licenseKey: viper.GetString("NEW_RELIC_LICENSE_KEY"),
		},
		sdc: statsDConfig{
			host:      viper.GetString("STATSD_HOST"),
			port:      viper.GetString("STATSD_PORT"),
			namespace: viper.GetString("STATSD_NAMESPACE"),
		},
	}

}

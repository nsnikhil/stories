package config

import "github.com/spf13/viper"

type NewRelicConfig struct {
	appName    string
	licenseKey string
}

func (nrc NewRelicConfig) AppName() string {
	return nrc.appName
}

func (nrc NewRelicConfig) LicenseKey() string {
	return nrc.licenseKey
}

func newNewRelicConfig() NewRelicConfig {
	return NewRelicConfig{
		appName:    viper.GetString("NEW_RELIC_APP_NAME"),
		licenseKey: viper.GetString("NEW_RELIC_LICENSE_KEY"),
	}
}

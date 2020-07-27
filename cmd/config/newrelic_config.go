package config

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
		appName:    getString("NEW_RELIC_APP_NAME"),
		licenseKey: getString("NEW_RELIC_LICENSE_KEY"),
	}
}

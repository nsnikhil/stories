package reporters

import (
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/nsnikhil/stories/pkg/config"
)

func NewNewRelicApp(nrc config.NewRelicConfig) (*newrelic.Application, error) {
	return newrelic.NewApplication(
		newrelic.ConfigAppName(nrc.AppName()),
		newrelic.ConfigLicense(nrc.LicenseKey()),
	)
}

package reporters

import (
	"fmt"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/nsnikhil/stories/pkg/config"
)

func NewNewRelicApp(nrc config.NewRelicConfig) *newrelic.Application {
	var err error
	var nrApp *newrelic.Application

	nrApp, err = newrelic.NewApplication(
		newrelic.ConfigAppName(nrc.AppName()),
		newrelic.ConfigLicense(nrc.LicenseKey()),
	)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	return nrApp
}

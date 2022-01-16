package metrics

import (
	"github.com/kelseyhightower/envconfig"
	newrelic "github.com/newrelic/go-agent"
	"log"
	"time"
)

type (
	parameters struct {
		AppName string `envconfig:"APP_NAME"`
		License string `envconfig:"LICENSE"`
	}

	newrelicmetrics struct {
		cli newrelic.Application
	}
)

func New() IMetrics {
	var cfg parameters

	if err := envconfig.Process("NEWRELIC", &cfg); err != nil {
		log.Println("Metrics config not loaded", err)
		return nil
	}

	app, err := newrelic.NewApplication(newrelic.Config{
		Enabled:           true,
		AppName:           cfg.AppName,
		License:           cfg.License,
		DistributedTracer: struct{ Enabled bool }{Enabled: true},
	})
	if err != nil {
		log.Println("Metrics are not available:", err)
		return nil
	}

	return &newrelicmetrics{
		cli: app,
	}
}

func (m *newrelicmetrics) Start() error {
	if m.cli == nil {
		return nil
	}

	return m.cli.WaitForConnection(20 * time.Second)
}

func (m *newrelicmetrics) Stop() {
	if m.cli == nil {
		return
	}

	m.cli.Shutdown(5 * time.Second)
}
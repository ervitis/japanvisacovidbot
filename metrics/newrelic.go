package metrics

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/newrelic/go-agent/v3/newrelic"
	"log"
	"time"
)

type (
	parameters struct {
		AppName string `envconfig:"APP_NAME"`
		License string `envconfig:"LICENSE"`
	}

	newrelicmetrics struct {
		cli *newrelic.Application
	}
)

func New() IMetrics {
	var cfg parameters

	if err := envconfig.Process("NEWRELIC", &cfg); err != nil {
		log.Println("Metrics config not loaded", err)
		return nil
	}

	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName(cfg.AppName),
		newrelic.ConfigDistributedTracerEnabled(true),
		newrelic.ConfigLicense(cfg.License),
		newrelic.ConfigEnabled(true),
	)
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

func (m *newrelicmetrics) ExecuteWithSegment(name string, fn func() error) error {
	txn := m.cli.StartTransaction(name)
	defer txn.End()
	return fn()
}

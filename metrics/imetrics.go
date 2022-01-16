package metrics

type (
	IMetrics interface {
		Start() error
		Stop()
	}
)

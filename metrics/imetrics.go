package metrics

type (
	IMetrics interface {
		Start() error
		Stop()
		ExecuteWithSegment(name string, fn func() error) error
	}
)

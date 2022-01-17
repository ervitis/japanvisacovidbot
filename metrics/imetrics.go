package metrics

import "context"

type (
	IMetrics interface {
		Start() error
		Stop()
		ExecuteWithSegment(ctx context.Context, name string, fn func() error) error
	}
)

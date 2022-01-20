package metrics

import "context"

type (
	IMetrics interface {
		Start() error
		Stop()
		ExecuteWithSegment(context.Context, string, func(context context.Context) error) error
	}
)

package metrics

import (
	"context"
	"fmt"
	"log"
)

type (
	nometric struct{}
)

func (n *nometric) Start() error {
	return nil
}

func (n *nometric) Stop() {}

func (n nometric) ExecuteWithSegment(ctx context.Context, s string, f func(context context.Context) error) error {
	log.Println(fmt.Sprintf("executing segment %s", s))
	return f(ctx)
}

func NoMetricNew() IMetrics {
	return &nometric{}
}

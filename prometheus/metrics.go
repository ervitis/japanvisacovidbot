package prometheus

import (
	"context"
	"errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"time"
)

type (
	metrics struct {
		srv *http.Server
	}
)

const (
	connTimeout = 15 * time.Second
)

func New() *metrics {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	srv := &http.Server{
		Addr:         ":8585",
		Handler:      mux,
		ReadTimeout:  connTimeout,
		IdleTimeout:  connTimeout,
		WriteTimeout: connTimeout,
	}

	return &metrics{
		srv: srv,
	}
}

func (m *metrics) StartMetricsServer() {
	log.Println("starting metrics server")

	if err := m.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Panic("error listening metrics server", err)
	}
}

func (m *metrics) Close() {
	log.Println("shutting down metrics server")
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := m.srv.Shutdown(ctx); err != nil {
		log.Fatalln("error shutting metrics server", err)
	}
}

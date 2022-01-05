package japanvisacovidbot

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"time"
)

type (
	httpServer struct {
		srv *http.Server
	}
)

const (
	connTimeout = 30 * time.Second
)

func NewServer() *httpServer {
	port := os.Getenv("PORT")
	if port == "" {
		port = ":8085"
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/health", handleHealthChecker())

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadTimeout:       connTimeout,
		ReadHeaderTimeout: connTimeout,
		WriteTimeout:      connTimeout,
		IdleTimeout:       connTimeout,
	}

	return &httpServer{srv: srv}
}

func (h *httpServer) StartServer() {
	log.Println("starting http server")

	if err := h.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Panic("error listening server", err)
	}
}

func (h *httpServer) Close() {
	log.Println("shutting down http server")
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := h.srv.Shutdown(ctx); err != nil {
		log.Fatalln("error shutting down server", err)
	}
}

func handleHealthChecker() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}

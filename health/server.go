package health

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type server struct {
	mux *http.ServeMux
}

func (s server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

type ProbeServer struct {
	logger log.Logger
	srv    *http.Server
}

// Start binds the probe to a listening port
func (s *ProbeServer) Start() {
	go func() {
		if err := s.srv.ListenAndServe(); err != nil {
			level.Error(s.logger).Log("message", "failed to start kit sidecar server", "error", err)
		}
	}()
}

// Stop shutsdown the probe
func (s *ProbeServer) Stop() error {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	return s.srv.Shutdown(ctx)
}

func NewServer(port int, logger log.Logger, checkers map[string]Checker) *ProbeServer {
	mux := http.NewServeMux()

	mux.HandleFunc("/live", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("/ready", Handler(logger, checkers))
	mux.Handle("/metrics", promhttp.Handler())

	addr := fmt.Sprintf(":%d", port)
	s := &ProbeServer{
		srv: &http.Server{Addr: addr, Handler: server{mux: mux}},
	}

	return s
}

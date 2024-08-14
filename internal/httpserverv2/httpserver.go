package httpserverv2

import (
	"log/slog"
	"net/http"
)

type Config struct {
	Addr string
}

// HttpServer implements concurrent simple HTTP 1.1 Server using stdlib/http package.
type HttpServer struct {
	config Config

	logger *slog.Logger

	server *http.Server
}

func New(config Config, logger *slog.Logger) *HttpServer {
	s := &HttpServer{
		config: config,
		logger: logger,
	}

	mux := http.NewServeMux()

	api := NewAPI(logger)
	api.Register(mux)

	httpServer := &http.Server{
		Addr:    s.config.Addr,
		Handler: mux,
	}
	s.server = httpServer
	return s
}

func (s *HttpServer) Start() error {

	if err := s.server.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

package apiserver

import (
	"net/http"
)

type APIServer struct {
	httpServer *http.Server
	config     *Config
}

func New(config *Config) *APIServer {
	return &APIServer{
		config: config,
	}
}

func (s *APIServer) Run(r http.Handler) error {
	s.httpServer = &http.Server{
		Addr:    s.config.Addr,
		Handler: r,
	}
	return s.httpServer.ListenAndServe()
}

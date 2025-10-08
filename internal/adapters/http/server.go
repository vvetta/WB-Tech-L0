package http

import (
	"context"
	"net/http"
	"time"

	"WB-Tech-L0/internal/usecase"
)

type Deps struct {
	OrderSvc usecase.OrderReader
	Logger   usecase.Logger
}

type Server struct {
	http *http.Server
}

func NewServer(addr string, d Deps) *Server {
	mux := NewRouter(d)
	s := &Server{
		http: &http.Server{
			Addr:              addr,
			Handler:           mux,
			ReadHeaderTimeout: 5 * time.Second,
			WriteTimeout:      15 * time.Second,
			IdleTimeout:       60 * time.Second,
		},
	}
	return s
}

func (s *Server) Start() error                       { return s.http.ListenAndServe() }
func (s *Server) Shutdown(ctx context.Context) error { return s.http.Shutdown(ctx) }

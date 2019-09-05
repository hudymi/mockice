package service

import (
	"fmt"
	"net/http"
)

type mux interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
}

type service struct {
	host string
	mux  mux
}

// New creates a new service that will listen at provided host
func New(host string) *service {
	return &service{
		host: host,
		mux:  http.NewServeMux(),
	}
}

// Start starts serving service
func (s *service) Start() error {
	service := &http.Server{Addr: s.host, Handler: s.mux}

	return service.ListenAndServe()
}

// Handler is a interface of endpoint that can be registered in service
type Handler interface {
	Handle(writer http.ResponseWriter, request *http.Request)
	Name() string
}

// Register adds endpoint to service
func (s *service) Register(handler Handler) {
	s.mux.HandleFunc(fmt.Sprintf("/%s", handler.Name()), handler.Handle)
}

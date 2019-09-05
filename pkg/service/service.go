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

func New(host string) *service {
	return &service{
		host: host,
		mux:  http.NewServeMux(),
	}
}

func (s *service) Start() error {
	service := &http.Server{Addr: s.host, Handler: s.mux}

	return service.ListenAndServe()
}

type Handler interface {
	Handle(writer http.ResponseWriter, request *http.Request)
	Name() string
}

func (s *service) Register(handler Handler) {
	s.mux.HandleFunc(fmt.Sprintf("/%s", handler.Name()), handler.Handle)
}

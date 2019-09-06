package service

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type mux interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
}

// Service is a structure responsible for serving HTTP server
type Service struct {
	host string
	mux  mux
}

// New creates a new service that will listen at provided host
func New(host string) *Service {
	return &Service{
		host: host,
		mux:  http.NewServeMux(),
	}
}

// Start starts serving service
func (s *Service) Start(ctx context.Context) error {
	service := &http.Server{Addr: s.host, Handler: s.mux}

	go func() {
		if err := service.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Errorf("Error while starting HTTP service: %v", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	return service.Shutdown(context.Background())
}

// Handler is a interface of endpoint that can be registered in service
type Handler interface {
	Handle(writer http.ResponseWriter, request *http.Request)
	Name() string
}

// Register adds endpoint to service
func (s *Service) Register(handler Handler) error {
	if handler.Name() == "" {
		return errors.New("Endpoint name cannot be empty")
	}

	s.mux.HandleFunc(fmt.Sprintf("/%s", handler.Name()), handler.Handle)
	return nil
}

package endpoint

import (
	"net/http"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Config stores endpoint configuration
type Config struct {
	Name                       string
	Methods                    []string
	DefaultResponseCode        *int
	DefaultResponseContent     string
	DefaultResponseContentType *string
}

// DefaultConfig generates default configuration with /hello endpoint
func DefaultConfig() []Config {
	return []Config{
		{
			Name:                   "hello",
			DefaultResponseContent: "Hello World! Mockice here!",
		},
	}
}

// Endpoint is a structure responsible for handling HTTP requests
type Endpoint struct {
	config Config
	log    *logrus.Entry
}

// New creates a new endpoint from config
func New(config Config) *Endpoint {
	return &Endpoint{
		config: config,
		log:    logrus.WithField("Endpoint", config.Name),
	}
}

// Handle is responsible for handling incoming requests
func (e *Endpoint) Handle(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()
	e.log.Infof("Handle %s request from %s", request.Method, request.RemoteAddr)

	if len(e.config.Methods) > 0 && !e.contains(e.config.Methods, request.Method) {
		http.Error(writer, "Invalid request method", http.StatusMethodNotAllowed)
		e.log.Errorf("Invalid request method %s", request.Method)
		return
	}

	if e.config.DefaultResponseContentType != nil {
		writer.Header().Set("Content-Type", *e.config.DefaultResponseContentType)
	} else {
		writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	}

	if e.config.DefaultResponseCode != nil {
		writer.WriteHeader(*e.config.DefaultResponseCode)
	}

	_, err := writer.Write([]byte(e.config.DefaultResponseContent))
	if err != nil {
		err = errors.Wrapf(err, "while writing response")
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		e.log.Error(err)
		return
	}
}

// Name returns name of the endpoint
func (e *Endpoint) Name() string {
	return e.config.Name
}

func (e *Endpoint) contains(values []string, value string) bool {
	for _, item := range values {
		if item == value {
			return true
		}
	}
	return false
}

package endpoint

import (
	"net/http"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type EndpointConfig struct {
	Name                       string
	Methods                    []string
	DefaultResponseCode        *int
	DefaultResponseContent     string
	DefaultResponseContentType *string
}

func DefaultConfig() []EndpointConfig {
	return []EndpointConfig{
		{
			Name:                   "hello",
			DefaultResponseContent: "Hello World! Mockice here!",
		},
	}
}

type endpoint struct {
	config EndpointConfig
}

func New(config EndpointConfig) *endpoint {
	return &endpoint{
		config: config,
	}
}

func (e *endpoint) Handle(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()

	if len(e.config.Methods) > 0 && !e.contains(e.config.Methods, request.Method) {
		http.Error(writer, "Invalid request method", http.StatusMethodNotAllowed)
		logrus.Errorf("Invalid request method %s", request.Method)
		return
	}

	if e.config.DefaultResponseContentType != nil {
		writer.Header().Set("Content-Type", *e.config.DefaultResponseContentType)
	} else {
		writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	}

	_, err := writer.Write([]byte(e.config.DefaultResponseContent))
	if err != nil {
		err = errors.Wrapf(err, "while writing response from /%s endpoint", e.config.Name)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		logrus.Error(err)
		return
	}

	if e.config.DefaultResponseCode != nil {
		writer.WriteHeader(*e.config.DefaultResponseCode)
	}
}

func (e *endpoint) Name() string {
	return e.config.Name
}

func (e *endpoint) contains(values []string, value string) bool {
	for _, item := range values {
		if item == value {
			return true
		}
	}
	return false
}

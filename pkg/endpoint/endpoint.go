package endpoint

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Config stores endpoint configuration
type Config struct {
	Name                       string   `yaml:"name"`
	Methods                    []string `yaml:"methods"`
	DefaultResponseCode        *int     `yaml:"defaultResponseCode"`
	DefaultResponseContent     string   `yaml:"defaultResponseContent"`
	DefaultResponseContentType *string  `yaml:"defaultResponseContentType"`
	DefaultResponseFile        *string  `yaml:"defaultResponseFile"`
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
		log:    logrus.WithField("endpoint", config.Name),
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

	writer.Header().Set("Content-Type", e.defaultContentType())

	var reader io.ReadCloser
	if e.config.DefaultResponseFile != nil {
		file, err := e.fileReader(writer)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			e.log.Error(err)
			return
		}

		reader = file
	} else {
		reader = ioutil.NopCloser(strings.NewReader(e.config.DefaultResponseContent))
	}
	defer reader.Close()

	if e.config.DefaultResponseCode != nil {
		writer.WriteHeader(*e.config.DefaultResponseCode)
	}

	_, err := io.Copy(writer, reader)
	if err != nil {
		err = errors.Wrapf(err, "while writing response")
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		e.log.Error(err)
		return
	}
}

func (e *Endpoint) defaultContentType() string {
	if e.config.DefaultResponseContentType != nil {
		return *e.config.DefaultResponseContentType
	}
	return "text/plain; charset=utf-8"
}

func (e *Endpoint) fileReader(writer io.Writer) (io.ReadCloser, error) {
	file, err := os.Open(*e.config.DefaultResponseFile)
	if err != nil {
		return nil, errors.Wrapf(err, "while opening file %s", *e.config.DefaultResponseFile)
	}

	return file, nil
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

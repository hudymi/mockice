package main

import (
	"flag"
	"io/ioutil"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"github.com/michal-hudy/mockice/pkg/endpoint"
	"github.com/michal-hudy/mockice/pkg/log"
	"github.com/michal-hudy/mockice/pkg/service"
)

type options struct {
	verbose bool
	config  string
}

type config struct {
	Address                string                    `yaml:"address"`
	EndpointsConfiguration []endpoint.EndpointConfig `yaml:"endpoints"`
}

func gatherOptions() options {
	o := options{}

	flag.StringVar(&o.config, "config", "", "The path to the configuration file")
	flag.BoolVar(&o.verbose, "verbose", false, "Enable verbose output")
	flag.Parse()

	return o
}

func main() {
	options := gatherOptions()
	log.Setup(options.verbose)

	cfg, err := loadConfig(options)
	if err != nil {
		logrus.Fatal(err)
	}

	svc := service.New(cfg.Address)
	for _, config := range cfg.EndpointsConfiguration {
		endpoint := endpoint.New(config)
		logrus.Infof("Registering /%s endpoint", endpoint.Name())
		svc.Register(endpoint)
	}

	logrus.Infof("Service listen at %s", cfg.Address)
	err = svc.Start()
	if err != nil {
		logrus.Error(err)
	}
}

func loadConfig(ops options) (config, error) {
	cfg := config{
		Address: ":8080",
	}

	if ops.config == "" {
		cfg.EndpointsConfiguration = endpoint.DefaultConfig()
		return cfg, nil
	}

	content, err := ioutil.ReadFile(ops.config)
	if err != nil {
		return config{}, errors.Wrapf(err, "while reading configuration from %s", ops.config)
	}

	err = yaml.Unmarshal(content, &cfg)
	if err != nil {
		return config{}, errors.Wrapf(err, "while loading configuration from %s", ops.config)
	}

	return cfg, nil
}

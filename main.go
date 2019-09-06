package main

import (
	"flag"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"github.com/michal-hudy/mockice/pkg/endpoint"
	"github.com/michal-hudy/mockice/pkg/log"
	"github.com/michal-hudy/mockice/pkg/service"
	"github.com/michal-hudy/mockice/pkg/signal"
)

type options struct {
	verbose bool
	config  string
}

type config struct {
	Address                string            `yaml:"address"`
	EndpointsConfiguration []endpoint.Config `yaml:"endpoints"`
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

	ctx := signal.Context()

	cfg, err := loadConfig(options.config)
	if err != nil {
		logrus.Fatal(err)
	}

	svc := service.New(cfg.Address)
	for _, config := range cfg.EndpointsConfiguration {
		endpoint := endpoint.New(config)

		logrus.Infof("Registering /%s endpoint", endpoint.Name())
		err := svc.Register(endpoint)
		if err != nil {
			logrus.Fatal(err)
		}
	}

	logrus.Infof("Service listen at %s", cfg.Address)
	err = svc.Start(ctx)
	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
}

func loadConfig(path string) (config, error) {
	cfg := config{
		Address: ":8080",
	}

	if path == "" {
		cfg.EndpointsConfiguration = endpoint.DefaultConfig()
		return cfg, nil
	}

	content, err := ioutil.ReadFile(path)
	if err != nil {
		return config{}, errors.Wrapf(err, "while reading configuration from %s", path)
	}

	err = yaml.Unmarshal(content, &cfg)
	if err != nil {
		return config{}, errors.Wrapf(err, "while loading configuration from %s", path)
	}

	return cfg, nil
}

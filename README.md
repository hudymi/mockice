# Mockice

[![Build Status](https://github.com/michal-hudy/mockice/workflows/build/badge.svg)](https://github.com/michal-hudy/mockice/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/michal-hudy/mockice/actions)](https://goreportcard.com/report/github.com/michal-hudy/mockice/actions)

Mockice is a simple HTTP service that provides configurable endpoints. It can be used for testing or serving a static content.

## Installation

To install Mockice, run:

```bash
go get -u -v github.com/michal-hudy/mockice
```

## Usage

### Run latest version

To run latest version, install a Mockice with command from [Installation](#Installation) section and run:

```bash
mockice --verbose
```

The service listens on port `8080` and has one endpoint `http://localhost:8080/hello`

### Run a Docker image

To run Mockice latest Docker image, run:

```bash
docker run -p 8080:8080 hudymi/mockice:latest --verbose
```

The service listens on port `8080` and has one endpoint `http://localhost:8080/hello`

### Run from sources

To run Mockice from sources, run:

```bash
GO111MODULE=on go run main.go --verbose
``` 

The service listens on port `8080` and has one endpoint `http://localhost:8080/hello`

### Command line parameters

You can use the following parameters:

| Name | Description | Default Value |
| ---- | ----------- | ------------- |
| `--config` | Path to the configuration file. If not provided the default configuration will be used. | |
| `--verbose` | Enable verbose logging | `false` |

## Configuration

By default Mockice listens on every interfaces on port `8080` and has one endpoint `http://localhost:8080/hello`. If configuration is provided then default endpoint is disabled.

Configuration file must be in `yaml` format and for up-to-date list of available fields see [config](main.go) structure.

### File structure

```yaml
# The service address
address: :8080
# The list of endpoints
endpoints:
- name: hello # Name of the endpoint
  # The list of valid methods
  methods:
  - GET
  - POST
  # Default HTTP response code, if not provided 200
  defaultResponseCode: 200 
  # Default response content
  defaultResponseContent: "# Sample service"  
  # Default response content-type, if not provided "text/plain; charset=utf-8"
  defaultResponseContentType: text/plain; charset=utf-8

```
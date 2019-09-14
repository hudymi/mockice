# Mockice

[![Build Status](https://github.com/michal-hudy/mockice/workflows/build/badge.svg)](https://github.com/michal-hudy/mockice/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/michal-hudy/mockice/actions)](https://goreportcard.com/report/github.com/michal-hudy/mockice/actions)

Mockice is a simple HTTP service that provides configurable endpoints. Use it for testing or serving static content.

## Installation

To install Mockice, run:

```bash
go get -u -v github.com/michal-hudy/mockice
```

## Usage

To use Mockice, you can run its latest version, run its latest Docker image, or run it from sources. In each case, by default, the service listens on port `8080` and has one endpoint - `http://localhost:8080/hello`.

### Run the latest version

To run Mockice latest version, install Mockice with the command from the [Installation](#Installation) section and run:

```bash
mockice --verbose
```

### Run a Docker image

To run Mockice latest Docker image, use the following command:

```bash
docker run -p 8080:8080 hudymi/mockice:latest --verbose
```

### Run from sources

To run Mockice from sources, use the following command:

```bash
GO111MODULE=on go run main.go --verbose
```

### Command line parameters

The table contains the command line parameters available for the service:

| Name | Description | Default Value |
| ---- | ----------- | ------------- |
| `--config` | A path to the configuration file. If not provided, the default configuration is used. | |
| `--verbose` | Enables verbose logging. | `false` |

## Configuration

By default, in every interface Mockice listens on port `8080` and has one endpoint - `http://localhost:8080/hello`. If you provide any configuration, the default endpoint is disabled.

The configuration file must be in the `yaml` format. See the [config](main.go) structure for the up-to-date list of the available fields.

### File structure

```yaml
# The service address
address: :8080
# The list of endpoints
endpoints:
- name: hello # Name of the endpoint
  # The list of valid methods, if not set validation is skipped
  methods:
  - GET
  - POST
  # Default HTTP response code, if not provided 200
  defaultResponseCode: 200
  # Default response content
  defaultResponseContent: "Sample service"  
  # Default response content-type, if not provided "text/plain; charset=utf-8"
  defaultResponseContentType: text/plain; charset=utf-8
  # Path to the file that is returned by default, if provided then defaultResponseContent is ignored
  defaultResponseFile: "mockice/index.html"

```

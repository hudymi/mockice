# Build the service binary
FROM golang:1.12.9 as builder

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY main.go main.go
COPY pkg/ pkg/

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o service main.go

# Certificates
FROM alpine:latest as certs
RUN apk --update add ca-certificates

FROM scratch
LABEL source = git@github.com:michal-hudy/mockice.git
WORKDIR /
COPY --from=builder /workspace/service .
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
ENTRYPOINT ["/service"]
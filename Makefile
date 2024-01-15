TARGETOS ?= linux
TARGETARCH ?= arm

ifeq ($(TARGETARCH),arm)
	TARGETARM ?= 7
else
	TARGETARM ?=
endif


.PHONY: all - Default target
all: build

.PHONY: build - Build the collector
build: openvpnreceiver/metadata.go pireceiver/metadata.go
	builder --config builder-config.yaml --name otelcol-dev

.PHONY: release - Build the collector for release targeting TARGETOS and TARGETARCH
release: openvpnreceiver/metadata.go pireceiver/metadata.go
	GOOS=$(TARGETOS) GOARCH=$(TARGETARCH) GOARM=$(TARGETARM) builder --config builder-config.yaml --name otelcol-$(TARGETOS)-$(TARGETARCH)

.PHONY: docker - Build the collector using docker
docker:
	DOCKERBUILD_KIT=1 docker build -o ./build .

.PHONY: test - Run tests
test: test-openvpnreceiver test-pireceiver

.PHONY: test-openvpnreceiver - Run tests for openvpnreceiver
test-openvpnreceiver:
	cd openvpnreceiver && go test -v ./...

.PHONY: test-pireceiver - Run tests for pireceiver
test-pireceiver:
	cd pireceiver && go test -v ./...

.PHONY: clean - Clean build artifacts
clean:
	rm -rf ./build

.PHONY: setup - Install dependencies
setup:
	@if ! command -v go > /dev/null; then \
		echo "go version 1.19 or greater is required"; \
		exit 1; \
	fi
	@VERSION=$$(go version | awk -F. '{ gsub(/go/, "", $$1); printf("%d.%d", $$1, $$2) }'); \
	MAJOR=$$(echo "$$VERSION" | cut -d. -f1); \
	MINOR=$$(echo "$$VERSION" | cut -d. -f2); \
	if [ $$MAJOR -lt 2 ] && [ $$MINOR -lt 19 ]; then \
		echo "go version $$VERSION is installed, but version 1.19 or greater is required"; \
		exit 1; \
	fi
	go install go.opentelemetry.io/collector/cmd/builder@latest
	go install github.com/open-telemetry/opentelemetry-collector-contrib/cmd/mdatagen@latest

openvpnreceiver/metadata.go: openvpnreceiver/metadata.yaml
	cd openvpnreceiver && mdatagen metadata.yaml

pireceiver/metadata.go: pireceiver/metadata.yaml
	cd pireceiver && mdatagen metadata.yaml
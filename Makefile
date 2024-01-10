.PHONY: all - Default target
all: build

.PHONY: build - Build the collector
build: openvpnreceiver/metadata.go pireceiver/metadata.go
	builder --config builder-config.yaml

.PHONY: docker - Build docker image
docker:
	docker build -t otelcol:latest .

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

.PHONY: metadata - Generate metadata for openvpnreceiver
openvpnreceiver/metadata.go: openvpnreceiver/metadata.yaml
	cd openvpnreceiver && mdatagen metadata.yaml

.PHONY: metadata - Generate metadata for pireceiver
pireceiver/metadata.go: pireceiver/metadata.yaml
	cd pireceiver && mdatagen metadata.yaml
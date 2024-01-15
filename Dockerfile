FROM golang:1.21-bullseye AS build-env

ARG TARGETOS=linux
ARG TARGETARCH=arm
ARG TARGETARM=7

RUN go install go.opentelemetry.io/collector/cmd/builder@v0.91.0
RUN go install github.com/open-telemetry/opentelemetry-collector-contrib/cmd/mdatagen@v0.91.0

WORKDIR /otelcol
COPY ./openvpnreceiver ./openvpnreceiver
COPY ./pireceiver ./pireceiver
COPY builder-config.yaml builder-config.yaml
RUN cd openvpnreceiver && mdatagen metadata.yaml
RUN cd pireceiver && mdatagen metadata.yaml

RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} GOARM=${TARGETARM} builder --config builder-config.yaml --name otelcol-${TARGETOS}-${TARGETARCH}

FROM scratch
COPY --from=build-env /otelcol/build/otelcol-* /

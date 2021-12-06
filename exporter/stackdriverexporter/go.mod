module github.com/open-telemetry/opentelemetry-collector-contrib/exporter/stackdriverexporter

go 1.14

require (
	cloud.google.com/go/container v1.0.0 // indirect
	cloud.google.com/go/monitoring v1.1.0 // indirect
	cloud.google.com/go/trace v1.0.0 // indirect
	contrib.go.opencensus.io/exporter/stackdriver v0.13.4
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace v0.2.2-0.20200728233621-2752da7eaab7
	github.com/census-instrumentation/opencensus-proto v0.3.0
	github.com/stretchr/testify v1.6.1
	go.opencensus.io v0.23.0
	go.opentelemetry.io/collector v0.9.0
	go.opentelemetry.io/otel v0.9.0
	go.uber.org/zap v1.15.0
	google.golang.org/api v0.61.0
	google.golang.org/genproto v0.0.0-20211118181313-81c1377c94b1
	google.golang.org/grpc v1.40.0
	google.golang.org/grpc/examples v0.0.0-20200728194956-1c32b02682df // indirect
	google.golang.org/protobuf v1.27.1
)

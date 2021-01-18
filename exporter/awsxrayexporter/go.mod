module github.com/open-telemetry/opentelemetry-collector-contrib/exporter/awsxrayexporter

go 1.14

replace github.com/open-telemetry/opentelemetry-collector-contrib/internal/common/awsxray => ../../internal/common/awsxray

require (
	github.com/aws/aws-sdk-go v1.36.28
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/common/awsxray v0.0.0-00010101000000-000000000000
	github.com/open-telemetry/opentelemetry-proto v0.4.0
	github.com/stretchr/testify v1.6.1
	go.opentelemetry.io/collector v0.9.0
	go.uber.org/zap v1.15.0
	golang.org/x/net v0.0.0-20201110031124-69a78807bb2b
	google.golang.org/grpc/examples v0.0.0-20200728194956-1c32b02682df // indirect
)

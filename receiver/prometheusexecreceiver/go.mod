module github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusexecreceiver

go 1.14

require (
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51
	github.com/prometheus/common v0.18.0
	github.com/prometheus/prometheus v1.8.2-0.20210217141258-a6be548dbc17
	github.com/satori/go.uuid v1.2.1-0.20181028125025-b2ce2384e17b // indirect
	github.com/stretchr/testify v1.7.0
	go.opentelemetry.io/collector v0.22.0
	go.uber.org/zap v1.16.0
	sigs.k8s.io/structured-merge-diff/v3 v3.0.0 // indirect
)

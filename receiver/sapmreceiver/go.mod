module github.com/open-telemetry/opentelemetry-collector-contrib/receiver/sapmreceiver

go 1.14

require (
	github.com/gorilla/mux v1.8.0
	github.com/hashicorp/go-plugin v1.3.0 // indirect
	github.com/jaegertracing/jaeger v1.21.0
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/common v0.0.0-00010101000000-000000000000
	github.com/open-telemetry/opentelemetry-proto v0.4.0
	github.com/sectioneight/md-to-godoc v0.0.0-20161108233149-55e43be6c335 // indirect
	github.com/signalfx/sapm-proto v0.5.3
	github.com/stretchr/testify v1.7.0
	go.opencensus.io v0.22.6
	go.opentelemetry.io/collector v0.21.0
	go.uber.org/zap v1.16.0
	gopkg.in/jcmturner/gokrb5.v7 v7.5.0 // indirect
)

replace github.com/open-telemetry/opentelemetry-collector-contrib/internal/common => ../../internal/common

replace github.com/googleapis/gnostic => github.com/googleapis/gnostic v0.3.1

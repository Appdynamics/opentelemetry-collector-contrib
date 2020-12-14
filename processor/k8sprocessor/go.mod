module github.com/open-telemetry/opentelemetry-collector-contrib/processor/k8sprocessor

go 1.14

require (
	github.com/census-instrumentation/opencensus-proto v0.3.0
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/common v0.0.0-00010101000000-000000000000 // indirect
	github.com/stretchr/testify v1.6.1
	go.opencensus.io v0.22.5
	go.opentelemetry.io/collector v0.16.0
	go.uber.org/zap v1.16.0
	k8s.io/api v0.19.2
	k8s.io/apimachinery v0.19.2
	k8s.io/client-go v0.19.2
	sigs.k8s.io/structured-merge-diff/v3 v3.0.0 // indirect
)

replace github.com/open-telemetry/opentelemetry-collector-contrib/internal/common => ../../internal/common

// Yet another hack that we need until kubernetes client moves to the new github.com/googleapis/gnostic
replace github.com/googleapis/gnostic => github.com/googleapis/gnostic v0.3.1

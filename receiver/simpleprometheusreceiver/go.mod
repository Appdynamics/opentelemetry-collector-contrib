module github.com/open-telemetry/opentelemetry-collector-contrib/receiver/simpleprometheusreceiver

go 1.14

require (
	github.com/go-sql-driver/mysql v1.5.0 // indirect
	github.com/prometheus/common v0.15.0
	github.com/prometheus/prometheus v1.8.2-0.20201105135750-00f16d1ac3a4
	github.com/satori/go.uuid v1.2.1-0.20181028125025-b2ce2384e17b // indirect
	github.com/stretchr/testify v1.6.1
	go.opentelemetry.io/collector v0.17.0
	go.uber.org/zap v1.16.0
	google.golang.org/grpc/examples v0.0.0-20200728194956-1c32b02682df // indirect
	k8s.io/client-go v0.19.2
	sigs.k8s.io/structured-merge-diff/v3 v3.0.0 // indirect
)

// Yet another hack that we need until kubernetes client moves to the new github.com/googleapis/gnostic
replace github.com/googleapis/gnostic => github.com/googleapis/gnostic v0.3.1

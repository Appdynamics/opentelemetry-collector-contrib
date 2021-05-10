module github.com/open-telemetry/opentelemetry-collector-contrib/internal/common

go 1.14

require (
	github.com/Microsoft/go-winio v0.4.14 // indirect
	github.com/cenkalti/backoff v2.2.1+incompatible
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/docker/docker v17.12.1-ce+incompatible
	github.com/docker/go-connections v0.4.0
	github.com/docker/go-units v0.4.0 // indirect
	github.com/google/go-cmp v0.5.1 // indirect
	github.com/googleapis/gnostic v0.4.1 // indirect
	github.com/imdario/mergo v0.3.10 // indirect
	github.com/json-iterator/go v1.1.10 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.0.1 // indirect
	github.com/stretchr/testify v1.6.1
	golang.org/x/net v0.0.0-20200625001655-4c5254603344 // indirect
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d // indirect
	golang.org/x/sys v0.0.0-20200625212154-ddb9806d33ae // indirect
	golang.org/x/text v0.3.3 // indirect
	golang.org/x/time v0.0.0-20200630173020-3af7569d3a1e // indirect
	google.golang.org/appengine v1.6.6 // indirect
	google.golang.org/protobuf v1.25.0 // indirect
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
	gopkg.in/yaml.v3 v3.0.0-20200603094226-e3079894b1e8 // indirect
	k8s.io/client-go v0.18.8
	k8s.io/utils v0.0.0-20200724153422-f32512634ab7 // indirect
)

// Yet another hack that we need until kubernetes client moves to the new github.com/googleapis/gnostic
replace github.com/googleapis/gnostic => github.com/googleapis/gnostic v0.3.1

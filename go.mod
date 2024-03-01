module github.com/devlibx/gox-metrics

go 1.16

replace github.com/cactus/go-statsd-client => github.com/cactus/go-statsd-client v3.1.0+incompatible

require (
	github.com/cactus/go-statsd-client v3.1.0+incompatible
	github.com/b2pacific/gox-base v0.0.0-20240301210626-54d926d9c8ec
	github.com/devlibx/gox-http v0.0.79
	github.com/m3db/prometheus_client_golang v0.8.1 // indirect
	github.com/m3db/prometheus_client_model v0.1.0 // indirect
	github.com/m3db/prometheus_common v0.1.0 // indirect
	github.com/m3db/prometheus_procfs v0.8.1 // indirect
	github.com/opentracing/opentracing-go v1.2.0
	github.com/stretchr/testify v1.8.4
	github.com/twmb/murmur3 v1.1.5 // indirect
	github.com/uber-go/tally v3.4.0+incompatible
	go.uber.org/ratelimit v0.2.0
	gopkg.in/DataDog/dd-trace-go.v1 v1.58.1
)

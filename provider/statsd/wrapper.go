package statsd

import (
	statsd3 "github.com/cactus/go-statsd-client/statsd"
	"github.com/devlibx/gox-base/errors"
	"github.com/devlibx/gox-base/metrics"
	"github.com/uber-go/tally"
	_ "github.com/uber-go/tally"
	statsd2 "github.com/uber-go/tally/statsd"
	"io"
	"time"
)

// Wrapper of tally timer
type timer struct {
	timer tally.Timer
}

func (t *timer) Record(value time.Duration) {
	t.timer.Record(value)
}

func (t *timer) Start() metrics.Stopwatch {
	return t.timer.Start()
}

// Wrapper of tally histogram
type histogram struct {
	histogram tally.Histogram
}

func (h *histogram) RecordValue(value float64) {
	h.histogram.RecordValue(value)
}

func (h *histogram) RecordDuration(value time.Duration) {
	h.histogram.RecordDuration(value)
}

func (h *histogram) Start() metrics.Stopwatch {
	return h.histogram.Start()
}

// Wrapper of tally scope class
type statsdMetrics struct {
	scope  tally.Scope
	closer io.Closer
}

func (s *statsdMetrics) Counter(name string) metrics.Counter {
	return s.scope.Counter(name)
}

func (s *statsdMetrics) Gauge(name string) metrics.Gauge {
	return s.scope.Gauge(name)
}

func (s *statsdMetrics) Timer(name string) metrics.Timer {
	return &timer{timer: s.scope.Timer(name)}
}

func (s *statsdMetrics) Histogram(name string, buckets metrics.Buckets) metrics.Histogram {
	return &histogram{histogram: s.scope.Histogram(name, buckets)}
}

func (s *statsdMetrics) Tagged(tags map[string]string) metrics.Scope {
	return &statsdMetrics{scope: s.scope.Tagged(tags)}
}

func (s *statsdMetrics) SubScope(name string) metrics.Scope {
	return &statsdMetrics{scope: s.scope.SubScope(name)}
}

func (s *statsdMetrics) Capabilities() metrics.Capabilities {
	return s.scope.Capabilities()
}

func NewRootScope(config metrics.Config) (metrics.Scope, error) {

	// Build client
	statsdClient, err := statsd3.NewBufferedClient(
		config.Statsd.Address,
		"stats",
		time.Duration(config.Statsd.FlushIntervalMs)*time.Millisecond,
		config.Statsd.FlushBytes,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create statsd client: config=%v", config)
	}

	// Create a new Statsd reported
	opts := statsd2.Options{}
	reporter := statsd2.NewReporter(statsdClient, opts)

	// Create tally specific scope object to use
	scope, closer := tally.NewRootScope(
		tally.ScopeOptions{
			Prefix:   config.Prefix,
			Tags:     map[string]string{},
			Reporter: reporter,
		},
		time.Duration(config.ReportingIntervalMs)*time.Millisecond,
	)

	return &statsdMetrics{scope:  scope, closer: closer,}, nil
}

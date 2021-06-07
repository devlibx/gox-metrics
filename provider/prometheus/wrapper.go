package prometheus

import (
	"github.com/devlibx/gox-base/metrics"
	"github.com/uber-go/tally"
	promreporter "github.com/uber-go/tally/prometheus"
	"io"
	"net/http"
	"strings"
	"sync"
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
type prometheusMetrics struct {
	scope     tally.Scope
	closer    io.Closer
	closeOnce sync.Once
	reporter  metrics.Reporter
}

func (s *prometheusMetrics) Counter(name string) metrics.Counter {
	return s.scope.Counter(name)
}

func (s *prometheusMetrics) Gauge(name string) metrics.Gauge {
	return s.scope.Gauge(name)
}

func (s *prometheusMetrics) Timer(name string) metrics.Timer {
	return &timer{timer: s.scope.Timer(name)}
}

func (s *prometheusMetrics) Histogram(name string, buckets metrics.Buckets) metrics.Histogram {
	return &histogram{histogram: s.scope.Histogram(name, buckets)}
}

func (s *prometheusMetrics) Tagged(tags map[string]string) metrics.Scope {
	return &prometheusMetrics{scope: s.scope.Tagged(tags)}
}

func (s *prometheusMetrics) SubScope(name string) metrics.Scope {
	return &prometheusMetrics{scope: s.scope.SubScope(name)}
}

func (s *prometheusMetrics) Capabilities() metrics.Capabilities {
	return s.scope.Capabilities()
}

func (s *prometheusMetrics) Stop() error {
	s.closeOnce.Do(func() {
		_ = s.closer.Close()
	})
	return nil
}

func (s *prometheusMetrics) HTTPHandler() http.Handler {
	return s.reporter.HTTPHandler()
}

func NewRootScope(config metrics.Config) (metrics.ClosableScope, error) {

	reporter := promreporter.NewReporter(promreporter.Options{})

	// Replace "." to "_" - prometheus does not work with "."
	prefix := strings.ReplaceAll(config.Prefix, ".", "_")

	// Create tally specific scope object to use
	scope, closer := tally.NewRootScope(
		tally.ScopeOptions{
			Prefix:         prefix,
			Tags:           map[string]string{},
			CachedReporter: reporter,
			Separator:      promreporter.DefaultSeparator,
		},
		time.Duration(config.ReportingIntervalMs)*time.Millisecond,
	)

	return &prometheusMetrics{scope: scope, closer: closer, closeOnce: sync.Once{}, reporter: reporter}, nil
}

package prometheus

import (
	"fmt"
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
type PrometheusMetrics struct {
	Scope     tally.Scope
	closer    io.Closer
	closeOnce sync.Once
	reporter  metrics.Reporter
}

func (s *PrometheusMetrics) Counter(name string) metrics.Counter {
	return s.Scope.Counter(name)
}

func (s *PrometheusMetrics) Gauge(name string) metrics.Gauge {
	return s.Scope.Gauge(name)
}

func (s *PrometheusMetrics) Timer(name string) metrics.Timer {
	return &timer{timer: s.Scope.Timer(name)}
}

func (s *PrometheusMetrics) Histogram(name string, buckets metrics.Buckets) metrics.Histogram {
	return &histogram{histogram: s.Scope.Histogram(name, buckets)}
}

func (s *PrometheusMetrics) Tagged(tags map[string]string) metrics.Scope {
	return &PrometheusMetrics{Scope: s.Scope.Tagged(tags)}
}

func (s *PrometheusMetrics) SubScope(name string) metrics.Scope {
	return &PrometheusMetrics{Scope: s.Scope.SubScope(name)}
}

func (s *PrometheusMetrics) Capabilities() metrics.Capabilities {
	return s.Scope.Capabilities()
}

func (s *PrometheusMetrics) Stop() error {
	s.closeOnce.Do(func() {
		_ = s.closer.Close()
	})
	return nil
}

func (s *PrometheusMetrics) HTTPHandler() http.Handler {
	return s.reporter.HTTPHandler()
}

func NewRootScope(config metrics.Config) (metrics.ClosableScope, error) {

	reporter := promreporter.NewReporter(promreporter.Options{
		OnRegisterError: func(err error) {
			fmt.Printf("error registering prometheus reporter: %v", err)
		},
	})

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

	return &PrometheusMetrics{Scope: scope, closer: closer, closeOnce: sync.Once{}, reporter: reporter}, nil
}

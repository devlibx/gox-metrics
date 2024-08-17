package multi

import (
	"github.com/devlibx/gox-base/v2/errors"
	"github.com/devlibx/gox-base/v2/metrics"
	"github.com/devlibx/gox-metrics/v2/provider/prometheus"
	"github.com/devlibx/gox-metrics/v2/provider/statsd"
	"io"
	"net/http"
	"sync"
	"time"
)

var postfix = ""

func NewRootScope(config metrics.Config) (metrics.ClosableScope, error) {
	s, err := statsd.NewRootScope(config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build statsd root scope")
	}

	p, err := prometheus.NewRootScope(config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build prometheus root scope")
	}

	m := &multiMetrics{
		statsdMetrics:              s,
		prometheusMetrics:          p,
		statsdMetricsCloseOnce:     sync.Once{},
		prometheusMetricsCloseOnce: sync.Once{},
	}

	return m, nil
}

func (s *multiMetrics) HTTPHandler() http.Handler {
	reporter := s.prometheusMetrics.(metrics.Reporter)
	return reporter.HTTPHandler()
}

type multiMetrics struct {
	statsdMetrics              metrics.Scope
	prometheusMetrics          metrics.Scope
	statsdMetricsCloseOnce     sync.Once
	prometheusMetricsCloseOnce sync.Once
	statsdMetricsCloser        io.Closer
	prometheusMetricsCloser    io.Closer
}

func (m *multiMetrics) Counter(name string) metrics.Counter {
	c := compositeCounter{
		statsd:     m.statsdMetrics.Counter(name),
		prometheus: m.prometheusMetrics.Counter(name + postfix),
	}
	return &c
}

func (m *multiMetrics) Gauge(name string) metrics.Gauge {
	c := compositeGauge{
		statsd:     m.statsdMetrics.Gauge(name),
		prometheus: m.prometheusMetrics.Gauge(name + postfix),
	}
	return &c
}

func (m *multiMetrics) Timer(name string) metrics.Timer {
	c := compositeTimer{
		statsd:     m.statsdMetrics.Timer(name),
		prometheus: m.prometheusMetrics.Timer(name + postfix),
	}
	return &c
}

func (m *multiMetrics) Histogram(name string, buckets metrics.Buckets) metrics.Histogram {
	c := compositeHistogram{
		statsd:     m.statsdMetrics.Histogram(name, buckets),
		prometheus: m.prometheusMetrics.Histogram(name+postfix, buckets),
	}
	return &c
}

func (m *multiMetrics) Tagged(tags map[string]string) metrics.Scope {
	return &multiMetrics{
		statsdMetrics:     m.statsdMetrics.Tagged(tags),
		prometheusMetrics: m.prometheusMetrics.Tagged(tags),
	}
}

func (m *multiMetrics) SubScope(name string) metrics.Scope {
	return &multiMetrics{
		statsdMetrics:     m.statsdMetrics.SubScope(name),
		prometheusMetrics: m.prometheusMetrics.SubScope(name + postfix),
	}
}

func (m *multiMetrics) Capabilities() metrics.Capabilities {
	return m.prometheusMetrics.Capabilities()
}

func (m *multiMetrics) Stop() error {
	m.statsdMetricsCloseOnce.Do(func() {
		closable := m.statsdMetricsCloser.(metrics.ClosableScope)
		_ = closable.Stop()
	})
	m.prometheusMetricsCloseOnce.Do(func() {
		closable := m.prometheusMetrics.(metrics.ClosableScope)
		_ = closable.Stop()
	})
	return nil
}

type compositeCounter struct {
	statsd     metrics.Counter
	prometheus metrics.Counter
}

func (c *compositeCounter) Inc(delta int64) {
	c.statsd.Inc(delta)
	c.prometheus.Inc(delta)
}

type compositeGauge struct {
	statsd     metrics.Gauge
	prometheus metrics.Gauge
}

func (c *compositeGauge) Update(value float64) {
	c.statsd.Update(value)
	c.prometheus.Update(value)
}

type compositeTimer struct {
	statsd     metrics.Timer
	prometheus metrics.Timer
}

func (c *compositeTimer) Record(value time.Duration) {
	c.statsd.Record(value)
	c.prometheus.Record(value)
}

func (c *compositeTimer) Start() metrics.Stopwatch {
	sw1 := c.statsd.Start()
	sw2 := c.prometheus.Start()
	return &compositeStopwatch{
		statsd:     sw1,
		prometheus: sw2,
	}
}

type compositeStopwatch struct {
	statsd     metrics.Stopwatch
	prometheus metrics.Stopwatch
}

func (c *compositeStopwatch) Stop() {
	c.statsd.Stop()
	c.prometheus.Stop()
}

type compositeHistogram struct {
	statsd     metrics.Histogram
	prometheus metrics.Histogram
}

func (c *compositeHistogram) RecordValue(value float64) {
	c.statsd.RecordValue(value)
	c.prometheus.RecordValue(value)
}

func (c *compositeHistogram) RecordDuration(value time.Duration) {
	c.statsd.RecordDuration(value)
	c.prometheus.RecordDuration(value)
}

func (c *compositeHistogram) Start() metrics.Stopwatch {
	sw1 := c.statsd.Start()
	sw2 := c.prometheus.Start()
	return &compositeStopwatch{
		statsd:     sw1,
		prometheus: sw2,
	}
}

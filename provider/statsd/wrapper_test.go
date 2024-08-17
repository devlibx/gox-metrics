package statsd

import (
	"fmt"
	"github.com/devlibx/gox-base/v2/metrics"
	"github.com/devlibx/gox-base/v2/test"
	"github.com/devlibx/gox-base/v2/util"
	"github.com/devlibx/gox-metrics/v2"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func init() {
	gox_metrics.Init()
}

func TestStatsd_Wrapper(t *testing.T) {
	testStatsd := gox_metrics.GetTestStringConfig("real.statsd")
	if util.IsStringEmpty(testStatsd) {
		t.Skip("set -real.statsd=true to enable test")
	}

	config := metrics.Config{
		Prefix:              "some",
		ReportingIntervalMs: 10,
		Statsd: metrics.StatsdConfig{
			Address:         "127.0.0.1:8125",
			FlushIntervalMs: 10,
			FlushBytes:      1440,
			Properties:      map[string]interface{}{"comma_appended_stats_reporter": true},
		},
	}
	statsdService, err := NewRootScope(config)
	assert.NoError(t, err)
	defer statsdService.Stop()

	counter := statsdService.Counter("some_counter")
	go func() {
		for i := 0; i < 100000; i++ {
			counter.Inc(1)
			time.Sleep(1 * time.Second)
			fmt.Println("submitting statsd counter")
		}
	}()

	time.Sleep(10 * time.Second)
}

func TestUsageWithCf(t *testing.T) {
	testStatsd := gox_metrics.GetTestStringConfig("real.statsd")
	if util.IsStringEmpty(testStatsd) {
		t.Skip("set -real.statsd=true to enable test")
	}

	c := metrics.Config{
		Prefix:              "some",
		ReportingIntervalMs: 10,
		Statsd: metrics.StatsdConfig{
			Address:         "127.0.0.1:8125",
			FlushIntervalMs: 10,
			FlushBytes:      1440,
		},
	}
	m, err := NewRootScope(c)
	assert.NoError(t, err)
	defer m.Stop()

	cf, _ := test.MockCf(t, m)

	counter := cf.Metric().Counter("TestUsageWithCf_Counter")
	go func() {
		for i := 0; i < 100000; i++ {
			counter.Inc(1)
			time.Sleep(1 * time.Second)
			fmt.Println("submitting statsd counter")
		}
	}()

	time.Sleep(10 * time.Second)
}

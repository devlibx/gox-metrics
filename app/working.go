package main

import (
	"fmt"
	"github.com/devlibx/gox-base/metrics"
	"github.com/devlibx/gox-metrics/provider/statsd"
	"math/rand"
	"os"
	"time"
)

//goland:noinspection GoUnreachableCode
func main() {
	var err error

	host := os.Getenv("__STATSD__")
	scope, err := statsd.NewRootScope(metrics.Config{
		Prefix:              "test_metric_stage",
		ReportingIntervalMs: 10000,
		Statsd: metrics.StatsdConfig{
			Address:         host,
			FlushIntervalMs: 10,
			FlushBytes:      1400 * 1000,
			StatsReporter:   statsd.NewCommaPerpetratedStatsReporter(true),
			Properties:      map[string]interface{}{"comma_perpetrated_stats_reporter": true},
		},
	})

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	i := 0
	for {
		ctr := scope.Tagged(map[string]string{"key": fmt.Sprintf("k_%d", r1.Intn(30)), "status": fmt.Sprintf("%d", r1.Intn(30))}).Counter("test_11")
		ctr.Inc(1)
		i++
	}

	ctr := scope.Tagged(map[string]string{"key": "111", "status": "ok"}).Counter("test_11")
	_ = ctr

	c := scope.Tagged(map[string]string{"key": "111", "status": "ok"}).Counter("test_11")
	g := scope.Tagged(map[string]string{"key": "111", "status": "ok"}).Gauge("g_test_11")
	if err == nil {
		for i := 0; i < 1000; i++ {
			c.Inc(1)
			g.Update(float64(r1.Intn(30)))
			// time.Sleep(100 * time.Millisecond)
		}
	}
	time.Sleep(10 * time.Second)
}

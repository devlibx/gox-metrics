package main

import (
	"fmt"
	"github.com/devlibx/gox-base/metrics"
	"github.com/devlibx/gox-metrics/provider/statsd"
	"go.uber.org/ratelimit"
	"os"
	"time"
)

//goland:noinspection GoUnreachableCode
func main() {
	host := os.Getenv("__STATSD__")
	scope, err := statsd.NewRootScope(metrics.Config{
		Prefix:              "test_metric_stage",
		ReportingIntervalMs: 1000,
		Statsd: metrics.StatsdConfig{
			Address:         host,
			FlushIntervalMs: 1000,
			FlushBytes:      1400 * 1000,
			StatsReporter:   statsd.NewCommaPerpetratedStatsReporter(true),
			Properties:      map[string]interface{}{"comma_perpetrated_stats_reporter": false},
		},
	})
	if err != nil {
		panic(err)
	}

	rl := ratelimit.New(200)
	for {
		rl.Take()
		ctr := scope.Tagged(map[string]string{"key": fmt.Sprintf("k_%d", 1), "status": fmt.Sprintf("%d", 1)}).Counter("test_16")
		ctr.Inc(1)
	}
	time.Sleep(10 * time.Second)
}

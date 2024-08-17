package statsd

import (
	"fmt"
	"github.com/devlibx/gox-base/v2/util"
	"github.com/uber-go/tally"
	"sort"
	"strings"
	"time"
)

type printStatsReporter struct {
	reporter  tally.StatsReporter
	printLogs bool
}

func NewCommaPerpetratedStatsReporter(printLogs bool) CustomStatsReporter {
	return &printStatsReporter{printLogs: printLogs}
}

func getCommaSeparatedKeyValuePairs(keyValuePairs map[string]string) string {
	keys := make([]string, 0, len(keyValuePairs))
	for k := range keyValuePairs {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var returnValue []string
	for _, key := range keys {
		if val, ok := keyValuePairs[key]; ok && !util.IsStringEmpty(val) {
			returnValue = append(returnValue, fmt.Sprintf("%s=%s", key, val))
		}
	}
	return strings.Join(returnValue, ",")
}

func (r *printStatsReporter) Init(reporter tally.StatsReporter) error {
	r.reporter = reporter
	return nil
}

func (r *printStatsReporter) ReportCounter(name string, v map[string]string, value int64) {
	f := fmt.Sprintf("%s,%s", name, getCommaSeparatedKeyValuePairs(v))
	if r.printLogs {
		fmt.Printf("Report counter = %s, value=%d \n", f, value)
	}
	r.reporter.ReportCounter(f, v, value)
}

func (r *printStatsReporter) ReportGauge(name string, v map[string]string, value float64) {
	f := fmt.Sprintf("%s,%s", name, getCommaSeparatedKeyValuePairs(v))
	if r.printLogs {
		fmt.Printf("Report gauge = %s, value=%f \n", f, value)
	}
	r.reporter.ReportGauge(f, v, value)
}

func (r *printStatsReporter) ReportTimer(name string, v map[string]string, interval time.Duration) {
	f := fmt.Sprintf("%s,%s", name, getCommaSeparatedKeyValuePairs(v))
	if r.printLogs {
		fmt.Printf("Report timer = %s, interval=%d \n", f, interval.Milliseconds())
	}
	r.reporter.ReportTimer(f, v, interval)
}

func (r *printStatsReporter) ReportHistogramValueSamples(
	name string,
	v map[string]string,
	b tally.Buckets,
	bucketLowerBound,
	bucketUpperBound float64,
	samples int64,
) {
	f := fmt.Sprintf("%s,%s", name, getCommaSeparatedKeyValuePairs(v))
	if r.printLogs {
		fmt.Printf("Report histogram = %s \n", f)
	}
	r.reporter.ReportHistogramValueSamples(f, v, b, bucketLowerBound, bucketUpperBound, samples)
}

func (r *printStatsReporter) ReportHistogramDurationSamples(
	name string,
	v map[string]string,
	b tally.Buckets,
	bucketLowerBound,
	bucketUpperBound time.Duration,
	samples int64,
) {
	f := fmt.Sprintf("%s,%s", name, getCommaSeparatedKeyValuePairs(v))
	if r.printLogs {
		fmt.Printf("Report histogram = %s \n", f)
	}
	r.reporter.ReportHistogramDurationSamples(f, v, b, bucketLowerBound, bucketUpperBound, samples)
}

func (r *printStatsReporter) Capabilities() tally.Capabilities {
	return r.reporter.Capabilities()
}

func (r *printStatsReporter) Reporting() bool {
	r.reporter.Flush()
	return true
}

func (r *printStatsReporter) Tagging() bool {
	return false
}

func (r *printStatsReporter) Flush() {
	r.reporter.Flush()
	if r.printLogs {
		fmt.Printf("Flush Metric \n")
	}
}

package stats

import (
	"github.com/uber-go/tally"
	"strings"
)

type TallyScopeWrapper struct {
	Scope tally.Scope
}

func (t TallyScopeWrapper) fixNames(name string) string {
	return strings.ReplaceAll(name, "-", "_")
}

func (t TallyScopeWrapper) Counter(name string) tally.Counter {
	return t.Scope.Counter(t.fixNames(name))
}

func (t TallyScopeWrapper) Gauge(name string) tally.Gauge {
	return t.Scope.Gauge(t.fixNames(name))
}

func (t TallyScopeWrapper) Timer(name string) tally.Timer {
	return t.Scope.Timer(t.fixNames(name))
}

func (t TallyScopeWrapper) Histogram(name string, buckets tally.Buckets) tally.Histogram {
	return t.Scope.Histogram(t.fixNames(name), buckets)
}

func (t TallyScopeWrapper) Tagged(tags map[string]string) tally.Scope {
	m := make(map[string]string)
	for k, v := range tags {
		m[t.fixNames(k)] = v
	}
	return TallyScopeWrapper{t.Scope.Tagged(m)}
}

func (t TallyScopeWrapper) SubScope(name string) tally.Scope {
	return TallyScopeWrapper{t.Scope.SubScope(t.fixNames(name))}
}

func (t TallyScopeWrapper) Capabilities() tally.Capabilities {
	return t.Scope.Capabilities()
}

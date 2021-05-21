package gox_metrics

import (
	"flag"
	"sync"
)

var testStatsd = ""
var doOnce = sync.Once{}

func Init() {
	doOnce.Do(func() {
		flag.StringVar(&testStatsd, "real.statsd", "", "run tests for statsd")
	})
}

func GetTestStringConfig(name string) string {
	switch name {
	case "real.statsd":
		return testStatsd
	}
	return ""
}

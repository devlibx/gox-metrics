package stats

import (
	"context"
	"fmt"
	"github.com/devlibx/gox-base/config"
	"github.com/devlibx/gox-base/metrics"
	"github.com/devlibx/gox-base/util"
	httpCommand "github.com/devlibx/gox-http/command/http"
	"github.com/devlibx/gox-metrics/provider/multi"
	"github.com/devlibx/gox-metrics/provider/prometheus"
	"github.com/devlibx/gox-metrics/provider/statsd"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/opentracer"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	goLogger "log"
	"net/http"
)

func NewMetricService(metricConfig *metrics.Config, appConfig *config.App) (metrics.Scope, *MetricHandler, error) {
	var toRet metrics.Scope
	var err error
	if metricConfig.Enabled {

		// Setup default if missing & build metric
		metricConfig.SetupDefaults()

		// Use correct type of metric exporter
		if metricConfig.EnablePrometheus && metricConfig.EnableStatsd {
			metricConfig.Statsd.StatsReporter = statsd.NewCommaPerpetratedStatsReporter(false)
			toRet, err = multi.NewRootScope(*metricConfig)
		} else if metricConfig.EnableStatsd {
			metricConfig.Statsd.StatsReporter = statsd.NewCommaPerpetratedStatsReporter(false)
			toRet, err = statsd.NewRootScope(*metricConfig)
		} else {
			toRet, err = prometheus.NewRootScope(*metricConfig)
		}
	} else {
		toRet, err = metrics.NoOpMetric(), nil
	}

	// Setup a reporter
	var mh metrics.Reporter
	if reporter, ok := toRet.(metrics.Reporter); ok {
		mh = reporter
	}

	// Setup DD
	host := metricConfig.Tracing.DD.Host
	port := metricConfig.Tracing.DD.Port
	env := metricConfig.Tracing.DD.Env
	if !metricConfig.Tracing.Enabled {
		goLogger.Println("************* datadog tracer is not enabled *************")
	} else if util.IsStringEmpty(host) || port <= 0 {
		goLogger.Println("datadog tracer is not enabled - host or port not provided")
	} else {
		agentAddr := fmt.Sprintf("%s:%d", host, port)
		goLogger.Println("Setting datadog tracer", "host=", host, "post=", port, "url=", agentAddr)
		if util.IsStringEmpty(env) {
			env = "dev"
		}

		// Set global tracer
		t := opentracer.New(tracer.WithAgentAddr(agentAddr), tracer.WithRuntimeMetrics(), tracer.WithServiceName(appConfig.AppName), tracer.WithEnv(env))
		opentracing.SetGlobalTracer(t)
		//tracer.Start(tracer.WithRuntimeMetrics())

		// Override the opentracing wrapper for Gox-Http to work
		httpCommand.DefaultStartSpanFromContextFunc = func(ctx context.Context, operationName string, opts ...opentracing.StartSpanOption) (opentracing.Span, context.Context) {
			s, ctx := tracer.StartSpanFromContext(ctx, operationName)
			return &DdSpanWrapper{s}, ctx
		}
	}

	return toRet, &MetricHandler{MetricsReporter: mh}, err
}

type MetricHandler struct {
	MetricsReporter metrics.Reporter
}

func (mh *MetricHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if mh.MetricsReporter == nil {
		w.WriteHeader(200)
	} else {
		mh.MetricsReporter.HTTPHandler().ServeHTTP(w, r)
	}
}

func (mh *MetricHandler) HTTPHandler() http.Handler {
	return mh
}

// DdSpanWrapper wraps the opentracing Span to work with Data Dog
type DdSpanWrapper struct {
	s ddtrace.Span
}

func (f DdSpanWrapper) Finish() {
	f.s.Finish()
}

func (f DdSpanWrapper) FinishWithOptions(opts opentracing.FinishOptions) {
	f.s.Finish()
}

func (f DdSpanWrapper) Context() opentracing.SpanContext {
	return f.s.Context()
}

func (f DdSpanWrapper) SetOperationName(operationName string) opentracing.Span {
	f.s.SetOperationName(operationName)
	return f
}

func (f DdSpanWrapper) SetTag(key string, value interface{}) opentracing.Span {
	f.s.SetTag(key, value)
	return f
}

func (f DdSpanWrapper) LogFields(fields ...log.Field) {
}

func (f DdSpanWrapper) LogKV(alternatingKeyValues ...interface{}) {
	// f.LogKV(alternatingKeyValues...)
}

func (f DdSpanWrapper) SetBaggageItem(restrictedKey, value string) opentracing.Span {
	f.s.SetBaggageItem(restrictedKey, value)
	return f
}

func (f DdSpanWrapper) BaggageItem(restrictedKey string) string {
	return f.s.BaggageItem(restrictedKey)
}

func (f DdSpanWrapper) Tracer() opentracing.Tracer {
	return opentracing.GlobalTracer()
}

func (f DdSpanWrapper) LogEvent(event string) {
}

func (f DdSpanWrapper) LogEventWithPayload(event string, payload interface{}) {
}

func (f DdSpanWrapper) Log(data opentracing.LogData) {
}

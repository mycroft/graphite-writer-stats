package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

var (
	metricPathCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "metrics_path_total",
		Help: "The total number of metrics paths events",
	}, []string{"metric_path", "application", "application_type"})
	dataPointTometricErrorCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "metrics_error_total",
		Help: "The total number of bad parsed metrics paths",
	})
	metricPathDidNotMatchAnyRulesCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "metric_path_did_not_match_rules_total",
		Help: "The total number of bad parsed metrics paths",
	})
)

// GetPrometheusHTTPHandler returns the Prometheus http handler
func GetPrometheusHTTPHandler() http.Handler {
	return promhttp.Handler()
}

// IncDataPointToMetricErrorCounter increments the dataPointTometricErrorCount counter
func IncDataPointToMetricErrorCounter() {
	dataPointTometricErrorCount.Inc()
}

// IncMetricPathDidNotMatchAnyRules increments the dataPointTometricErrorCount counter
func IncMetricPathDidNotMatchAnyRules() {
	dataPointTometricErrorCount.Inc()
}

// IncMetricPathCounter increments an application counter based on its extracted metric
func IncMetricPathCounter(extractedMetric string, applicationName string, applicationType string) {
	metricPathCount.WithLabelValues(extractedMetric, applicationName, applicationType).Inc()
}

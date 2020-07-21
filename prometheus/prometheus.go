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
	metricProcessedEvents = promauto.NewCounter(prometheus.CounterOpts{
		Name: "metrics_processed_events",
		Help: "The total number of processed metrics",
	})
	metricLatestTimestampGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "metrics_timestamp_value",
		Help: "Lowest and Highest Timestamp processed",
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

// IncMetricProcessedEvents increments number of processed metrics in total
func IncMetricProcessedEvents() {
	metricProcessedEvents.Inc()
}

// SetMetricLatestTimestamp sets the latest processed timestamp
func SetMetricLatestTimestamp(ts float64) {
	metricLatestTimestampGauge.Set(ts)
}

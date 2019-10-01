package stats

import (
	"bytes"
	"github.com/criteo/graphite-writer-stats/prometheus"
	"go.uber.org/zap"
)

type Stats struct {
	Logger         *zap.Logger
	MetricMetadata MetricMetadata
}

func (stats *Stats) process(metricPath string) {
	metric := stats.getMetric(metricPath)
	if ce := stats.Logger.Check(zap.DebugLevel, "metrics"); ce != nil {
		ce.Write(zap.Any("metric", metric))
	}
	prometheus.IncMetricPathCounter(metric.ExtractedMetric, metric.ApplicationName, string(metric.ApplicationType))
}

func (stats *Stats) Process(dataPoint []byte) bool {
	metricPath, succeed := extractMetricPath(dataPoint)
	if succeed {
		stats.process(metricPath)
	} else {
		stats.Logger.Error("fail to convert datapoint to metricpath", zap.ByteString("datapoint", dataPoint))
		prometheus.IncDataPointToMetricErrorCounter()
	}
	return succeed
}

func extractMetricPath(metric []byte) (string, bool) {
	index := 0
	indexSpace := bytes.IndexByte(metric, ' ')
	if indexSpace != -1 {
		index = indexSpace
	}
	return string(metric[:index]), indexSpace != -1
}

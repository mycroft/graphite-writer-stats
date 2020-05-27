package stats

import (
	"bytes"
	"errors"

	"github.com/Shopify/sarama"
	"github.com/criteo/graphite-writer-stats/prometheus"
	"go.uber.org/zap"
)

type Stats struct {
	Logger         *zap.Logger
	MetricMetadata MetricMetadata
}

type Metric struct {
	// A Metric without its timestamp or value
	Path string
	Tags map[string]string
}

func BuildMetricFromMessage(message *sarama.ConsumerMessage) (Metric, error) {
	metric := Metric{}
	metric.Tags = make(map[string]string, 0)

	// Path
	index := 0
	indexSpace := bytes.IndexByte(message.Value, ' ')
	if indexSpace != -1 {
		index = indexSpace
	}

	if indexSpace == -1 {
		return metric, errors.New("Failed to parse metric name")
	}

	metric.Path = string(message.Value[:index])

	for _, kv := range message.Headers {
		metric.Tags[string(kv.Key)] = string(kv.Value)
	}

	return metric, nil
}

func (stats *Stats) process(metric Metric) {
	extractedMetric := stats.getMetric(metric.Path, metric.Tags)
	if ce := stats.Logger.Check(zap.DebugLevel, "metrics"); ce != nil {
		ce.Write(zap.Any("metric", metric.Path))
	}
	prometheus.IncMetricPathCounter(extractedMetric.ExtractedMetric, extractedMetric.ApplicationName, string(extractedMetric.ApplicationType))
}

func (stats *Stats) Process(message *sarama.ConsumerMessage) error {
	metric, err := BuildMetricFromMessage(message)
	if err != nil {
		return err
	}

	stats.process(metric)

	return nil
}

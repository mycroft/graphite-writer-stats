package stats

import (
	"bytes"
	"errors"

	"github.com/Shopify/sarama"
	"github.com/criteo/graphite-writer-stats/prometheus"
	"go.uber.org/zap"
)

// Stats is used to log messages & configuration
type Stats struct {
	MetricMetadata MetricMetadata
}

// The Metric structure only contains its path & tags; it doesn't store timestamp or value
type Metric struct {
	// A Metric without its timestamp or value
	Path string
	Tags map[string]string
}

// BuildMetricFromMessage is retrieving a Metric from a consumed message.
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

// Process a consumer kafka message (building metric & processing)
func (stats *Stats) Process(logger *zap.Logger, message *sarama.ConsumerMessage) error {
	metric, err := BuildMetricFromMessage(message)
	if err != nil {
		return err
	}

	extractedMetric := stats.getMetric(logger, metric.Path, metric.Tags)
	if ce := logger.Check(zap.DebugLevel, "metrics"); ce != nil {
		ce.Write(zap.Any("metric", metric.Path))
	}

	prometheus.IncMetricPathCounter(extractedMetric.ExtractedMetric, extractedMetric.ApplicationName, string(extractedMetric.ApplicationType))

	return nil
}

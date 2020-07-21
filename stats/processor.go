package stats

import (
	"bytes"
	"errors"
	"strconv"

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
	Path      string
	Tags      map[string]string
	Timestamp uint32
	Value     string
}

// BuildMetricFromMessage is retrieving a Metric from a consumed message.
func BuildMetricFromMessage(message *sarama.ConsumerMessage) (Metric, error) {
	var timestamp uint64
	var err error

	metric := Metric{}
	metric.Tags = make(map[string]string, 0)

	// Path
	index := 0
	indexSpace := bytes.IndexByte(message.Value, ' ')
	if indexSpace != -1 {
		index = indexSpace
	}

	if indexSpace == -1 {
		return metric, errors.New("Invalid indexSapce while parsing metric name")
	}

	metric.Path = string(message.Value[:index])

	for _, kv := range message.Headers {
		metric.Tags[string(kv.Key)] = string(kv.Value)
	}

	lastIndexSpace := bytes.LastIndexByte(message.Value, ' ')
	if lastIndexSpace == -1 || lastIndexSpace == indexSpace {
		return metric, errors.New("Invalid lastIndexSpace while parsing metric name")
	}

	timestamp, err = strconv.ParseUint(string(message.Value[lastIndexSpace+1:]), 10, 32)
	if err != nil {
		return metric, err
	}

	metric.Timestamp = uint32(timestamp)

	return metric, nil
}

// Process a consumer kafka message (building metric & processing)
func (stats *Stats) Process(logger *zap.Logger, message *sarama.ConsumerMessage) error {
	prometheus.IncMetricProcessedEvents()

	metric, err := BuildMetricFromMessage(message)
	if err != nil {
		prometheus.IncDataPointToMetricErrorCounter()
		return err
	}

	extractedMetric := stats.getMetric(logger, metric.Path, metric.Tags)
	if ce := logger.Check(zap.DebugLevel, "metrics"); ce != nil {
		ce.Write(zap.Any("metric", metric.Path))
	}

	if message.Offset%1000 == 0 && metric.Timestamp != 0 {
		prometheus.SetMetricLatestTimestamp(float64(metric.Timestamp))
	}

	prometheus.IncMetricPathCounter(extractedMetric.ExtractedMetric, extractedMetric.ApplicationName, string(extractedMetric.ApplicationType))

	return nil
}

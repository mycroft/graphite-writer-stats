package prometheus

import (
	"fmt"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// PartitionContext owns the prometheus metrics for the given partition
type PartitionContext struct {
	messagesConsumedCounter prometheus.Counter
	timeLagGauge            prometheus.Gauge
	offsetLagGauge          prometheus.Gauge
}

// KafkaConsumerMetrics owns the mandatory fields to be used while generating prometheus metrics
type KafkaConsumerMetrics struct {
	ID          string
	kafkaConfig *sarama.Config
}

// RegisterKafkaConsumerMetrics prepare the KafkaConsumerMetrics structure and register it in prometheus.
func RegisterKafkaConsumerMetrics(id string, kafkaConfig *sarama.Config) error {
	m := &KafkaConsumerMetrics{
		id,
		kafkaConfig,
	}
	return prometheus.Register(m)
}

// Describe the Prometheus metrics
func (k *KafkaConsumerMetrics) Describe(chan<- *prometheus.Desc) {
	return
}

// Collect the Prometheus metrics
func (k *KafkaConsumerMetrics) Collect(c chan<- prometheus.Metric) {
	metrics := k.kafkaConfig.MetricRegistry.GetAll()
	for key, value := range metrics {
		key = strings.ReplaceAll(key, "-", "_") // Because prometheus does not handle - in metric
		if value["count"] != nil {              // counter see MetricRegistry.GetAll() implementation
			value := value["count"].(int64)
			counter := prometheus.NewCounter(prometheus.CounterOpts{
				Namespace:   "kafka",
				Subsystem:   "consumer",
				Name:        key,
				Help:        "",
				ConstLabels: nil,
			})
			counter.Add(float64(value))
			c <- counter
		} else if value["value"] != nil { // gauge see MetricRegistry.GetAll() implementation
			gauge := prometheus.NewGauge(prometheus.GaugeOpts{
				Namespace:   "kafka",
				Subsystem:   "consumer",
				Name:        key,
				Help:        "",
				ConstLabels: nil,
			})
			gauge.Set(value["value"].(float64))
			c <- gauge
		}

	}
}

// NewPartitionContext prepares metrics for given partition
func NewPartitionContext(partition int32) PartitionContext {
	// promauto automatically register metrics

	messagesConsumedCounter := promauto.NewCounter(prometheus.CounterOpts{
		Namespace:   "kafka",
		Subsystem:   "consumer",
		Name:        "messages",
		Help:        "number of message consumed",
		ConstLabels: prometheus.Labels{"partition": fmt.Sprintf("%d", partition)},
	})

	timeLagGauge := promauto.NewGauge(prometheus.GaugeOpts{
		Namespace:   "kafka",
		Subsystem:   "consumer",
		Name:        "time_lag",
		Help:        "consumer time lag in seconds",
		ConstLabels: prometheus.Labels{"partition": fmt.Sprintf("%d", partition)},
	})

	offsetLagGauge := promauto.NewGauge(prometheus.GaugeOpts{
		Namespace:   "kafka",
		Subsystem:   "consumer",
		Name:        "offset_lag",
		Help:        "consumer offset lag",
		ConstLabels: prometheus.Labels{"partition": fmt.Sprintf("%d", partition)},
	})

	return PartitionContext{
		messagesConsumedCounter: messagesConsumedCounter,
		timeLagGauge:            timeLagGauge,
		offsetLagGauge:          offsetLagGauge,
	}
}

// MonitorConsumerLag computes & store lag data in metrics
func MonitorConsumerLag(context PartitionContext, claim sarama.ConsumerGroupClaim, message *sarama.ConsumerMessage) {
	context.messagesConsumedCounter.Inc()
	context.offsetLagGauge.Set(float64(claim.HighWaterMarkOffset() - message.Offset))
	context.timeLagGauge.Set(time.Now().Sub(message.Timestamp).Seconds())
}

// Destroy PartitionContext: Unregister its metrics from prometheus
func (pc PartitionContext) Destroy() {
	prometheus.Unregister(pc.messagesConsumedCounter)
	prometheus.Unregister(pc.timeLagGauge)
	prometheus.Unregister(pc.offsetLagGauge)
}

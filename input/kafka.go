package input

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/criteo/graphite-writer-stats/stats"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
)

// KafkaConfiguration stores kafka related configuration
type KafkaConfiguration struct {
	topics []string
	group  string
	config *sarama.Config
	oldest bool
}

// KafkaProcessor is the entry point to configuring, starting the process & managing metrics
type KafkaProcessor struct {
	logger      *zap.Logger
	kafkaConfig KafkaConfiguration
	client      sarama.Client
	consumer    sarama.ConsumerGroup
	ctx         context.Context
	cancel      context.CancelFunc
	wg          *sync.WaitGroup
	stats       stats.Stats
}

// The BrokerStatus has some broker status informations.
type BrokerStatus struct {
	ID              int32
	Addr            string
	Rack            string
	Connected       bool
	ConnectionError error
}

// The KafkaStatus is the list & statuses of the brokers
type KafkaStatus struct {
	Brokers []BrokerStatus
	Closed  bool
	Metrics map[string]map[string]interface{}
}

// CreateProcessor initialize the main KafkaProcessor structure
func CreateProcessor(logger *zap.Logger) *KafkaProcessor {
	return &KafkaProcessor{
		logger: logger,
	}
}

// SetupConsumer initializes the sarama client & consumer group, and returns a Processor ready to be used
func (processor *KafkaProcessor) SetupConsumer(brokers string, group string, topic string, oldest bool) error {
	var err error

	processor.kafkaConfig.group = group

	config := sarama.NewConfig()
	config.Version = sarama.V2_3_0_0
	if oldest {
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	}

	processor.kafkaConfig.config = config

	if topic == "" {
		return fmt.Errorf("can not use empty topic")
	}

	processor.kafkaConfig.topics = []string{topic}

	processor.ctx, processor.cancel = context.WithCancel(context.Background())

	processor.client, err = sarama.NewClient(strings.Split(brokers, ","), config)
	if err != nil {
		return fmt.Errorf("error creating kafka client: %v", err)
	}

	processor.consumer, err = sarama.NewConsumerGroupFromClient(group, processor.client)
	if err != nil {
		return fmt.Errorf("error creating kafka consumer: %v", err)
	}

	processor.wg = &sync.WaitGroup{}
	processor.wg.Add(1)

	return nil
}

// Run starts the consumer
func (processor *KafkaProcessor) Run(stats stats.Stats) {
	processor.stats = stats

	go func() {
		defer processor.wg.Done()
		for {
			fmt.Println(processor.ctx, processor.kafkaConfig.topics, processor)
			if err := processor.consumer.Consume(processor.ctx, processor.kafkaConfig.topics, processor); err != nil {
				processor.logger.Panic("Error from consumer", zap.Error(err))
			}
			// check if context was cancelled, signaling that the consumer should stop
			if processor.ctx.Err() != nil {
				return
			}
		}
	}()
}

// Wait endlessly for an unrecoverable error or a INT/TERM stopping signal.
func (processor *KafkaProcessor) Wait() {
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-processor.ctx.Done():
		processor.logger.Info("terminating: context cancelled")
	case <-sigterm:
		processor.logger.Info("terminating: via signal")
	}
}

// Close the Kafka topic before exiting.
func (processor *KafkaProcessor) Close() {
	processor.cancel()
	processor.wg.Wait()
	err := processor.client.Close()
	if err != nil {
		processor.logger.Panic("Error closing client: %v", zap.Error(err))
	}
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (processor *KafkaProcessor) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (processor *KafkaProcessor) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (processor *KafkaProcessor) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		processor.logger.Debug("Message", zap.ByteString("message", message.Value), zap.Time("timestamp", message.Timestamp), zap.ByteString("key", message.Key))
		processor.stats.Process(processor.logger, message)
		session.MarkMessage(message, "")
	}

	return nil
}

// Status queries the brokers and fills the KafkaStatus structure
func (processor *KafkaProcessor) Status() KafkaStatus {
	status := KafkaStatus{}
	for _, broker := range processor.client.Brokers() {
		brokerStatus := BrokerStatus{
			ID:   broker.ID(),
			Addr: broker.Addr(),
			Rack: broker.Rack(),
		}
		brokerStatus.Connected, brokerStatus.ConnectionError = broker.Connected()
		status.Brokers = append(status.Brokers, brokerStatus)
	}
	status.Closed = processor.client.Closed()
	status.Metrics = processor.kafkaConfig.config.MetricRegistry.GetAll()
	return status
}

// GetStatusHTTPHandler returns the http handler to query and retrieve the KafkaStatus structure in json format
func (processor *KafkaProcessor) GetStatusHTTPHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bytes, err := json.MarshalIndent(processor.Status(), "", "  ")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(bytes)
	})
}

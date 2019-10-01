package input

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/criteo/graphite-writer-stats/stats"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
)

type Kafka struct {
	logger *zap.Logger
	topic  []string
	client sarama.ConsumerGroup
	ctx    context.Context
	cancel context.CancelFunc
	wg     *sync.WaitGroup
	stats  stats.Stats
}

func SetupConsumer(logger *zap.Logger, oldest bool, group string, brokers string, topic string, stats stats.Stats) *Kafka {
	config := sarama.NewConfig()
	config.Version = sarama.V2_3_0_0
	if oldest {
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	}
	if topic == "" {
		logger.Panic("Error topic config empty")
	}
	ctx, cancel := context.WithCancel(context.Background())
	client, err := sarama.NewConsumerGroup(strings.Split(brokers, ","), group, config)
	if err != nil {
		logger.Panic("Error creating consumer group client", zap.Error(err))
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	return &Kafka{logger: logger, topic: []string{topic}, client: client, ctx: ctx, cancel: cancel, wg: wg, stats: stats}
}

func (kafka *Kafka) Run() {
	go func() {
		defer kafka.wg.Done()
		for {
			if err := kafka.client.Consume(kafka.ctx, kafka.topic, kafka); err != nil {
				kafka.logger.Panic("Error from consumer", zap.Error(err))
			}
			// check if context was cancelled, signaling that the consumer should stop
			if kafka.ctx.Err() != nil {
				return
			}
		}
	}()
}
func (kafka *Kafka) Wait() {
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-kafka.ctx.Done():
		kafka.logger.Info("terminating: context cancelled")
	case <-sigterm:
		kafka.logger.Info("terminating: via signal")
	}
}
func (kafka *Kafka) Close() {
	kafka.cancel()
	kafka.wg.Wait()
	err := kafka.client.Close()
	if err != nil {
		kafka.logger.Panic("Error closing client: %v", zap.Error(err))
	}
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (kafka *Kafka) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (kafka *Kafka) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (kafka *Kafka) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		kafka.logger.Debug("Message", zap.ByteString("message", message.Value), zap.Time("timestamp", message.Timestamp), zap.ByteString("key", message.Key))
		kafka.stats.Process(message.Value)
		session.MarkMessage(message, "")
	}

	return nil
}

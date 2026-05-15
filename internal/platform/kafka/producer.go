package kafka

import (
	"strings"

	"github.com/IBM/sarama"
	"github.com/tbikbulatov/go-pulseops/internal/platform/config"
)

type Producer struct {
	producer sarama.SyncProducer
}

func NewProducer(cfg config.KafkaConfig) (*Producer, error) {
	saramaCfg := sarama.NewConfig()
	saramaCfg.Producer.Return.Successes = true

	brokers := strings.Split(cfg.Brokers, ",")
	producer, err := sarama.NewSyncProducer(brokers, saramaCfg)
	if err != nil {
		return nil, err
	}

	return &Producer{producer: producer}, nil
}

func (p *Producer) Publish(topic string, key string, value []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(value),
	}

	_, _, err := p.producer.SendMessage(msg)
	return err
}

func (p *Producer) Close() error {
	return p.producer.Close()
}

package kafka

import (
	"context"

	kafka "github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(pBrokers []string, pTopic string) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(pBrokers...),
			Topic:    pTopic,
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (p *Producer) Publish(ctx context.Context, pkey, pValue []byte) error {
	return p.writer.WriteMessages(ctx, kafka.Message{
		Key:   pkey,
		Value: pValue,
	})
}

func (p *Producer) Close() error {
	return p.writer.Close()
}

package kafka

import (
	"context"

	kafka "github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
}

func (c *Consumer) Ping(ctx context.Context) error {

	_, lErr := c.reader.FetchMessage(ctx)
	if lErr != nil {
		return lErr
	}
	return nil
}

func NewConsumer(pBrokers []string, pTopic, pGroupID string) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: pBrokers,
			Topic:   pTopic,
			GroupID: pGroupID,
		}),
	}
}

func (c *Consumer) Consume(pCtx context.Context, pHandler func([]byte) error) error {
	for {
		lMsg, lErr := c.reader.ReadMessage(pCtx)
		if lErr != nil {
			return lErr
		}
		if lErr := pHandler(lMsg.Value); lErr != nil {
			return lErr
		}
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}

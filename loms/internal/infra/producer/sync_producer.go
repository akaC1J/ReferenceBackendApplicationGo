package producer

import (
	"fmt"

	"github.com/IBM/sarama"
)

func NewSyncProducer(brokers []string, opts ...Option) (sarama.SyncProducer, error) {
	config := PrepareConfig(opts...)

	syncProducer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("NewSyncProducer failed: %w", err)
	}

	return syncProducer, nil
}

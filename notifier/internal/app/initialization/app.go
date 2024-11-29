package initialization

import (
	"github.com/IBM/sarama"
	"log"
	"route256/notifier/internal/infra/consumer_group"
	"route256/notifier/internal/pkg/service/processors/eventsprocessor"
)

type App struct {
	ConsumerGroup *consumer_group.ConsumerGroup
}

func New(config *Config) (*App, error) {
	log.Println("[cart] Starting application initialization")
	eventService := eventsprocessor.New()
	group, err := consumer_group.NewConsumerGroup(config.KafkaConfig.Brokers, config.KafkaConfig.GroupId, []string{config.KafkaConfig.Topic},
		consumer_group.NewConsumerGroupHandler(eventService), eventService, consumer_group.WithOffsetsInitial(sarama.OffsetNewest))
	if err != nil {
		log.Fatal(err)
	}
	return &App{
		ConsumerGroup: group,
	}, nil
}

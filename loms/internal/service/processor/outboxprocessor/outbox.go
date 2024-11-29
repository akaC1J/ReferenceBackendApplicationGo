package outboxprocessor

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/jackc/pgx/v5"
	"log"
	"route256/loms/internal/model"
	"route256/loms/internal/repository/outboxrepository"
	transactionmanager "route256/loms/internal/service/transactionamanger"
	"strconv"
	"time"
)

var _ OutboxRepository = (*outboxrepository.Repository)(nil)

type OutboxRepository interface {
	GetPendingOutboxEvents(ctx context.Context, tx pgx.Tx, limit int) ([]*model.OutboxEvent, error)
	MarkOutboxEventProcessed(ctx context.Context, tx pgx.Tx, eventID int32) error
}

var _ TransactionManager = (*transactionmanager.TransactionManager)(nil)

type TransactionManager interface {
	Begin(ctx context.Context) (tx pgx.Tx, closeFunc func(), err error)
	Commit(tx pgx.Tx, closeFunc func()) error
}

type OutboxProcessor struct {
	repo     OutboxRepository
	producer sarama.SyncProducer
	interval time.Duration
	tm       TransactionManager
	topic    string
}

func NewOutboxProcessor(repo OutboxRepository, producer sarama.SyncProducer, interval time.Duration, tm TransactionManager, topicName string) *OutboxProcessor {
	return &OutboxProcessor{
		repo:     repo,
		producer: producer,
		interval: interval,
		tm:       tm,
		topic:    topicName,
	}
}

func (p *OutboxProcessor) Start(ctx context.Context) {
	return
	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Stopping OutboxProcessor")
			return
		case <-ticker.C:
			p.processEvents(ctx)
		}
	}
}

func (p *OutboxProcessor) processEvents(ctx context.Context) {
	tx, closeTx, err := p.tm.Begin(ctx)
	if err != nil {
		log.Printf("Failed to start transaction: %v", err)
		return
	}
	defer closeTx()
	events, err := p.repo.GetPendingOutboxEvents(ctx, tx, 10) // Получаем до 10 событий за раз
	if err != nil {
		log.Printf("Error fetching outbox events: %v", err)
		return
	}

	for _, event := range events {
		msg := p.createProducerMessage(event)
		partition, offset, err := p.producer.SendMessage(msg)
		if err != nil {
			log.Printf("Failed to send message: %v", err)
			return
		}
		log.Printf("Sent message to partition %d at offset %d", partition, offset)

		if err = p.repo.MarkOutboxEventProcessed(ctx, tx, event.ID); err != nil {
			log.Printf("Failed to mark event as processed: %v", err)
			return
		}
	}
	err = p.tm.Commit(tx, closeTx)
	if err != nil {
		log.Printf("Failed to commit transaction: %v", err)
		return
	}

}

func (p *OutboxProcessor) createProducerMessage(event *model.OutboxEvent) *sarama.ProducerMessage {
	value := sarama.ByteEncoder(event.Payload)
	key := sarama.StringEncoder(strconv.FormatInt(event.OrderID, 10))
	return &sarama.ProducerMessage{
		Topic:     p.topic,
		Key:       key,
		Value:     value,
		Timestamp: event.CreatedAt,
	}
}

package outboxrepository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"route256/loms/internal/model"
)

type Repository struct {
}

// NewRepository создает новый репозиторий с транзакцией
func NewRepository() *Repository {
	return &Repository{}
}

// SaveOutboxEvent сохраняет событие в таблицу outbox в рамках переданной транзакции
func (r *Repository) SaveOutboxEvent(ctx context.Context, tx pgx.Tx, event *model.OutboxEvent) error {
	repTx := New(tx)
	return repTx.SaveOutboxEvent(ctx, &SaveOutboxEventParams{
		OrderID: event.OrderID,
		Payload: event.Payload,
	})
}

// GetPendingOutboxEvents получает неотправленные события из таблицы outbox в рамках транзакции
func (r *Repository) GetPendingOutboxEvents(ctx context.Context, tx pgx.Tx, limit int) ([]*model.OutboxEvent, error) {
	repTx := New(tx)
	eventsFromDB, err := repTx.GetPendingOutboxEvents(ctx, int32(limit))
	if err != nil {
		return nil, fmt.Errorf("unable to get pending outbox events: %w", err)
	}

	// Преобразуем результат в нужный формат
	events := make([]*model.OutboxEvent, len(eventsFromDB))
	for i, event := range eventsFromDB {
		events[i] = &model.OutboxEvent{
			ID:        event.ID,
			OrderID:   event.OrderID,
			Payload:   event.Payload,
			CreatedAt: event.CreatedAt.Time,
		}
	}
	return events, nil
}

// MarkOutboxEventProcessed помечает событие как обработанное в рамках транзакции
func (r *Repository) MarkOutboxEventProcessed(ctx context.Context, tx pgx.Tx, eventID int32) error {
	repTx := New(tx)
	return repTx.MarkOutboxEventProcessed(ctx, eventID)
}

package model

import "time"

// OutboxEvent представляет событие, связанное с заказом, для обработки через механизм outbox.
type OutboxEvent struct {
	ID        int32     // Уникальный идентификатор события в таблице outbox
	OrderID   int64     // Идентификатор заказа, с которым связано событие
	Payload   string    // Дополнительные данные события в формате JSON
	CreatedAt time.Time // Временная метка создания события
	Processed bool      // Флаг обработки, указывает, было ли отправлено событие
}

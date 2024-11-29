package eventsprocessor

import (
	"log"
	"route256/notifier/internal/pkg/model"
)

type EventService struct {
}

func New() *EventService {
	return &EventService{}
}

func (ep *EventService) Process(event *model.OrderEvent) {
	switch event.State {
	case model.NEW:
		log.Printf("Заказ создан: %+v\n", event)

	case model.AWAITING_PAYMENT:
		log.Printf("Заказ ожидает оплаты: %+v\n", event)

	case model.FAILED:
		log.Printf("Заказ не оплачен: %+v\n", event)

	case model.PAYED:
		log.Printf("Заказ оплачен: %+v\n", event)

	case model.CANCELLED:
		log.Printf("Заказ отменен: %+v\n", event)
	}
}

func (ep *EventService) ProcessError(err error) {
	log.Printf("Ошибка обработки события: %v\n", err)
}

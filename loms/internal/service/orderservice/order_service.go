package orderservice

import (
	"context"
	"log"
	appErr "route256/loms/internal/errors"
	"route256/loms/internal/model"
	"route256/loms/internal/repository/orderrepository"
)

var _ Repository = (*orderrepository.Repository)(nil)

type Repository interface {
	SaveOrder(ctx context.Context, order *model.Order) (*model.Order, error)
	UpdateOrder(ctx context.Context, order *model.Order) error
	GetById(ctx context.Context, orderID int64) (*model.Order, error)
}

type StockService interface {
	Reserve(ctx context.Context, items []*model.Item) error
	ReserveRemove(ctx context.Context, items []*model.Item) error
	ReserveCancel(ctx context.Context, items []*model.Item) error
}

type Service struct {
	repository   Repository
	stockService StockService
}

func NewService(repository Repository, stockService StockService) *Service {
	return &Service{repository: repository, stockService: stockService}
}

func (s *Service) Create(ctx context.Context, order *model.Order) (orderID int64, err error) {
	_ = order.SetState(model.NEW)

	order, err = s.repository.SaveOrder(ctx, order)
	if err != nil {
		log.Printf("[order_service] Error saving order: %v", err)
		return 0, err
	}

	err = s.stockService.Reserve(ctx, order.Items)
	if err != nil {
		_ = order.SetState(model.FAILED)
		if updateErr := s.repository.UpdateOrder(ctx, order); updateErr != nil {
			log.Printf("[order_service] Error updating order state: %v", updateErr)
			return 0, updateErr
		}
		log.Printf("[order_service] Error reserving stock: %v", err)
		return 0, err
	}

	_ = order.SetState(model.AWAITING_PAYMENT)
	err = s.repository.UpdateOrder(ctx, order)
	if err != nil {
		log.Printf("[order_service] Error updating order state: %v", err)
		return 0, err
	}

	return order.ID, nil
}

func (s *Service) GetById(ctx context.Context, orderID int64) (*model.Order, error) {
	order, err := s.repository.GetById(ctx, orderID)
	if err != nil {
		log.Printf("[order_service] Error getting order: %v", err)
		return nil, err
	}
	return order, nil
}

func (s *Service) OrderPay(ctx context.Context, orderID int64) error {
	order, err := s.repository.GetById(ctx, orderID)
	if err != nil {
		log.Printf("[order_service] Error getting order: %v", err)
		return err
	}

	if order.State() != model.AWAITING_PAYMENT {
		log.Printf("[order_service] Invalid order state: %v", order.State())
		return appErr.ErrOrderState
	}

	err = s.stockService.ReserveRemove(ctx, order.Items)
	if err != nil {
		log.Printf("[order_service] Error removing stock reservation: %v", err)
		return err
	}

	_ = order.SetState(model.PAYED)
	err = s.repository.UpdateOrder(ctx, order)
	if err != nil {
		log.Printf("[order_service] Error updating order state: %v", err)
		return err
	}

	return nil
}

func (s *Service) OrderCancel(ctx context.Context, orderID int64) error {
	order, err := s.repository.GetById(ctx, orderID)
	if err != nil {
		log.Printf("[order_service] Error getting order: %v", err)
		return err
	}

	if order.State() != model.AWAITING_PAYMENT {
		log.Printf("[order_service] Invalid order state: %v", order.State())
		return appErr.ErrOrderState
	}

	err = s.stockService.ReserveCancel(ctx, order.Items)
	if err != nil {
		log.Printf("[order_service] Error cancelling stock reservation: %v", err)
		return err
	}

	_ = order.SetState(model.CANCELLED)
	err = s.repository.UpdateOrder(ctx, order)
	if err != nil {
		log.Printf("[order_service] Error updating order state: %v", err)
		return err
	}

	return nil
}

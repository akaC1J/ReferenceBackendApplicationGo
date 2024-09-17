package stockservice

import (
	"context"
	"fmt"
	"log"
	appErr "route256/loms/internal/errors"
	"route256/loms/internal/model"
	"route256/loms/internal/service/orderservice"
)

var _ orderservice.StockService = (*Service)(nil)

type Repository interface {
	GetStock(ctx context.Context, sku model.SKUType) (model.Stock, error)
	UpdateStock(ctx context.Context, items []model.Stock) error
}

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) Reserve(ctx context.Context, items []*model.Item) error {
	itemMap := aggregateItem(items)

	var updateStocks []model.Stock

	for sku, totalCount := range itemMap {
		stock, err := s.repository.GetStock(ctx, sku)
		if err != nil {
			log.Printf("[stock_service] Error getting stock for SKU %v: %v", sku, err)
			return err
		}
		availableCount := stock.TotalCount - stock.ReservedCount
		if availableCount < totalCount {
			log.Printf("[stock_service] Not enough stock for SKU %v: requested %v, available %v", sku, totalCount, availableCount)
			return fmt.Errorf("not enough stock for SKU %v: %w", sku, appErr.ErrStockInsufficient)
		}

		stock.ReservedCount += totalCount
		updateStocks = append(updateStocks, stock)
	}

	err := s.repository.UpdateStock(ctx, updateStocks)
	if err != nil {
		log.Printf("[stock_service] Error updating stock: %v", err)
		return err
	}

	return nil
}

func (s *Service) ReserveRemove(ctx context.Context, items []*model.Item) error {
	itemMap := aggregateItem(items)

	var updateStocks []model.Stock

	for sku, totalCount := range itemMap {
		stock, err := s.repository.GetStock(ctx, sku)
		if err != nil {
			log.Printf("[stock_service] Error getting stock for SKU %v: %v", sku, err)
			return err
		}
		if stock.ReservedCount < totalCount {
			log.Printf("[stock_service] Reserved count less than requested for SKU %v: requested %v, reserved %v", sku, totalCount, stock.ReservedCount)
			return appErr.ErrNegativeReserved
		}
		if stock.TotalCount < totalCount {
			log.Printf("[stock_service] Total count less than requested for SKU %v: requested %v, total %v", sku, totalCount, stock.TotalCount)
			return appErr.ErrStockInsufficient
		}

		stock.ReservedCount -= totalCount
		stock.TotalCount -= totalCount

		updateStocks = append(updateStocks, stock)
	}
	err := s.repository.UpdateStock(ctx, updateStocks)
	if err != nil {
		log.Printf("[stock_service] Error updating stock: %v", err)
		return err
	}

	return nil
}

func (s *Service) ReserveCancel(ctx context.Context, items []*model.Item) error {
	itemMap := aggregateItem(items)

	var updateStocks []model.Stock

	for sku, totalCount := range itemMap {
		stock, err := s.repository.GetStock(ctx, sku)
		if err != nil {
			log.Printf("[stock_service] Error getting stock for SKU %v: %v", sku, err)
			return err
		}
		if stock.ReservedCount < totalCount {
			log.Printf("[stock_service] Reserved count less than requested for SKU %v: requested %v, reserved %v", sku, totalCount, stock.ReservedCount)
			return appErr.ErrNegativeReserved
		}

		stock.ReservedCount -= totalCount

		updateStocks = append(updateStocks, stock)
	}
	err := s.repository.UpdateStock(ctx, updateStocks)
	if err != nil {
		log.Printf("[stock_service] Error updating stock: %v", err)
		return err
	}
	return nil
}

func (s *Service) GetBySKUAvailableCount(ctx context.Context, sku model.SKUType) (uint64, error) {
	stock, err := s.repository.GetStock(ctx, sku)
	if err != nil {
		log.Printf("[stock_service] Error getting stock for SKU %v: %v", sku, err)
		return 0, err
	}
	if stock.TotalCount < stock.ReservedCount {
		log.Printf("[stock_service] Total count less than reserved for SKU %v: total %v, reserved %v", sku, stock.TotalCount, stock.ReservedCount)
		return 0, appErr.ErrNegativeAvailable
	}
	availableCount := stock.TotalCount - stock.ReservedCount

	return uint64(availableCount), nil
}

func aggregateItem(items []*model.Item) map[model.SKUType]uint32 {
	itemMap := make(map[model.SKUType]uint32)
	for _, item := range items {
		itemMap[item.SKU] += item.Count
	}
	return itemMap
}

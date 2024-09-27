package stockservice

import (
	"context"
	"fmt"
	"log"
	appErr "route256/loms/internal/errors"
	"route256/loms/internal/model"
	"route256/loms/internal/repository/stockrepository"
	"route256/loms/internal/service/orderservice"
)

var _ orderservice.StockService = (*Service)(nil)
var _ Repository = (*stockrepository.Repository)(nil)

type Repository interface {
	GetStocks(ctx context.Context, sku []model.SKUType) ([]*model.Stock, error)
	UpdateStock(ctx context.Context, items map[model.SKUType]*model.Stock) error
}

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) Reserve(ctx context.Context, items []*model.Item) error {
	return s.processItems(ctx, items, func(stock *model.Stock, neededCount uint32) error {
		availableCount := stock.TotalCount - stock.ReservedCount
		if availableCount < neededCount {
			return fmt.Errorf("not enough stock for SKU %v: %w", stock.SKU, appErr.ErrStockInsufficient)
		}
		stock.ReservedCount += neededCount
		return nil
	})
}

func (s *Service) ReserveRemove(ctx context.Context, items []*model.Item) error {
	return s.processItems(ctx, items, func(stock *model.Stock, neededCount uint32) error {
		if stock.ReservedCount < neededCount {
			return appErr.ErrNegativeReserved
		}
		if stock.TotalCount < neededCount {
			return appErr.ErrStockInsufficient
		}
		stock.ReservedCount -= neededCount
		stock.TotalCount -= neededCount
		return nil
	})
}

func (s *Service) ReserveCancel(ctx context.Context, items []*model.Item) error {
	return s.processItems(ctx, items, func(stock *model.Stock, neededCount uint32) error {
		if stock.ReservedCount < neededCount {
			return appErr.ErrNegativeReserved
		}
		stock.ReservedCount -= neededCount
		return nil
	})
}

func (s *Service) GetBySKUAvailableCount(ctx context.Context, sku model.SKUType) (uint64, error) {
	stocks, err := s.repository.GetStocks(ctx, []model.SKUType{sku})
	if err != nil {
		log.Printf("[stock_service] Error getting stock for SKU %v: %v", sku, err)
		return 0, err
	}
	stock := stocks[0]
	if stock.TotalCount < stock.ReservedCount {
		return 0, appErr.ErrNegativeAvailable
	}
	return uint64(stock.TotalCount - stock.ReservedCount), nil
}

func (s *Service) processItems(ctx context.Context, items []*model.Item, processFunc func(*model.Stock, uint32) error) error {
	itemMap := makeSkuCountMap(items)
	skus := getSKUList(itemMap)

	stocks, err := s.repository.GetStocks(ctx, skus)
	if err != nil {
		log.Printf("[stock_service] Error getting stocks: %v", err)
		return err
	}

	updateStocks := make(map[model.SKUType]*model.Stock)
	for _, stock := range stocks {
		neededCount := itemMap[stock.SKU]
		err = processFunc(stock, neededCount)
		if err != nil {
			log.Printf("[stock_service] Error processing SKU %v: %v", stock.SKU, err)
			return err
		}
		updateStocks[stock.SKU] = stock
	}

	return s.repository.UpdateStock(ctx, updateStocks)
}

func getSKUList(itemMap map[model.SKUType]uint32) []model.SKUType {
	skus := make([]model.SKUType, 0, len(itemMap))
	for sku := range itemMap {
		skus = append(skus, sku)
	}
	return skus
}

func makeSkuCountMap(items []*model.Item) map[model.SKUType]uint32 {
	itemMap := make(map[model.SKUType]uint32)
	for _, item := range items {
		itemMap[item.SKU] += item.Count
	}
	return itemMap
}

package lomsservice

import (
	"context"
	"go.uber.org/multierr"
	"log"
	"route256/cart/internal/generated/api/loms/v1"
	"route256/cart/internal/pkg/model"
	"route256/cart/internal/pkg/service/cartservice"
	"route256/cart/internal/pkg/utils"
)

var _ cartservice.LomsService = (*LomsService)(nil)

type LomsService struct {
	client loms.LomsClient
}

func NewOrderRequest(userId model.UserId, cart map[model.SKU]model.CartItem) (*loms.Order, error) {
	orderInRq := &loms.Order{User: int64(userId)}
	var resultError error
	for sku, item := range cart {
		safeSku, err := utils.SafeInt64ToUint32(int64(sku))
		resultError = multierr.Append(resultError, err)
		orderInRq.Items = append(orderInRq.Items, &loms.Item{
			Sku: safeSku,
			//unsafe cast uint16 to uint32
			Count: uint32(item.Count),
		})
	}
	if resultError != nil {
		return nil, resultError
	}
	return orderInRq, nil
}

func NewLomsService(client loms.LomsClient) *LomsService {
	return &LomsService{client: client}
}

func (s *LomsService) CreateOrder(ctx context.Context, userId model.UserId, cart map[model.SKU]model.CartItem) (orderId int64, err error) {
	orderInRq, err := NewOrderRequest(userId, cart)
	if err != nil {
		log.Printf("[orderservice] Error creating order request: %v", err)
		return 0, err
	}
	orderIdRs, err := s.client.OrderCreate(ctx, &loms.OrderCreateRequest{Order: orderInRq})
	if err != nil {
		log.Printf("[orderservice] Error creating order: %v", err)
		return 0, err
	}
	log.Printf("Success creating order for user %d, order ID %d", userId, orderIdRs.OrderId)
	return orderIdRs.OrderId, err
}

func (s *LomsService) GetStockInfo(ctx context.Context, sku model.SKU) (availableCountStock uint64, err error) {
	safeSku, err := utils.SafeInt64ToUint32(int64(sku))
	if err != nil {
		return 0, err
	}
	availableCount, err := s.client.StocksInfo(ctx, &loms.StocksInfoRequest{Sku: safeSku})
	if err != nil {
		log.Printf("[orderservice] Error getting stock info: %v", err)
		return 0, err
	}
	log.Printf("Success getting stock info for SKU %d, available %d", sku, availableCount.GetCount())
	return availableCount.GetCount(), nil

}

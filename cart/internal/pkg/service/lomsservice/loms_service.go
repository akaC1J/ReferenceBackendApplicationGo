package lomsservice

import (
	"context"
	"log"
	"route256/cart/internal/pkg/model"
	"route256/cart/internal/pkg/service/cartservice"
	"route256/loms/pkg/api/loms/v1"
)

var _ cartservice.LomsService = (*LomsService)(nil)

type LomsService struct {
	client loms.LomsClient
}

func NewLomsService(client loms.LomsClient) *LomsService {
	return &LomsService{client: client}
}

func (s *LomsService) CreateOrder(ctx context.Context, userId model.UserId, cart map[model.SKU]model.CartItem) (int64, error) {
	orderInRq := &loms.Order{User: int64(userId)}
	for sku, item := range cart {
		orderInRq.Items = append(orderInRq.Items, &loms.Item{
			Sku:   uint32(sku),
			Count: uint32(item.Count),
		})
	}

	orderIdRs, err := s.client.OrderCreate(ctx, &loms.OrderCreateRequest{Order: orderInRq})
	if err != nil {
		log.Printf("[orderservice] Error creating order: %v", err)
		return 0, err
	}
	log.Printf("Success creating order for user %d, order ID %d", userId, orderIdRs.OrderId)
	return orderIdRs.OrderId, err
}

func (s *LomsService) GetStockInfo(ctx context.Context, sku model.SKU) (uint64, error) {
	availableCount, err := s.client.StocksInfo(ctx, &loms.StocksInfoRequest{Sku: uint32(sku)}) //два тз противоречат другу, только вот как быть с этим?
	// везде пилить по проверке - это кошмар
	if err != nil {
		log.Printf("[orderservice] Error getting stock info: %v", err)
		return 0, err
	}
	log.Printf("Success getting stock info for SKU %d, available %d", sku, availableCount.GetCount())
	return availableCount.GetCount(), nil

}

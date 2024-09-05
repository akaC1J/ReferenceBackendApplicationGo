package cartservice

import (
	"context"
	"fmt"
	"log"
	"route256/cart/internal/pkg/model"
)

type CartRepository interface {
	InsertItem(context.Context, model.CartItem) (*model.CartItem, error)
	RemoveItem(_ context.Context, userId model.UserId, sku model.SKU) error
	RemoveByUserId(_ context.Context, userId model.UserId) error
	GetItem(_ context.Context, userId model.UserId) (map[model.SKU]model.CartItem, error)
}

type ProductService interface {
	// GetProductInfo получает информацию о продукте по SKU.
	// Возвращает либо валидный продукт, если он существует, либо ошибку в случае,
	// если продукт не найден или данные о нем некорректны.
	GetProductInfo(ctx context.Context, sku model.SKU) (*model.Product, error)
}

type CartService struct {
	repository     CartRepository
	productService ProductService
}

func NewService(repository CartRepository, service ProductService) *CartService {
	return &CartService{repository: repository, productService: service}
}

func (s *CartService) AddCartItem(ctx context.Context, cartItem model.CartItem) (*model.CartItem, error) {
	if errSku := checkFieldMustPositive(int64(cartItem.SKU), "sku"); errSku != nil {
		log.Printf("[cartService] Failed to add item to cart: SKU validation failed for SKU %d", cartItem.SKU)
		return nil, errSku
	}

	if errUserId := checkFieldMustPositive(int64(cartItem.UserId), "user_id"); errUserId != nil {
		log.Printf("[cartService] Failed to add item to cart: UserID validation failed for UserID %d", cartItem.UserId)
		return nil, errUserId
	}

	if errCount := checkFieldMustPositive(int64(cartItem.Count), "count"); errCount != nil {
		log.Printf("[cartService] Failed to add item to cart: Count validation failed for Count %d", cartItem.Count)
		return nil, errCount
	}

	log.Printf("[cartService] Fetching product info for SKU %d", cartItem.SKU)
	_, err := s.productService.GetProductInfo(ctx, cartItem.SKU)
	if err != nil {
		log.Printf("[cartService] Failed to add item to cart: Product info for SKU %d not found", cartItem.SKU)
		return nil, err
	}

	item, err := s.repository.InsertItem(ctx, cartItem)
	if err != nil {
		log.Printf("[cartService] Failed to add item to cart for user %d: %v", cartItem.UserId, err)
		return nil, err
	}

	log.Printf("[cartService] Item with SKU %d successfully added to the cart for user %d, count: %d", cartItem.SKU, cartItem.UserId, cartItem.Count)

	return item, nil
}

func (s *CartService) DeleteCartItem(ctx context.Context, userId model.UserId, sku model.SKU) error {
	// Валидация SKU и UserID
	if errSku := checkFieldMustPositive(int64(sku), "sku"); errSku != nil {
		log.Printf("[cartService] Failed to delete item from cart: SKU validation failed for SKU %d", sku)
		return errSku
	}

	if errUserId := checkFieldMustPositive(int64(userId), "user_id"); errUserId != nil {
		log.Printf("[cartService] Failed to delete item from cart: UserID validation failed for UserID %d", userId)
		return errUserId
	}

	err := s.repository.RemoveItem(ctx, userId, sku)
	if err != nil {
		log.Printf("[cartService] Failed to delete item with SKU %d from cart for user %d: %v", sku, userId, err)
		return err
	}

	// Логируем успешное удаление товара из корзины
	log.Printf("[cartService] Item with SKU %d successfully removed from the cart for user %d", sku, userId)

	return nil
}

func (s *CartService) CleanUpCart(ctx context.Context, userId model.UserId) error {
	if errUserId := checkFieldMustPositive(int64(userId), "user_id"); errUserId != nil {
		log.Printf("[cartService] Failed to clean up cart: UserID validation failed for UserID %d", userId)
		return errUserId
	}

	err := s.repository.RemoveByUserId(ctx, userId)
	if err != nil {
		log.Printf("[cartService] Failed to clean up cart for user %d: %v", userId, err)
		return err
	}

	log.Printf("[cartService] Cart for user %d successfully cleaned up", userId)

	return nil
}

func (s *CartService) GetCartItem(ctx context.Context, userId model.UserId) (*CartContent, error) {
	if errUserId := checkFieldMustPositive(int64(userId), "user_id"); errUserId != nil {
		log.Printf("[cartService] Failed to retrieve cart: validation failed: for UserID %d", userId)
		return nil, errUserId
	}

	userCart, err := s.repository.GetItem(ctx, userId)
	if err != nil {
		log.Printf("[cartService] Failed to retrieve cart for user %d: %v", userId, err)
		return nil, err
	}
	log.Printf("[cartService] Retrieving cart successed %+v", userCart)

	// Обогащение данных о товарах
	var cartContent CartContent
	for keySku, item := range userCart {
		productInfo, err := s.productService.GetProductInfo(ctx, keySku)
		if err != nil {
			log.Printf("[cartService] Failed to retrieve product info for SKU %d while retrieving cart for user %d", keySku, userId)
			return nil, err
		}
		cartContent.Items = append(cartContent.Items, createEnrichedCartItemDTO(item, *productInfo))
		cartContent.TotalPrice += productInfo.Price * uint32(item.Count)
	}

	log.Printf("[cartService] Cart for user %d successfully retrieved with %d items, total price: %d", userId, len(cartContent.Items), cartContent.TotalPrice)

	return &cartContent, nil
}

func checkFieldMustPositive(value int64, fieldName string) error {
	if value < 1 {
		return fmt.Errorf("field " + fieldName + " must be positive")
	}
	return nil
}

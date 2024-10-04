package cartservice

import (
	"context"
	"fmt"
	"log"
	"route256/cart/internal/infra/errgroup"
	"route256/cart/internal/pkg/model"
	"sync"
)

type CartRepository interface {
	InsertItem(context.Context, model.CartItem) (*model.CartItem, error)
	RemoveItem(_ context.Context, userId model.UserId, sku model.SKU) error
	RemoveByUserId(_ context.Context, userId model.UserId) error
	GetCartByUserId(_ context.Context, userId model.UserId) (map[model.SKU]model.CartItem, error)
}

type ProductService interface {
	// GetProductInfo получает информацию о продукте по SKU.
	// Возвращает либо валидный продукт, если он существует, либо ошибку в случае,
	// если продукт не найден или данные о нем некорректны.
	GetProductInfo(ctx context.Context, sku model.SKU) (*model.Product, error)
}

type LomsService interface {
	CreateOrder(ctx context.Context, userId model.UserId, cart map[model.SKU]model.CartItem) (int64, error)
	GetStockInfo(ctx context.Context, sku model.SKU) (availableCountStock uint64, err error)
}

type CartService struct {
	repository     CartRepository
	productService ProductService
	lomsService    LomsService
}

func NewService(repository CartRepository, service ProductService, lomsService LomsService) *CartService {
	return &CartService{repository: repository, productService: service, lomsService: lomsService}
}

func (s *CartService) AddCartItem(ctx context.Context, cartItem model.CartItem) (*model.CartItem, error) {
	if errValidate := cartItem.Validate(); errValidate != nil {
		return nil, fmt.Errorf("errors during cartservice validate %w", errValidate)
	}

	log.Printf("[cartService] Fetching product info for SKU %d", cartItem.SKU)
	_, err := s.productService.GetProductInfo(ctx, cartItem.SKU)
	if err != nil {
		log.Printf("[cartService] Failed to add item to cart: Product info for SKU %d not found", cartItem.SKU)
		return nil, err
	}

	availableCount, err := s.lomsService.GetStockInfo(ctx, cartItem.SKU)
	if err != nil {
		log.Printf("[cartService] Failed get info from stock for SKU %d", cartItem.SKU)
		return nil, err
	}
	//unsafe cast uint16 to uint64
	if uint64(cartItem.Count) > availableCount {
		log.Printf("[cartService] Failed to add item to cart: Not enough items in stock for SKU %d", cartItem.SKU)
		return nil, fmt.Errorf("not enough items in stock for SKU %d", cartItem.SKU)

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

	userCart, err := s.repository.GetCartByUserId(ctx, userId)
	if err != nil {
		log.Printf("[cartService] Failed to retrieve cart for user %d: %v", userId, err)
		return nil, err
	}
	log.Printf("[cartService] Retrieving cart successed %+v", userCart)

	// Обогащение данных о товарах
	var cartContent CartContent
	errGroup, cancelCtx := errgroup.NewErrGroup(ctx)
	var mx sync.Mutex
	for keySku, item := range userCart {
		keySku := keySku
		item := item
		errGroup.Go(func() error {
			productInfo, err := s.productService.GetProductInfo(cancelCtx, keySku)
			if err != nil {
				log.Printf("[cartService] Failed to retrieve product info for SKU %d while retrieving cart for user %d", keySku, userId)
				return err
			}
			mx.Lock()
			defer mx.Unlock()
			cartContent.Items = append(cartContent.Items, createEnrichedCartItemDTO(item, *productInfo))
			cartContent.TotalPrice += productInfo.Price * uint32(item.Count)
			return nil
		})
	}

	if err := errGroup.Wait(); err != nil {
		return nil, err
	}
	log.Printf("[cartService] Cart for user %d successfully retrieved with %d items, total price: %d", userId, len(cartContent.Items), cartContent.TotalPrice)

	return &cartContent, nil
}

func (s *CartService) Checkout(ctx context.Context, userId model.UserId) (orderId int64, err error) {
	if errUserId := checkFieldMustPositive(int64(userId), "user_id"); errUserId != nil {
		log.Printf("[cartService] Failed to retrieve cart: validation failed: for UserID %d", userId)
		return 0, errUserId
	}
	userCart, err := s.repository.GetCartByUserId(ctx, userId)
	if err != nil {
		log.Printf("[cartService] Failed to retrieve cart for user %d: %v", userId, err)
		return 0, err
	}
	log.Printf("[cartService] Retrieving cart successed %+v", userCart)
	orderId, err = s.lomsService.CreateOrder(ctx, userId, userCart)
	if err != nil {
		log.Printf("[cartService] Failed to create order for user %d: %v", userId, err)
		return 0, err
	}
	err = s.repository.RemoveByUserId(ctx, userId)
	if err != nil {
		log.Printf("[cartService] Failed to clean up cart for user %d: %v", userId, err)
		return 0, err

	}
	return orderId, nil
}

func checkFieldMustPositive(value int64, fieldName string) error {
	if value < 1 {
		return fmt.Errorf("field " + fieldName + " must be positive")
	}
	return nil
}

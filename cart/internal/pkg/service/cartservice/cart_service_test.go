package cartservice

import (
	"context"
	"errors"
	"fmt"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
	"route256/cart/internal/pkg/model"
	"testing"
)

func TestCartService_AddCartItem(t *testing.T) {
	mc := minimock.NewController(t)

	ctx := context.Background()

	cartItem := model.CartItem{
		UserId: 1,
		SKU:    123,
		Count:  1,
	}

	cartItemWithIvalidData := model.CartItem{
		UserId: 0,
		SKU:    123,
		Count:  1,
	}

	product := &model.Product{
		Name:  "TestProduct",
		Price: 100,
	}

	type testStruct struct {
		name           string
		mockRepo       func() *CartRepositoryMock
		mockProductSvc func() *ProductServiceMock
		mockLomsSvc    func() *LomsServiceMock
		args           struct {
			ctx      context.Context
			cartItem model.CartItem
		}
		want       *model.CartItem
		wantErr    bool
		checkMocks func(*testing.T, *CartRepositoryMock, *ProductServiceMock, *LomsServiceMock)
	}

	tests := []testStruct{
		{
			name: "success - add item to cart",
			mockRepo: func() *CartRepositoryMock {
				repoMock := NewCartRepositoryMock(mc)
				repoMock.InsertItemMock.Return(&cartItem, nil)
				return repoMock
			},
			mockProductSvc: func() *ProductServiceMock {
				productServiceMock := NewProductServiceMock(mc)
				productServiceMock.GetProductInfoMock.Return(product, nil)
				return productServiceMock
			},
			mockLomsSvc: func() *LomsServiceMock {
				lomsServiceMock := NewLomsServiceMock(mc)
				var enoughCount uint64 = 99999
				lomsServiceMock.GetStockInfoMock.Return(enoughCount, nil)
				return lomsServiceMock
			},
			args: struct {
				ctx      context.Context
				cartItem model.CartItem
			}{
				ctx:      ctx,
				cartItem: cartItem,
			},
			want:    &cartItem,
			wantErr: false,
			checkMocks: func(t *testing.T, repoMock *CartRepositoryMock, productServiceMock *ProductServiceMock, lomsServiceMock *LomsServiceMock) {
				assert.Equal(t, 1, len(repoMock.InsertItemMock.Calls()))
				assert.Equal(t, 1, len(productServiceMock.GetProductInfoMock.Calls()))
				assert.Equal(t, 1, len(lomsServiceMock.GetStockInfoMock.Calls()))
			},
		},

		{
			name: "error - product not found",
			mockRepo: func() *CartRepositoryMock {
				repoMock := NewCartRepositoryMock(mc)
				return repoMock
			},
			mockProductSvc: func() *ProductServiceMock {
				productServiceMock := NewProductServiceMock(mc)
				productServiceMock.GetProductInfoMock.Expect(ctx, cartItem.SKU).Return(nil, fmt.Errorf("product not found for SKU %d", cartItem.SKU))
				return productServiceMock
			},
			mockLomsSvc: func() *LomsServiceMock {
				lomsServiceMock := NewLomsServiceMock(mc)
				return lomsServiceMock
			},
			args: struct {
				ctx      context.Context
				cartItem model.CartItem
			}{
				ctx:      ctx,
				cartItem: cartItem,
			},
			want:    nil,
			wantErr: true,
			checkMocks: func(t *testing.T, repoMock *CartRepositoryMock, productServiceMock *ProductServiceMock, lomsServiceMock *LomsServiceMock) {
				assert.Equal(t, 0, len(repoMock.InsertItemMock.Calls()))
				assert.Equal(t, 1, len(productServiceMock.GetProductInfoMock.Calls()))
				assert.Equal(t, 0, len(lomsServiceMock.CreateOrderMock.Calls()))
			},
		},
		{
			name: "error - validate item",
			mockRepo: func() *CartRepositoryMock {
				repoMock := NewCartRepositoryMock(mc)
				return repoMock
			},
			mockProductSvc: func() *ProductServiceMock {
				productServiceMock := NewProductServiceMock(mc)
				return productServiceMock
			},
			mockLomsSvc: func() *LomsServiceMock {
				lomsServiceMock := NewLomsServiceMock(mc)
				return lomsServiceMock
			},
			args: struct {
				ctx      context.Context
				cartItem model.CartItem
			}{
				ctx:      ctx,
				cartItem: cartItemWithIvalidData,
			},
			want:    nil,
			wantErr: true,
			checkMocks: func(t *testing.T, repoMock *CartRepositoryMock, productServiceMock *ProductServiceMock, lomsServiceMock *LomsServiceMock) {
				assert.Equal(t, 0, len(repoMock.InsertItemMock.Calls()))
				assert.Equal(t, 0, len(productServiceMock.GetProductInfoMock.Calls()))
				assert.Equal(t, 0, len(lomsServiceMock.CreateOrderMock.Calls()))
			},
		},
		{
			name: "error - repository insert fails",
			mockRepo: func() *CartRepositoryMock {
				repoMock := NewCartRepositoryMock(mc)
				repoMock.InsertItemMock.Return(nil, errors.New("insert error"))
				return repoMock
			},
			mockProductSvc: func() *ProductServiceMock {
				productServiceMock := NewProductServiceMock(mc)
				productServiceMock.GetProductInfoMock.Return(product, nil)
				return productServiceMock
			},
			mockLomsSvc: func() *LomsServiceMock {
				lomsServiceMock := NewLomsServiceMock(mc)
				var enoughCount uint64 = 99999
				lomsServiceMock.GetStockInfoMock.Return(enoughCount, nil)
				return lomsServiceMock
			},
			args: struct {
				ctx      context.Context
				cartItem model.CartItem
			}{
				ctx:      ctx,
				cartItem: cartItem,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error - not enough available stock",
			mockRepo: func() *CartRepositoryMock {
				repoMock := NewCartRepositoryMock(mc)
				return repoMock
			},
			mockProductSvc: func() *ProductServiceMock {
				productServiceMock := NewProductServiceMock(mc)
				productServiceMock.GetProductInfoMock.Return(product, nil)
				return productServiceMock
			},
			mockLomsSvc: func() *LomsServiceMock {
				lomsServiceMock := NewLomsServiceMock(mc)
				notEnoughCount := uint64(cartItem.Count - 1) // Недостаточно товара
				lomsServiceMock.GetStockInfoMock.Return(notEnoughCount, nil)
				return lomsServiceMock
			},
			args: struct {
				ctx      context.Context
				cartItem model.CartItem
			}{
				ctx:      ctx,
				cartItem: cartItem,
			},
			want:    nil,
			wantErr: true,
			checkMocks: func(t *testing.T, repoMock *CartRepositoryMock, productServiceMock *ProductServiceMock, lomsServiceMock *LomsServiceMock) {
				assert.Equal(t, 0, len(repoMock.InsertItemMock.Calls()))
				assert.Equal(t, 1, len(productServiceMock.GetProductInfoMock.Calls()))
				assert.Equal(t, 1, len(lomsServiceMock.GetStockInfoMock.Calls()))
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			//безопасный параллельный запуск, тесты не используют общий контекст или данные
			t.Parallel()
			repoMock := tt.mockRepo()
			productServiceMock := tt.mockProductSvc()
			lomsServiceMock := tt.mockLomsSvc()
			s := NewService(repoMock, productServiceMock, lomsServiceMock)
			got, err := s.AddCartItem(tt.args.ctx, tt.args.cartItem)

			if tt.checkMocks != nil {
				tt.checkMocks(t, repoMock, productServiceMock, lomsServiceMock)
			}
			require.Equal(t, tt.want, got, "AddCartItem() got unexpected result")

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
func TestCartService_DeleteCartItem(t *testing.T) {
	mc := minimock.NewController(t)

	ctx := context.Background()

	type testStruct struct {
		name     string
		mockRepo func() *CartRepositoryMock
		args     struct {
			ctx    context.Context
			userId model.UserId
			sku    model.SKU
		}
		wantErr    bool
		checkMocks func(t *testing.T, repoMock *CartRepositoryMock)
	}

	tests := []testStruct{
		{
			name: "success - delete item from cart",
			mockRepo: func() *CartRepositoryMock {
				repoMock := NewCartRepositoryMock(mc)
				repoMock.RemoveItemMock.Return(nil)
				return repoMock
			},
			args: struct {
				ctx    context.Context
				userId model.UserId
				sku    model.SKU
			}{
				ctx:    ctx,
				userId: 1,
				sku:    123,
			},
			wantErr: false,
			checkMocks: func(t *testing.T, repoMock *CartRepositoryMock) {

				assert.Equal(t, 1, len(repoMock.RemoveItemMock.Calls()))
			},
		},
		{
			name: "error - invalid SKU",
			mockRepo: func() *CartRepositoryMock {
				repoMock := NewCartRepositoryMock(mc)
				return repoMock
			},
			args: struct {
				ctx    context.Context
				userId model.UserId
				sku    model.SKU
			}{
				ctx:    ctx,
				userId: 1,
				sku:    0, // Некорректный SKU
			},
			wantErr: true,
			checkMocks: func(t *testing.T, repoMock *CartRepositoryMock) {
				assert.Equal(t, 0, len(repoMock.RemoveItemMock.Calls()))
			},
		},
		{
			name: "error - invalid UserID",
			mockRepo: func() *CartRepositoryMock {
				repoMock := NewCartRepositoryMock(mc)
				return repoMock
			},
			args: struct {
				ctx    context.Context
				userId model.UserId
				sku    model.SKU
			}{
				ctx:    ctx,
				userId: 0, // Некорректный UserID
				sku:    123,
			},
			wantErr: true,
			checkMocks: func(t *testing.T, repoMock *CartRepositoryMock) {
				assert.Equal(t, 0, len(repoMock.RemoveItemMock.Calls()))
			},
		},
		{
			name: "error - repository failure",
			mockRepo: func() *CartRepositoryMock {
				repoMock := NewCartRepositoryMock(mc)
				// Ошибка при удалении из репозитория
				repoMock.RemoveItemMock.Return(fmt.Errorf("repository error"))
				return repoMock
			},
			args: struct {
				ctx    context.Context
				userId model.UserId
				sku    model.SKU
			}{
				ctx:    ctx,
				userId: 1,
				sku:    123,
			},
			wantErr: true,
			checkMocks: func(t *testing.T, repoMock *CartRepositoryMock) {
				assert.Equal(t, 1, len(repoMock.RemoveItemMock.Calls()))
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			//безопасный параллельный запуск, тесты не используют общий контекст или данные
			t.Parallel()
			repoMock := tt.mockRepo()
			s := NewService(repoMock, nil, nil)
			err := s.DeleteCartItem(tt.args.ctx, tt.args.userId, tt.args.sku)
			if tt.checkMocks != nil {
				tt.checkMocks(t, repoMock)
			}
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCartService_CleanUpCart(t *testing.T) {
	mc := minimock.NewController(t)

	ctx := context.Background()

	type testStruct struct {
		name     string
		mockRepo func() *CartRepositoryMock
		args     struct {
			ctx    context.Context
			userId model.UserId
		}
		wantErr    bool
		checkMocks func(t *testing.T, repoMock *CartRepositoryMock)
	}

	tests := []testStruct{
		{
			name: "success - clean up cart",
			mockRepo: func() *CartRepositoryMock {
				repoMock := NewCartRepositoryMock(mc)
				repoMock.RemoveByUserIdMock.Return(nil)
				return repoMock
			},
			args: struct {
				ctx    context.Context
				userId model.UserId
			}{
				ctx:    ctx,
				userId: 1,
			},
			wantErr: false,
			checkMocks: func(t *testing.T, repoMock *CartRepositoryMock) {
				assert.Equal(t, 1, len(repoMock.RemoveByUserIdMock.Calls()))
			},
		},
		{
			name: "error - invalid UserID",
			mockRepo: func() *CartRepositoryMock {
				repoMock := NewCartRepositoryMock(mc)
				return repoMock
			},
			args: struct {
				ctx    context.Context
				userId model.UserId
			}{
				ctx:    ctx,
				userId: 0, // Некорректный UserID
			},
			wantErr: true,
			checkMocks: func(t *testing.T, repoMock *CartRepositoryMock) {
				assert.Equal(t, 0, len(repoMock.RemoveByUserIdMock.Calls()))
			},
		},
		{
			name: "error - repository failure",
			mockRepo: func() *CartRepositoryMock {
				repoMock := NewCartRepositoryMock(mc)
				repoMock.RemoveByUserIdMock.Return(fmt.Errorf("repository error"))
				return repoMock
			},
			args: struct {
				ctx    context.Context
				userId model.UserId
			}{
				ctx:    ctx,
				userId: 1,
			},
			wantErr: true,
			checkMocks: func(t *testing.T, repoMock *CartRepositoryMock) {
				assert.Equal(t, 1, len(repoMock.RemoveByUserIdMock.Calls()))
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		//безопасный параллельный запуск, тесты не используют общий контекст или данные
		t.Run(tt.name, func(t *testing.T) {
			repoMock := tt.mockRepo()
			t.Parallel()

			s := NewService(repoMock, nil, nil)

			err := s.CleanUpCart(tt.args.ctx, tt.args.userId)

			if tt.checkMocks != nil {
				tt.checkMocks(t, repoMock)
			}

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCartService_GetCartItem(t *testing.T) {
	mc := minimock.NewController(t)

	ctx := context.Background()

	type testStruct struct {
		name           string
		mockRepo       func() *CartRepositoryMock
		mockProductSvc func() *ProductServiceMock
		args           struct {
			ctx    context.Context
			userId model.UserId
		}
		wantContent *CartContent
		wantErr     bool
		checkMocks  func(t *testing.T, repoMock *CartRepositoryMock, productServiceMock *ProductServiceMock)
	}

	tests := []testStruct{
		{
			name: "success - retrieve cart",
			mockRepo: func() *CartRepositoryMock {
				repoMock := NewCartRepositoryMock(mc)
				repoMock.GetCartByUserIdMock.Return(map[model.SKU]model.CartItem{
					123: {UserId: 1, SKU: 123, Count: 2},
				}, nil)
				return repoMock
			},
			mockProductSvc: func() *ProductServiceMock {
				productServiceMock := NewProductServiceMock(mc)
				productServiceMock.GetProductInfoMock.Return(&model.Product{Name: "Any test name", Price: 100}, nil)
				return productServiceMock
			},
			args: struct {
				ctx    context.Context
				userId model.UserId
			}{
				ctx:    ctx,
				userId: 1,
			},
			wantContent: &CartContent{
				Items: []EnrichedCartItem{
					{SKU: 123, Count: 2, Price: 100, Name: "Any test name"},
				},
				TotalPrice: 200,
			},
			wantErr: false,
			checkMocks: func(t *testing.T, repoMock *CartRepositoryMock, productServiceMock *ProductServiceMock) {
				// Проверяем, что GetCartByUserId был вызван 1 раз
				assert.Equal(t, 1, len(repoMock.GetCartByUserIdMock.Calls()))
				// Проверяем, что GetProductInfo был вызван 1 раз
				assert.Equal(t, 1, len(productServiceMock.GetProductInfoMock.Calls()))
			},
		},
		{
			name: "error - invalid UserID",
			mockRepo: func() *CartRepositoryMock {
				return NewCartRepositoryMock(mc)
			},
			mockProductSvc: func() *ProductServiceMock {
				return NewProductServiceMock(mc)
			},
			args: struct {
				ctx    context.Context
				userId model.UserId
			}{
				ctx:    ctx,
				userId: 0, // Некорректный UserID
			},
			wantContent: nil,
			wantErr:     true,
			checkMocks: func(t *testing.T, repoMock *CartRepositoryMock, productServiceMock *ProductServiceMock) {
				// Проверяем, что GetCartByUserId не был вызван
				assert.Equal(t, 0, len(repoMock.GetCartByUserIdMock.Calls()))
				// Проверяем, что GetProductInfo не был вызван
				assert.Equal(t, 0, len(productServiceMock.GetProductInfoMock.Calls()))
			},
		},
		{
			name: "error - repository failure",
			mockRepo: func() *CartRepositoryMock {
				repoMock := NewCartRepositoryMock(mc)
				// Ошибка при получении корзины
				repoMock.GetCartByUserIdMock.Return(nil, fmt.Errorf("repository error"))
				return repoMock
			},
			mockProductSvc: func() *ProductServiceMock {
				return NewProductServiceMock(mc)
			},
			args: struct {
				ctx    context.Context
				userId model.UserId
			}{
				ctx:    ctx,
				userId: 1,
			},
			wantContent: nil,
			wantErr:     true,
			checkMocks: func(t *testing.T, repoMock *CartRepositoryMock, productServiceMock *ProductServiceMock) {
				// Проверяем, что GetCartByUserId был вызван 1 раз
				assert.Equal(t, 1, len(repoMock.GetCartByUserIdMock.Calls()))
				// Проверяем, что GetProductInfo не был вызван
				assert.Equal(t, 0, len(productServiceMock.GetProductInfoMock.Calls()))
			},
		},
		{
			name: "error - product service failure",
			mockRepo: func() *CartRepositoryMock {
				repoMock := NewCartRepositoryMock(mc)
				repoMock.GetCartByUserIdMock.Return(map[model.SKU]model.CartItem{
					123: {UserId: 1, SKU: 123, Count: 2},
				}, nil)
				return repoMock
			},
			mockProductSvc: func() *ProductServiceMock {
				productServiceMock := NewProductServiceMock(mc)
				productServiceMock.GetProductInfoMock.Return(nil, fmt.Errorf("product service error"))
				return productServiceMock
			},
			args: struct {
				ctx    context.Context
				userId model.UserId
			}{
				ctx:    ctx,
				userId: 1,
			},
			wantContent: nil,
			wantErr:     true,
			checkMocks: func(t *testing.T, repoMock *CartRepositoryMock, productServiceMock *ProductServiceMock) {
				assert.Equal(t, 1, len(repoMock.GetCartByUserIdMock.Calls()))
				assert.Equal(t, 1, len(productServiceMock.GetProductInfoMock.Calls()))
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			//безопасный параллельный запуск, тесты не используют общий контекст или данные
			t.Parallel()

			// Проверка на утечку горутин, исключая текущую горутину, так как эксперименты показали, что она может отображаться как утечка.
			//Это связано с тем, что эта горутина, вероятно, управляет параллельными тестами.
			//Проверено: реальная утечка горутин корректно обнаруживается кодом ниже.
			//Для проверки можно добавить в тестируемый метод намеренно утекшую горутину и убедиться, что она будет выявлена.
			defer goleak.VerifyNone(t, goleak.IgnoreCurrent())
			repoMock := tt.mockRepo()
			productServiceMock := tt.mockProductSvc()

			s := NewService(repoMock, productServiceMock, nil)

			got, err := s.GetCartItem(tt.args.ctx, tt.args.userId)

			if tt.checkMocks != nil {
				tt.checkMocks(t, repoMock, productServiceMock)
			}

			require.Equal(t, tt.wantContent, got, "GetCartItem() got unexpected result")

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantContent, got)
			}
		})
	}
}

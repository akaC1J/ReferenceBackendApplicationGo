package repository

import (
	"context"
	"route256/cart/internal/pkg/apperrors"
	"route256/cart/internal/pkg/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	name     string
	mockRepo func() *AbstractStorageMock
	args     struct {
		ctx    context.Context
		userId model.UserId
		sku    model.SKU
		item   model.CartItem
	}
	wantErr bool
}

func TestRepository_RemoveItem(t *testing.T) {
	tests := []testStruct{
		{
			"Remove item successfully",
			func() *AbstractStorageMock {
				storageMock := NewAbstractStorageMock(t)
				storageMock.RemoveItemMock.Expect(model.UserId(1), model.SKU(101)).Return(nil)
				return storageMock
			},
			struct {
				ctx    context.Context
				userId model.UserId
				sku    model.SKU
				item   model.CartItem
			}{
				ctx:    context.Background(),
				userId: model.UserId(1),
				sku:    model.SKU(101),
			},
			false,
		},
		{
			"Remove item with error",
			func() *AbstractStorageMock {
				storageMock := NewAbstractStorageMock(t)
				storageMock.RemoveItemMock.Expect(model.UserId(1), model.SKU(101)).Return(apperrors.ErrCartNotFound)
				return storageMock
			},
			struct {
				ctx    context.Context
				userId model.UserId
				sku    model.SKU
				item   model.CartItem
			}{
				ctx:    context.Background(),
				userId: model.UserId(1),
				sku:    model.SKU(101),
			},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoMock := tt.mockRepo()
			repo := NewRepository(repoMock)

			err := repo.RemoveItem(tt.args.ctx, tt.args.userId, tt.args.sku)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

		})
	}
}

func TestRepository_GetItem(t *testing.T) {
	tests := []testStruct{
		{
			"Get item successfully",
			func() *AbstractStorageMock {
				storageMock := NewAbstractStorageMock(t)
				storageMock.GetCartMock.Expect(model.UserId(1)).Return(map[model.SKU]model.CartItem{
					model.SKU(101): {UserId: 1, SKU: 101, Count: 2},
				}, nil)
				return storageMock
			},
			struct {
				ctx    context.Context
				userId model.UserId
				sku    model.SKU
				item   model.CartItem
			}{
				ctx:    context.Background(),
				userId: model.UserId(1),
			},
			false,
		},
		{
			"Get item - user not found",
			func() *AbstractStorageMock {
				storageMock := NewAbstractStorageMock(t)
				storageMock.GetCartMock.Expect(model.UserId(1)).Return(nil, apperrors.ErrCartNotFound)
				return storageMock
			},
			struct {
				ctx    context.Context
				userId model.UserId
				sku    model.SKU
				item   model.CartItem
			}{
				ctx:    context.Background(),
				userId: model.UserId(1),
			},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoMock := tt.mockRepo()
			repo := NewRepository(repoMock)

			_, err := repo.GetItem(tt.args.ctx, tt.args.userId)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

		})
	}
}

func TestRepository_RemoveByUserId(t *testing.T) {
	tests := []testStruct{
		{
			"Remove user by ID successfully",
			func() *AbstractStorageMock {
				storageMock := NewAbstractStorageMock(t)
				storageMock.RemoveByUserIdMock.Expect(model.UserId(1)).Return(nil)
				return storageMock
			},
			struct {
				ctx    context.Context
				userId model.UserId
				sku    model.SKU
				item   model.CartItem
			}{
				ctx:    context.Background(),
				userId: model.UserId(1),
			},
			false,
		},
		{
			"Remove user by ID - user not found",
			func() *AbstractStorageMock {
				storageMock := NewAbstractStorageMock(t)
				storageMock.RemoveByUserIdMock.Expect(model.UserId(999)).Return(apperrors.ErrUserNotFound)
				return storageMock
			},
			struct {
				ctx    context.Context
				userId model.UserId
				sku    model.SKU
				item   model.CartItem
			}{
				ctx:    context.Background(),
				userId: model.UserId(999),
			},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoMock := tt.mockRepo()
			repo := NewRepository(repoMock)

			err := repo.RemoveByUserId(tt.args.ctx, tt.args.userId)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

		})
	}
}

func TestRepository_InsertItem(t *testing.T) {
	tests := []testStruct{
		{
			"Insert item successfully",
			func() *AbstractStorageMock {
				storageMock := NewAbstractStorageMock(t)
				storageMock.AddItemMock.Expect(model.UserId(1), model.CartItem{
					UserId: 1,
					SKU:    101,
					Count:  2,
				}).Return()
				return storageMock
			},
			struct {
				ctx    context.Context
				userId model.UserId
				sku    model.SKU
				item   model.CartItem
			}{
				ctx: context.Background(),
				item: model.CartItem{
					UserId: 1,
					SKU:    101,
					Count:  2,
				},
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoMock := tt.mockRepo()
			repo := NewRepository(repoMock)

			result, err := repo.InsertItem(tt.args.ctx, tt.args.item)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, &tt.args.item, result)
			}

		})
	}
}

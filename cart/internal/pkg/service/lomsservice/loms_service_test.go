package lomsservice

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"route256/cart/internal/generated/api/loms/v1"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"route256/cart/internal/pkg/model"
)

type LomsServiceSuite struct {
	suite.Suite
	mc             *minimock.Controller
	mockLomsClient *LomsClientMock
	lomsSvc        *LomsService
}

func TestLomsServiceSuite(t *testing.T) {
	suite.Run(t, new(LomsServiceSuite))
}

func (suite *LomsServiceSuite) SetupTest() {
	suite.mc = minimock.NewController(suite.T())
	suite.mockLomsClient = NewLomsClientMock(suite.mc)
	suite.lomsSvc = NewLomsService(suite.mockLomsClient)
}

func (suite *LomsServiceSuite) TestCreateOrder_Success() {
	userId := model.UserId(123)
	cart := map[model.SKU]model.CartItem{
		1001: {Count: 2},
		1002: {Count: 3},
	}

	orderRequest := &loms.OrderCreateRequest{
		Order: &loms.Order{
			User: int64(userId),
			Items: []*loms.Item{
				{Sku: 1001, Count: 2},
				{Sku: 1002, Count: 3},
			},
		},
	}

	orderResponse := &loms.OrderCreateResponse{
		OrderId: 1,
	}

	suite.mockLomsClient.OrderCreateMock.
		Set(func(ctx context.Context,
			in *loms.OrderCreateRequest,
			opts ...grpc.CallOption) (op1 *loms.OrderCreateResponse, err error) {
			assert.ElementsMatch(suite.T(), orderRequest.Order.Items, in.Order.Items)
			assert.Equal(suite.T(), orderRequest.Order.User, in.Order.User)
			return orderResponse, nil
		})

	orderId, err := suite.lomsSvc.CreateOrder(context.Background(), userId, cart)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(1), orderId)
}

func (suite *LomsServiceSuite) TestCreateOrder_Error() {
	userId := model.UserId(123)
	cart := map[model.SKU]model.CartItem{
		1001: {Count: 2},
	}

	orderRequest := &loms.OrderCreateRequest{
		Order: &loms.Order{
			User: int64(userId),
			Items: []*loms.Item{
				{Sku: 1001, Count: 2},
			},
		},
	}

	createErr := errors.New("failed to create order")

	suite.mockLomsClient.OrderCreateMock.
		Set(func(ctx context.Context,
			in *loms.OrderCreateRequest,
			opts ...grpc.CallOption) (op1 *loms.OrderCreateResponse, err error) {
			assert.ElementsMatch(suite.T(), orderRequest.Order.Items, in.Order.Items)
			assert.Equal(suite.T(), orderRequest.Order.User, in.Order.User)
			return nil, createErr
		})

	orderId, err := suite.lomsSvc.CreateOrder(context.Background(), userId, cart)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), int64(0), orderId)
	assert.Equal(suite.T(), createErr, err)
}

func (suite *LomsServiceSuite) TestGetStockInfo_Success() {
	sku := model.SKU(1001)
	expectedCount := uint64(150)

	stocksInfoRequest := &loms.StocksInfoRequest{
		Sku: 1001,
	}

	stocksInfoResponse := &loms.StocksInfoResponse{
		Count: expectedCount,
	}

	suite.mockLomsClient.StocksInfoMock.
		Expect(context.Background(), stocksInfoRequest).
		Return(stocksInfoResponse, nil)

	availableCount, err := suite.lomsSvc.GetStockInfo(context.Background(), sku)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedCount, availableCount)
}

func (suite *LomsServiceSuite) TestGetStockInfo_NotFound() {
	sku := model.SKU(9999)

	stocksInfoRequest := &loms.StocksInfoRequest{
		Sku: 9999,
	}

	notFoundErr := status.Error(codes.NotFound, "stock not found")

	suite.mockLomsClient.StocksInfoMock.
		Expect(context.Background(), stocksInfoRequest).
		Return(nil, notFoundErr)

	availableCount, err := suite.lomsSvc.GetStockInfo(context.Background(), sku)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), uint64(0), availableCount)
	assert.Equal(suite.T(), notFoundErr, err)
}

func (suite *LomsServiceSuite) TestGetStockInfo_Error() {
	sku := model.SKU(1001)

	stocksInfoRequest := &loms.StocksInfoRequest{
		Sku: 1001,
	}

	internalErr := errors.New("internal server error")

	suite.mockLomsClient.StocksInfoMock.
		Expect(context.Background(), stocksInfoRequest).
		Return(nil, internalErr)

	availableCount, err := suite.lomsSvc.GetStockInfo(context.Background(), sku)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), uint64(0), availableCount)
	assert.Equal(suite.T(), internalErr, err)
}

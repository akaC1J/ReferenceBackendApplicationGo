package productservice

import (
	"context"
	"encoding/json"
	"route256/cart/internal/infra/cache/redis"
	"route256/cart/internal/pkg/model"
	"strconv"
)

type SKUWrapper model.SKU

func (s SKUWrapper) String() string {
	return strconv.Itoa(int(s))
}

type ProductWrapper struct {
	value *model.Product
}

func (p *ProductWrapper) UnmarshalJSON(bytes []byte) error {
	if p.value == nil {
		p.value = new(model.Product)
	}

	return json.Unmarshal(bytes, p.value)
}

func (p *ProductWrapper) MarshalBinary() (data []byte, err error) {
	return json.Marshal(p.value)
}

type ProductCacherAdapterImpl struct {
	underlying *redis.Cacher[SKUWrapper, *ProductWrapper]
}

func NewProductCacherAdapterImpl(underlying *redis.Cacher[SKUWrapper, *ProductWrapper]) *ProductCacherAdapterImpl {
	return &ProductCacherAdapterImpl{underlying: underlying}
}

func (p *ProductCacherAdapterImpl) Put(ctx context.Context, key model.SKU, value *model.Product) {
	wrappedKey := SKUWrapper(key)
	wrappedValue := &ProductWrapper{value: value}
	p.underlying.Put(ctx, wrappedKey, wrappedValue)
}

func (p *ProductCacherAdapterImpl) Get(ctx context.Context, key model.SKU) (*model.Product, error) {
	wrappedKey := SKUWrapper(key)
	wrappedValue, err := p.underlying.Get(ctx, wrappedKey)

	if err != nil {
		return nil, err
	}
	return wrappedValue.value, nil
}

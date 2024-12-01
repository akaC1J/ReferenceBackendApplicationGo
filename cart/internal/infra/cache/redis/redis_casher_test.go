package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"math/rand"
	"reflect"
	"route256/cart/internal/logger"
	"strings"
	"testing"
)

// Типы для тестирования
type LargeExample struct {
	ID      int                 `json:"id"`
	Name    string              `json:"name"`
	Details map[string]string   `json:"details"`
	Items   []map[string]string `json:"items"`
}

func (e LargeExample) MarshalBinary() (data []byte, err error) {
	return json.Marshal(e)
}

type StringWrapper string

func (s StringWrapper) String() string {
	return string(s)
}

// Генерация большого JSON
func generateLargeJSON(size int) LargeExample {
	details := map[string]string{}
	for i := 0; i < size; i++ {
		details[fmt.Sprintf("key_%d", i)] = strings.Repeat("value", 10)
	}

	items := make([]map[string]string, size)
	for i := range items {
		items[i] = map[string]string{
			"item_key":   fmt.Sprintf("item_%d", i),
			"item_value": strings.Repeat("item_value", 10),
		}
	}

	return LargeExample{
		ID:      rand.Intn(1000),
		Name:    strings.Repeat("large_example_name", 10),
		Details: details,
		Items:   items,
	}
}

// Реализация с рефлексией
func (c *Cacher[K, V]) GetWithReflection(ctx context.Context, key K) (V, error) {
	data, err := c.client.Get(ctx, key.String()).Bytes()
	if err != nil {
		var zero V
		return zero, err
	}

	var res V
	resEl := reflect.New(reflect.TypeOf(res).Elem()).Interface()
	res = resEl.(V)

	err = json.Unmarshal(data, res)
	if err != nil {
		var zero V
		return zero, err
	}

	return res, nil
}

func (c *Cacher[K, V]) GetWithCheckPointerReflection(ctx context.Context, key K) (V, error) {
	data, err := c.client.Get(ctx, key.String()).Bytes()
	if err != nil {
		var zero V
		return zero, err
	}

	var res V

	// Проверяем, является ли V указателем
	if reflect.TypeOf(res).Kind() == reflect.Ptr {
		ptr := reflect.New(reflect.TypeOf(res).Elem()).Interface()
		err = json.Unmarshal(data, ptr)
		if err != nil {
			var zero V
			logger.Errorw(ctx, "Error unmarshal", "key", key, "error", err)
			return zero, err
		}
		return ptr.(V), nil
	} else {
		err = json.Unmarshal(data, &res)
		if err != nil {
			var zero V
			logger.Errorw(ctx, "Error unmarshal", "key", key, "error", err)
			return zero, err
		}
		return res, nil
	}
}

// Реализация без рефлексии
func (c *Cacher[K, V]) GetWithoutReflection(ctx context.Context, key K) ([]byte, error) {
	data, err := c.client.Get(ctx, key.String()).Bytes()
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Бенчмарк с большим JSON
func BenchmarkLargeJSON(b *testing.B) {
	client := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	cacher := Cacher[StringWrapper, *LargeExample]{client: client}
	ctx := context.Background()

	largeJSON := generateLargeJSON(1000)
	data, _ := json.Marshal(largeJSON)
	_ = cacher.client.Set(ctx, "large_example", data, 0).Err()

	b.Run("WithReflection", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = cacher.GetWithReflection(ctx, "large_example")
		}
	})

	b.Run("WithCheckPointerReflection", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = cacher.GetWithCheckPointerReflection(ctx, "large_example")
		}
	})

	b.Run("WithoutReflection", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			data, _ := cacher.GetWithoutReflection(ctx, StringWrapper("large_example"))
			var result LargeExample
			_ = json.Unmarshal(data, &result)
		}
	})
}

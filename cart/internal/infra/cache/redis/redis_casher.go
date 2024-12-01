package redis

import (
	"context"
	"encoding"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"reflect"
	"route256/cart/internal/logger"
	"time"
)

var NotAPointerType = errors.New("V must be a pointer to a struct")

// Cacher is a generic Redis caching utility.
//
// Requirements:
// - `K` must implement `fmt.Stringer` (used for converting keys to strings).
// - `V` must implement `encoding.BinaryMarshaler` (used for marshaling values).
// - `V` **must be a pointer to a struct** to allow proper unmarshaling with `json.Unmarshal`.
//
// Incorrect usage (causes runtime error):
//
//	type User struct { ID int }
//	cacher := NewRedisCacher[string, User](redisClient, time.Minute) // Fails
//
// Correct usage:
//
//	type User struct { ID int }
//	cacher := NewRedisCacher[string, *User](redisClient, time.Minute)
type Cacher[K fmt.Stringer, V encoding.BinaryMarshaler] struct {
	client *redis.Client
	ttl    time.Duration
}

// NewRedisCacher creates a new instance of the Cacher with the provided Redis client and TTL.
//
// Requirements:
// - `K` must implement `fmt.Stringer` (used for converting keys to strings).
// - `V` must implement `encoding.BinaryMarshaler` (used for marshaling values).
// - `V` **must be a pointer to a struct** for proper unmarshaling in `Get`.
//
// Incorrect usage (causes runtime error):
//
//	type User struct { ID int }
//	cacher := NewRedisCacher[string, User](redisClient, time.Minute) // Fails
//
// Correct usage:
//
//	type User struct { ID int }
//	cacher := NewRedisCacher[string, *User](redisClient, time.Minute)
func NewRedisCacher[K fmt.Stringer, V encoding.BinaryMarshaler](client *redis.Client, ttl time.Duration) *Cacher[K, V] {
	return &Cacher[K, V]{client: client, ttl: ttl}
}

func (c *Cacher[K, V]) Put(ctx context.Context, key K, value V) {
	c.client.Set(ctx, key.String(), value, c.ttl)
}

func (c *Cacher[K, V]) Get(ctx context.Context, key K) (V, error) {
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
		return res, NotAPointerType
	}
}

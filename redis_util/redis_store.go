package redis_util

import "context"

type RedisStore interface {
	Set(ctx context.Context, key string, value interface{}, expiryInMinutes int) error
	SetInSeconds(ctx context.Context, key string, value interface{}, expiryInSeconds int) error
	SetNX(ctx context.Context, key string, value interface{}, expiryInMinutes int) (bool, error)
	Get(ctx context.Context, key string, dest interface{}) error
	Delete(ctx context.Context, key string) error
	DeleteAll(ctx context.Context, pattern string) error
}

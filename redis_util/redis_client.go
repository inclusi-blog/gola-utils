package redis_util

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/gola-glitch/gola-utils/constants"
	"github.com/gola-glitch/gola-utils/logging"
	"go.opencensus.io/trace"
	"time"
)

//TODO: Remove standalone client once cluster has been deployed in all envs
type redisStore struct {
	rdb *redis.Client
}

func NewRedisClient(host string, port string, db int, readTimeout int, dialTimeout int, writeTimeout int) (RedisStore, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:         host + ":" + port,
		DB:           db,
		DialTimeout:  time.Duration(dialTimeout) * time.Second,
		ReadTimeout:  time.Duration(readTimeout) * time.Second,
		WriteTimeout: time.Duration(writeTimeout) * time.Second,
	})
	return redisStore{rdb: rdb}, rdb.Ping().Err()
}

func NewRedisClientWith(config RedisStoreConfig) (RedisStore, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:         config.Host + ":" + config.Port,
		DB:           config.Db,
		DialTimeout:  time.Duration(config.DialTimeoutInSeconds) * time.Second,
		ReadTimeout:  time.Duration(config.ReadTimeoutInSeconds) * time.Second,
		WriteTimeout: time.Duration(config.WriteTimeoutInSeconds) * time.Second,
	})
	return redisStore{rdb: rdb}, rdb.Ping().Err()
}

func (rs redisStore) Set(ctx context.Context, key string, value interface{}, expiryInMinutes int) error {
	initTracing(ctx, "Set")

	p, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return rs.rdb.Set(key, p, time.Duration(expiryInMinutes)*time.Minute).Err()
}

func (rs redisStore) SetInSeconds(ctx context.Context, key string, value interface{}, expiryInSeconds int) error {
	initTracing(ctx, "SetInSeconds")

	p, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return rs.rdb.Set(key, p, time.Duration(expiryInSeconds)*time.Second).Err()
}

func (rs redisStore) SetNX(ctx context.Context, key string, value interface{}, expiryInMinutes int) (bool, error) {
	initTracing(ctx, "SetNX")

	p, err := json.Marshal(value)
	if err != nil {
		return false, err
	}
	cmd := rs.rdb.SetNX(key, p, time.Duration(expiryInMinutes)*time.Minute)
	return cmd.Val(), cmd.Err()
}

func (rs redisStore) Get(ctx context.Context, key string, dest interface{}) error {
	initTracing(ctx, "Get")

	p, err := rs.rdb.Get(key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(p, dest)
}

func (rs redisStore) Delete(ctx context.Context, key string) error {
	initTracing(ctx, "Del")
	return rs.rdb.Del(key).Err()
}

func (rs redisStore) DeleteAll(ctx context.Context, pattern string) error {
	initTracing(ctx, "DeleteAll")

	cursor := uint64(constants.REDIS_SCAN_DEAFULT_CURSOR_VALUE)
	return rs.scanAndDelete(cursor, pattern, constants.REDIS_SCAN_DEFAULT_COUNT)
}

func (rs redisStore) scanAndDelete(cursor uint64, pattern string, count int64) error {
	logger := logging.NewLoggerEntry()

	keys, updatedCursor, err := rs.rdb.Scan(cursor, pattern, count).Result()
	if err != nil {
		logger.Errorf("RedisStore.scanAndDelete: Error in scanning keys. Error: %v", err)
		return err
	}

	keysLength := len(keys)
	logger.Infof("RedisStore.scanAndDelete: %v keys found", keysLength)
	if keysLength > 0 {
		pipeline := rs.rdb.Pipeline()
		for _, key := range keys {
			pipeline.Del(key)
		}

		_, err = pipeline.Exec()
		if err != nil {
			logger.Errorf("RedisStore.scanAndDelete: Error in pipeline. Error: %v", err)
			return err
		}
	}
	if updatedCursor == 0 {
		logger.Debugf("RedisStore.scanAndDelete: Zero value returned for cursor. Stopping execution")
		return nil
	}
	return rs.scanAndDelete(updatedCursor, pattern, count)
}

func initTracing(ctx context.Context, cmdName string) {
	if ginContext, ok := ctx.(*gin.Context); ok {
		ctx = ginContext.Request.Context()
	}
	_, span := trace.StartSpan(ctx, "redis.(*baseClient)."+cmdName)

	defer span.End()
}

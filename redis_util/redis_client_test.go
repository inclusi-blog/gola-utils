package redis_util

import (
	"context"
	"encoding/json"
	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v7"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type RedisStoreTestSuite struct {
	suite.Suite
	mockCtrl    *gomock.Controller
	redisClient RedisStore
	mockRedis   *miniredis.Miniredis
	ctx         context.Context
}

func (suite *RedisStoreTestSuite) SetupTest() {
	var err error
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.mockRedis, err = miniredis.Run()
	if err != nil {
		panic(err)
	}
	suite.ctx = context.Background()
	suite.redisClient, err = NewRedisClient(suite.mockRedis.Host(), suite.mockRedis.Port(), 0, 10, 10, 10, "")
	suite.Nil(err)
}

func (suite *RedisStoreTestSuite) TearDownTest() {
	suite.mockCtrl.Finish()
	suite.mockRedis.Close()
}

func TestRedisStoreTestSuite(t *testing.T) {
	suite.Run(t, new(RedisStoreTestSuite))
}

func (suite RedisStoreTestSuite) TestShouldSetObjectToRedis() {
	dataToBeCached := "gola_redis_value"
	key := "gola_redis_key"

	err := suite.redisClient.Set(suite.ctx, key, dataToBeCached, 10)
	suite.Nil(err)

	var actualData string
	dataFromCache, err := suite.mockRedis.Get(key)
	suite.Nil(err)
	err = json.Unmarshal([]byte(dataFromCache), &actualData)

	suite.Nil(err)
	suite.Equal(dataToBeCached, actualData)
}

func (suite RedisStoreTestSuite) TestShouldSetObjectToRedisInSeconds() {
	dataToBeCached := "gola_redis_value"
	key := "gola_redis_key"

	err := suite.redisClient.SetInSeconds(suite.ctx, key, dataToBeCached, 20)
	suite.Nil(err)

	var actualData string
	dataFromCache, err := suite.mockRedis.Get(key)
	suite.Nil(err)
	err = json.Unmarshal([]byte(dataFromCache), &actualData)

	suite.Nil(err)
	suite.Equal(dataToBeCached, actualData)
}

func (suite RedisStoreTestSuite) TestShouldGetObjectFromRedis() {
	dataToBeCached := "gola_redis_value"
	key := "gola_redis_key"
	marshalledData, _ := json.Marshal(dataToBeCached)

	err := suite.mockRedis.Set(key, string(marshalledData))
	suite.Nil(err)
	var dataFromCache string

	err = suite.redisClient.Get(suite.ctx, key, &dataFromCache)
	suite.Nil(err)
	suite.Equal(dataToBeCached, dataFromCache)
}

func (suite RedisStoreTestSuite) TestGetShouldReturnErrorWhenKeyIsNotPresentInRedis() {
	key := "key-not-present-in-redis"

	var userJourney string
	err := suite.redisClient.Get(suite.ctx, key, &userJourney)
	suite.Equal(redis.Nil, err)
}

func (suite RedisStoreTestSuite) TestShouldSetAndGetStringToRedis() {
	dataToBeCached := "gola_redis_value"
	key := "gola_redis_key"

	err := suite.redisClient.Set(suite.ctx, key, dataToBeCached, 10)
	suite.Nil(err)

	var actualData string
	err = suite.redisClient.Get(suite.ctx, key, &actualData)
	suite.Nil(err)

	suite.Equal(dataToBeCached, actualData)
}

func (suite RedisStoreTestSuite) TestShouldThrowErrorWhenTypeIsMismatchedInRedisFetch() {
	type Frm struct {
		Id string `json:"id"`
	}

	dataToBeCached := "gola_redis_value"
	key := "gola_redis_key"

	err := suite.redisClient.Set(suite.ctx, key, dataToBeCached, 10)
	suite.Nil(err)

	actualData := Frm{}
	err = suite.redisClient.Get(suite.ctx, key, &actualData)
	suite.Error(err)
}

func (suite RedisStoreTestSuite) TestShouldThrowErrorWhenRedisConnectionIsUnavailable() {
	dataToBeCached := "gola_redis_value"
	key := "gola_redis_key"

	suite.mockRedis.Close()

	err := suite.redisClient.Set(suite.ctx, key, dataToBeCached, 10)
	suite.Error(err)
}

func (suite RedisStoreTestSuite) TestShouldSetAndDeleteKeyFromRedis() {
	dataToBeCached := "gola_redis_value"
	key := "gola_redis_key"
	marshalledData, _ := json.Marshal(dataToBeCached)

	err := suite.mockRedis.Set(key, string(marshalledData))
	suite.Nil(err)

	err = suite.redisClient.Delete(suite.ctx, key)
	suite.Nil(err)
}

func (suite RedisStoreTestSuite) TestShouldNotReplaceExistingKeyWithNewValueWhenSetNXIsUsed() {
	key := "key"
	value := "value"
	differentValue := "different_value"

	val, err := suite.redisClient.SetNX(suite.ctx, key, value, 5)
	suite.True(val)
	suite.Nil(err)

	val, err = suite.redisClient.SetNX(suite.ctx, key, differentValue, 5)
	suite.False(val)
	suite.Nil(err)

	var actualValue string
	valueInCache, err := suite.mockRedis.Get(key)
	suite.Nil(err)

	err = json.Unmarshal([]byte(valueInCache), &actualValue)
	suite.Nil(err)
	suite.Equal(value, actualValue)
}

func (suite RedisStoreTestSuite) TestShouldDeleteAllKeysMatchingAPattern() {
	suite.addKeysToRedis("gola_1", "value 1")
	suite.addKeysToRedis("gola_2", "value 2")
	suite.addKeysToRedis("gola 3", "value 3")
	suite.addKeysToRedis("gola 4", "value 4")

	err := suite.redisClient.DeleteAll(suite.ctx, "gola*")

	value, _ := suite.mockRedis.Get("gola_1")
	suite.Nil(err)
	suite.Equal("", value)
}

func (suite RedisStoreTestSuite) TestShouldReturnNilErrorWhenNoKeysAreFound() {
	err := suite.redisClient.DeleteAll(suite.ctx, "gola*")

	suite.Nil(err)
}

func (suite RedisStoreTestSuite) addKeysToRedis(key, value string) {
	marshalledData, _ := json.Marshal(value)

	err := suite.mockRedis.Set(key, string(marshalledData))
	suite.Nil(err)
}

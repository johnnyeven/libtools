package redis_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/johnnyeven/libtools/redis"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestRedis_GetSet(t *testing.T) {
	tt := assert.New(t)

	r := redis.Redis{
		Host: "test.g7pay.chinawayltd.com",
		Port: 6379,
	}
	r.MarshalDefaults(&r)

	r.Init()

	err := r.GetCache().Set("test", []byte("test"), 7200*time.Second)
	tt.Nil(err)

	bytes, err := r.GetCache().Get("test")
	tt.Nil(err)

	tt.Equal("test", string(bytes))
}

type SetSuite struct {
	suite.Suite
	r        redis.Redis
	cacheKey string
}

func (suite *SetSuite) SetupTest() {
	suite.r = redis.Redis{
		Host: "test.g7pay.chinawayltd.com",
		Port: 6379,
	}
	suite.r.MarshalDefaults(&suite.r)
	suite.r.Init()

	rand.Seed(time.Now().UnixNano())
	suite.cacheKey = fmt.Sprintf("%d", rand.Int63())

	_, err := suite.r.GetCache().HSetWithExpire(suite.cacheKey, time.Second*3)
	suite.Nil(err)

	time.Sleep(time.Second)
	exist, err := suite.r.GetCache().Exists(suite.cacheKey)
	suite.Nil(err)
	suite.True(exist)
}

func (suite *SetSuite) TestHGetUint64() {
	var field uint64 = 12
	_, err := suite.r.GetCache().HSet(suite.cacheKey, fmt.Sprintf("%d", field), field)
	suite.Nil(err)
	value, err := suite.r.GetCache().HGetUint64(suite.cacheKey, fmt.Sprintf("%d", field))
	suite.Nil(err)
	suite.Equal(field, value)
}

func (suite *SetSuite) TestHGetInt64() {
	var field int64 = 12
	_, err := suite.r.GetCache().HSet(suite.cacheKey, fmt.Sprintf("%d", field), field)
	suite.Nil(err)
	value, err := suite.r.GetCache().HGetInt64(suite.cacheKey, fmt.Sprintf("%d", field))
	suite.Nil(err)
	suite.Equal(field, value)
}

func (suite *SetSuite) TestHGetString() {
	var field string = "12"
	_, err := suite.r.GetCache().HSet(suite.cacheKey, field, field)
	suite.Nil(err)
	value, err := suite.r.GetCache().HGetString(suite.cacheKey, field)
	suite.Nil(err)
	suite.Equal(field, value)
}

func TestSetSuite(t *testing.T) {
	suite.Run(t, new(SetSuite))
}

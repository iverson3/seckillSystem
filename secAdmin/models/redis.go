package models

import (
	"github.com/garyburd/redigo/redis"
)

var (
	RedisPool *redis.Pool
	RedisConf RedisConfig
)

type RedisConfig struct {
	Addr        string
	PassWd      string
	MaxIdle     int
	MaxActive   int
	IdleTimeout int
	ProductLeftKey string
}

type RedisModel struct {

}

func NewRedisModel() *RedisModel {
	return &RedisModel{}
}

func SetRedis(pool *redis.Pool, redisConf RedisConfig) {
	RedisPool = pool
	RedisConf = redisConf
}

// 从redis中获取指定活动商品的剩余数量
func (this *RedisModel) GetProductLeftNum(activityId int) (num int, err error) {
	conn := RedisPool.Get()
	defer conn.Close()

	num, err = redis.Int(conn.Do("hget", RedisConf.ProductLeftKey, activityId))
	return
}
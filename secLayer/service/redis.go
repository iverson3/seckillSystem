package service

import (
	"github.com/astaxie/beego/logs"
	"github.com/garyburd/redigo/redis"
	"time"
)

func initRedis() (err error) {
	secLayerContext.Proxy2LayerRedisPool, err = initAndGetRedisPool(secLayerContext.SecLayerConfig.Proxy2LayerRedis)
	if err != nil {
		logs.Error("init Proxy2LayerRedisPool failed! error: %v", err)
		return
	}

	secLayerContext.Layer2ProxyRedisPool, err = initAndGetRedisPool(secLayerContext.SecLayerConfig.Layer2ProxyRedis)
	if err != nil {
		logs.Error("init Layer2ProxyRedisPool failed! error: %v", err)
		return
	}

	return
}

func initAndGetRedisPool(redisConf RedisConf) (pool *redis.Pool, err error) {
	pool = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", redisConf.RedisAddr, redis.DialPassword(redisConf.RedisPassword))
		},
		MaxIdle:         redisConf.RedisMaxIdle,
		MaxActive:       redisConf.RedisMaxActive,
		IdleTimeout:     time.Duration(redisConf.RedisIdleTimeout) * time.Second,
	}

	conn := pool.Get()
	defer conn.Close()
	_, err = conn.Do("ping")
	if err != nil {
		logs.Error("ping redis failed! error: %v", err)
		return
	}
	return
}
package service

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"github.com/garyburd/redigo/redis"
)

func writeToRedis() {
	for {
		conn := SeckillConfig.Proxy2LayerRedisPool.Get()
		for {
			req := <-SeckillConfig.SecReqChan
			logs.Debug("got request from channel")

			data, err := json.Marshal(req)
			if err != nil {
				logs.Error("json Marshal for request failed! request: %v; error: %v", req, err)
				continue
			}

			_, err = conn.Do("RPUSH", SeckillConfig.Redis.RedisProxy2LayerQueueKey, data)
			if err != nil {
				logs.Error("rpush request to redis failed! request: %v; error: %v", req, err)
				break
			}
			logs.Debug("rpush request to redis success")
		}
		conn.Close()
	}
}

func readFromRedis() {
	for {
		conn := SeckillConfig.Proxy2LayerRedisPool.Get()
		for {
			values, err := redis.Values(conn.Do("BLPOP", SeckillConfig.Redis.RedisLayer2ProxyQueueKey, 0))
			if err != nil {
				logs.Error("lpop user_response from redis failed! error: %v", err)
				break
			}
			// redis中blpop返回的是元组对象，因此需要进行特殊处理
			data := B2S(values[1].([]uint8))
			logs.Debug("got response from redis! response: %s", data)

			var resp SecResponse
			err = json.Unmarshal([]byte(data), &resp)
			if err != nil {
				logs.Debug("json Unmarshal for user_response_str failed! user_response_str: %s; error: %v", data, err)
				continue
			}

			SeckillConfig.ReqMapLock.RLock()
			request, ok := SeckillConfig.SecRequestMap[resp.UserId]
			SeckillConfig.ReqMapLock.RUnlock()
			if !ok {
				continue
			}
			request.ResultChan <- &resp
		}
		conn.Close()
	}
}

// 将[]uint8类型的数据转为string类型
func B2S(bs []uint8) string {
	ba := []byte{}
	for _, b := range bs {
		ba = append(ba, byte(b))
	}
	return string(ba)
}
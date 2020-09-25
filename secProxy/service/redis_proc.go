package service

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"strings"
)

func writeToRedis() {
	retryTimes := 0
	var req *SecRequest
	var data []byte
	var err error
	for {
		conn := SeckillConfig.Proxy2LayerRedisPool.Get()
		for {
			if data == nil {
				req = <-SeckillConfig.SecReqChan
				logs.Debug("got request from channel")

				data, err = json.Marshal(req)
				if err != nil {
					logs.Error("json Marshal for request failed! request: %v; error: %v", req, err)
					continue
				}
			}

			_, err = conn.Do("RPUSH", SeckillConfig.Redis.RedisProxy2LayerQueueKey, string(data))
			if err != nil {
				if strings.Contains(err.Error(), "An existing connection was forcibly closed by the remote host") && retryTimes < 3 {
					retryTimes++
					logs.Error("===================== retry times: %d", retryTimes)
					break
				} else {
					resp := &SecResponse{
						UserId:     req.UserId,
						ActivityId: req.ActivityId,
						Token:      "",
						TokenTime:  0,
						Code:       ErrServiceBusy,
					}
					SendResponseToResultChan(resp)

					logs.Error("rpush request to redis failed! request: %v; error: %v", req, err)
					data = nil
					req = nil
					retryTimes = 0
					break
				}
			}
			data = nil
			req = nil
			retryTimes = 0
			logs.Debug("rpush request to redis success")
		}
		conn.Close()
	}
}

func readFromRedis() {
	for {
		conn := SeckillConfig.Proxy2LayerRedisPool.Get()
		for {
			values, err := conn.Do("BLPOP", SeckillConfig.Redis.RedisLayer2ProxyQueueKey, 0)
			if err != nil {
				logs.Error("lpop user_response from redis failed! error: %v", err)
				break
			}
			// redis中blpop返回的是元组对象，因此需要进行特殊处理
			dataArr, ok := values.([]interface{})
			if !ok || len(dataArr) != 2 {
				logs.Error("blpop response from redis failed！ error: Type Assertion is wrong.")
				continue
			}
			data, ok := dataArr[1].([]byte)
			if !ok {
				logs.Error("blpop response from redis failed！ error: Type Assertion is wrong.")
				continue
			}
			logs.Debug("got response from redis! response: %s", string(data))

			var resp SecResponse
			err = json.Unmarshal(data, &resp)
			if err != nil {
				logs.Debug("json Unmarshal for user_response_str failed! user_response_str: %s; error: %v", data, err)
				continue
			}

			SendResponseToResultChan(&resp)
		}
		conn.Close()
	}
}

func SendResponseToResultChan(resp *SecResponse) {
	SeckillConfig.ReqMapLock.RLock()
	request, ok := SeckillConfig.SecRequestMap[resp.UserId]
	SeckillConfig.ReqMapLock.RUnlock()
	if ok {
		request.ResultChan <- resp
	}
}

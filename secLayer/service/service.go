package service

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"math/rand"
	"time"
)

func Run() {
	for i := 0; i < secLayerContext.SecLayerConfig.ReadLayer2ProxyGoroutineNum; i++ {
		secLayerContext.SecWaitGroup.Add(1)
		go HandleReader()
	}
	for i := 0; i < secLayerContext.SecLayerConfig.WriteProxy2LayerGoroutineNum; i++ {
		secLayerContext.SecWaitGroup.Add(1)
		go HandleWriter()
	}

	for i := 0; i < secLayerContext.SecLayerConfig.HandleUserGoroutineNum; i++ {
		secLayerContext.SecWaitGroup.Add(1)
		go HandleSecUserRequest()
	}

	logs.Debug("all process goroutine have started")
	secLayerContext.SecWaitGroup.Wait()
	logs.Debug("all process goroutine exited")
}

func HandleReader()  {
	for {
		conn := secLayerContext.Proxy2LayerRedisPool.Get()
		for {
			logs.Debug("prepare to read request from redis queue")
			// BLPOP：阻塞式的从redis list队列中获取元素，超时时间设置为0 表示无限超时时间，没有元素则一直阻塞
			values, err := conn.Do("BLPOP", secLayerContext.SecLayerConfig.Proxy2LayerRedis.RedisProxy2LayerQueueKey, 0)
			if err != nil {
				logs.Error("lpop user_request from redis failed! error: %v", err)
				break
			}
			// redis中blpop返回的是元组对象，因此需要进行特殊处理
			dataArr, ok := values.([]interface{})
			if !ok || len(dataArr) != 2 {
				logs.Error("lpop user_request from redis failed！ error: Type Assertion is wrong.")
				continue
			}
			data, ok := dataArr[1].([]byte)
			if !ok {
				logs.Error("lpop user_request from redis failed！ error: Type Assertion is wrong.")
				continue
			}
			logs.Debug("got user_request from redis! request: %s", string(data))

			var req SecRequest
			err = json.Unmarshal(data, &req)
			if err != nil {
				logs.Debug("json Unmarshal user_request_str failed! user_request_str: %s; error: %v", data, err)
				continue
			}

			// 如果客户端请求已经超过了最大的等待时间 则直接将请求丢弃
			if time.Now().Unix() - req.AccessTime.Unix() >= int64(secLayerContext.SecLayerConfig.MaxRequestWaitTimeout) {
				logs.Warn("user request is expire! request: %v", req)
				continue
			}

			timer := time.NewTicker(time.Millisecond * time.Duration(secLayerContext.SecLayerConfig.SendToHandleChanTimeout))
			// 放入channel 让业务处理协程们去从中拿取并进行处理
			select {
			case secLayerContext.Read2HandleChan <- &req:
			case <-timer.C:
				logs.Warn("send request to channel timeout, request: %v", req)
				break
			}


		}
		conn.Close()
	}
}

func HandleWriter()  {
	for {
		conn := secLayerContext.Layer2ProxyRedisPool.Get()
		for {
			resp := <-secLayerContext.Handle2WriteChan
			logs.Debug("got user response from resp_channel, response: %v", resp)
			data, err := json.Marshal(resp)
			if err != nil {
				logs.Error("json Marshal for response failed!  response: %v, error: %v", resp, err)
				continue
			}

			logs.Debug("prepare push user response to redis.")
			_, err = conn.Do("RPUSH", secLayerContext.SecLayerConfig.Layer2ProxyRedis.RedisLayer2ProxyQueueKey, string(data))
			if err != nil {
				logs.Error("rpush response to redis failed! response: %s, error: %v", string(data), err)
				break
			}
			logs.Debug("rpush response to redis ok ok ok ok ok!!!.")
		}
		conn.Close()
	}
}

func HandleSecUserRequest() {
	logs.Debug("start handle user request")

	for req := range secLayerContext.Read2HandleChan {
		logs.Debug("got user request to handle; request: %v", req)

		resp, err := handleSeckill(req)
		if err != nil {
			logs.Warn("handle user seckill request failed! request: %v, error: %v", req, err)
			resp = &SecResponse{
				UserId:     req.UserId,
				ActivityId: req.ActivityId,
				Token:      "",
				Code:       ErrServiceBusy,
			}
		}

		timer := time.NewTicker(time.Millisecond * time.Duration(secLayerContext.SecLayerConfig.SendToWriteChanTimeout))
		// 将处理的结果放入channel
		select {
		case secLayerContext.Handle2WriteChan <- resp:
		case <-timer.C:
			logs.Warn("send response to channel timeout, response: %v", resp)
			break
		}
	}
}

func handleSeckill(req *SecRequest) (res *SecResponse, err error) {
	secLayerContext.SecActivityRwLock.RLock()
	defer secLayerContext.SecActivityRwLock.RUnlock()

	logs.Debug("start to handle user request: %v", req)

	res = &SecResponse{
		UserId:     req.UserId,
		ActivityId: req.ActivityId,
		Code:       ErrSeckillSuccess,
	}

	activity, ok := secLayerContext.SecLayerConfig.SecActivityListMap[req.ActivityId]
	// 找不到商品
	if !ok {
		logs.Error("not found activity: %v", req.ActivityId)
		res.Code = ErrActivityNotFound
		logs.Warn("handle user request end! request: %v, result: %v", req, "ErrActivityNotFound")
		return
	}

	// 活动被禁用或活动已结束
	if activity.Status == ActivityStatusDisable || activity.Status == ActivityStatusExpire {
		res.Code = ErrActivityOver
		logs.Warn("handle user request end! request: %v, result: %v", req, "ErrActivityOver")
		return
	}

	// 商品售罄
	if activity.Status == ActivityStatusSoldOut {
		res.Code = ErrActivityProductSoldOut
		logs.Warn("handle user request end! request: %v, result: %v", req, "ErrActivityProductSoldOut")
		return
	}

	// 限速检测 (每件商品在单位时间内的销售数量都是有相应限制的)
	now := time.Now().Unix()
	if activity.SecSoldLimit.Check(now) >= activity.MaxSoldLimit {
		res.Code = ErrReTry
		logs.Warn("handle user request end! request: %v, result: %v", req, "ErrReTry: SecSoldLimit")
		return
	}

	// 用户购买当前商品的数量限制
	secLayerContext.HistoryMapLock.Lock()
	userHistory, ok := secLayerContext.UserBuyHistoryMap[req.UserId]
	if !ok {
		userHistory = &UserBuyHistory{
			history: make(map[int]int, 16),
		}
		secLayerContext.UserBuyHistoryMap[req.UserId] = userHistory
	}
	secLayerContext.HistoryMapLock.Unlock()

	count := userHistory.GetProductBuyCount(activity.ActivityId)
	if count >= activity.UserMaxBuyLimit {
		res.Code = ErrAlreadyBuy
		logs.Warn("handle user request end! request: %v, result: %v", req, "ErrAlreadyBuy")
		return
	}

	// 秒杀商品已售出的数量到达了该商品的最大库存数
	soldCount := secLayerContext.ProductCountManager.Count(activity.ActivityId)
	if soldCount >= activity.Total {
		res.Code = ErrActivityProductSoldOut
		activity.Status = ActivityStatusSoldOut

		syncActivityChangeToEtcd(activity)
		logs.Warn("handle user request end! request: %v, result: %v", req, "ErrActivityProductSoldOut")
		return
	}

	// 根据当前秒杀商品的成功概率 进行随机秒杀，计算出用户本次请求能否抢到该商品
	rate := rand.Float64()
	if rate > activity.BuyRate {
		res.Code = ErrReTry
		logs.Warn("handle user request end! request: %v, result: %v", req, "ErrReTry: rate failed")
		return
	}

	// 到这里，说明用户此处秒杀成功，抢到当前商品
	userHistory.Increment(activity.ActivityId, 1)
	secLayerContext.ProductCountManager.Increment(activity.ActivityId, 1)

	// 计算秒杀活动剩余的商品数量
	//soldCount = secLayerContext.ProductCountManager.Count(activity.ActivityId)
	//activity.Left = activity.Total - soldCount
	//syncActivityChangeToEtcd(activity)

	// 为用户此次秒杀创建token
	curTime := time.Now().Unix()
	tokenStr := fmt.Sprintf("userid=%d&activityid=%d&timestamp=%d&security=%s",
		req.UserId, activity.ActivityId, curTime, secLayerContext.SecLayerConfig.SeckillTokenPasswd)
	res.Token = fmt.Sprintf("%x", md5.Sum([]byte(tokenStr)))
	res.TokenTime = curTime

	logs.Debug("handle user request is successfully! request: %v", req)

	return
}
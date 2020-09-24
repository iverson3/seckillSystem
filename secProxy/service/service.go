package service

import (
	"crypto/md5"
	"fmt"
	"github.com/astaxie/beego/logs"
	"time"
)

// 获取指定的秒杀活动的状态信息
func SecInfo(activityId int) (data []map[string]interface{}, code int, err error) {
	SeckillConfig.RwLock.RLock()
	defer SeckillConfig.RwLock.RUnlock()

	activity, ok := SeckillConfig.SecActivityListMap[activityId]
	if !ok {
		code = ErrNotFoundActivityId
		err = fmt.Errorf("not found activity_id: %d", activityId)
		return
	}

	data = make([]map[string]interface{}, 0)
	item := getSecInfoByActivityConf(activity)
	data = append(data, item)
	return
}

// 获取秒杀商品列表
func SecInfoList() (data []map[string]interface{}, code int, err error) {
	SeckillConfig.RwLock.RLock()
	defer SeckillConfig.RwLock.RUnlock()

	data = make([]map[string]interface{}, 0)
	for _, activity := range SeckillConfig.SecActivityListMap {
		item := getSecInfoByActivityConf(activity)
		data = append(data, item)
	}
	return
}

// 格式化秒杀商品的状态信息 (便于客户端显示)
func getSecInfoByActivityConf(activity *SecActivityConf) (item map[string]interface{}) {
	start := false
	end := false
	status := "seckill is not start"

	now := time.Now().Unix()
	if now >= activity.EndTime {
		end = true
		status = "seckill is already end"
	} else {
		if now >= activity.StartTime {
			start = true
			status = "seckill is already start"
		}
	}

	// 活动商品售罄了，则活动自动结束
	if activity.Left == 0 || activity.Status == ActivityStatusSoldOut {
		start = false
		end   = true
		status = "product had sale out"
	}
	if activity.Status == ActivityStatusDisable || activity.Status == ActivityStatusExpire {
		start = false
		end   = true
		status = "seckill is already end"
	}

	item = make(map[string]interface{})
	item["activity_id"] = activity.ActivityId
	item["product_id"]  = activity.ProductId
	item["start"]       = start
	item["end"]         = end
	item["status"]      = status
	return
}

func userCheck(req *SecRequest) (err error) {
	found := false
	for _, refer := range SeckillConfig.RefererWhiteList {
		if req.ClientReferer == refer {
			found = true
			break
		}
	}
	if !found {
		logs.Warn("user[%d] is reject by referer, req[%v]", req.UserId, req)
		return fmt.Errorf("invalid request")
	}

	authData := fmt.Sprintf("%d:%s", req.UserId, SeckillConfig.CookieSecretKey)
	authSign := fmt.Sprintf("%x", md5.Sum([]byte(authData)))
	if authSign != req.UserAuthSign {
		err = fmt.Errorf("invalid user cookie auth")
		return
	}
	return
}


// 处理秒杀请求
func SeckillProcess(req *SecRequest) (data map[string]interface{}, code int, err error) {
	SeckillConfig.RwLock.RLock()
	defer SeckillConfig.RwLock.RUnlock()

	//err = userCheck(req)
	//if err != nil {
	//	logs.Warn("userId[%d] invalid, auth check failed, req[%v]", req.UserId, req)
	//	return nil, ErrUserAuthCheckFailed, err
	//}

	//err = antiSpam(req)
	//if err != nil {
	//	logs.Warn("antiSpam error: over user access limit per second, req[%v]", req)
	//	return nil, ErrServiceBusy, err
	//}

	res, code, err := SecInfo(req.ActivityId)
	if err != nil {
		logs.Error("get userId[%d] secInfo failed! req: %v; error: %v", req.UserId, req, err)
		return
	}
	// 秒杀还未开始或已经结束 直接返回
	if res[0]["start"] == false || res[0]["end"] == true {
		logs.Debug("seckill activity is over or isn't started; req: %v", req)
		err = fmt.Errorf("seckill activity is over! reason: %s", res[0]["status"])
		code = ErrRequestSuccess
		return
	}

	logs.Debug("send request to channel; req: %v", req)
	SeckillConfig.SecReqChan <- req

	return
}

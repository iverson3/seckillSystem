package service

import (
	"crypto/md5"
	"fmt"
	"github.com/astaxie/beego/logs"
	"time"
)

// 获取指定的秒杀商品的状态信息
func SecInfo(pid int) (data []map[string]interface{}, code int, err error) {
	SeckillConfig.RwLock.RLock()
	defer SeckillConfig.RwLock.RUnlock()

	product, ok := SeckillConfig.SecProductInfo[pid]
	if !ok {
		code = ErrNotFoundProductId
		err = fmt.Errorf("not found product_id: %d", pid)
		return
	}

	data = make([]map[string]interface{}, 0)
	item := getSecInfoByProductConf(product)
	data = append(data, item)
	return
}

// 获取秒杀商品列表
func SecInfoList() (data []map[string]interface{}, code int, err error) {
	SeckillConfig.RwLock.RLock()
	defer SeckillConfig.RwLock.RUnlock()

	data = make([]map[string]interface{}, 0)
	for id, product := range SeckillConfig.SecProductInfo {
		logs.Info(id)
		item := getSecInfoByProductConf(product)
		data = append(data, item)
	}
	return
}

// 格式化秒杀商品的状态信息 (便于客户端显示)
func getSecInfoByProductConf(product *SecProductInfoConf) (item map[string]interface{}) {
	start := false
	end := false
	status := "seckill is not start"

	now := time.Now().Unix()
	if now >= product.EndTime {
		end = true
		status = "seckill is already end"
	} else {
		if now >= product.StartTime {
			start = true
			status = "seckill is already start"
		}
	}

	if product.Status == ProductStatusForceSoldOut || product.Status == ProductStatusSoldOut {
		start = false
		end   = true
		status = "product had sale out"
	}

	item = make(map[string]interface{})
	item["product_id"] = product.ProductId
	item["start"]      = start
	item["end"]        = end
	item["status"]     = status
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

	res, code, err := SecInfo(req.ProductId)
	if err != nil {
		logs.Error("get userId[%d] secInfo failed! req: %v; error: %v", req.UserId, req, err)
		return
	}
	// 秒杀还未开始或已经结束 直接返回
	if res[0]["start"] == false || res[0]["end"] == true {
		logs.Debug("seckill is over or isn't started; req: %v", req)
		err = fmt.Errorf("seckill is over")
		code = ErrRequestSuccess
		return
	}

	logs.Debug("send request to channel; req: %v", req)
	SeckillConfig.SecReqChan <- req

	return
}

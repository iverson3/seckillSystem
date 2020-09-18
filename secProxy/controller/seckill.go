package controller

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"seckillsystem/secProxy/service"
	"strings"
	"time"
)

type SeckillController struct {
	beego.Controller
}

func (this *SeckillController) SecKill()  {
	res := make(map[string]interface{})
	res["code"] = service.ErrRequestSuccess
	res["message"] = "success"

	defer func() {
		this.Data["json"] = res
		this.ServeJSON()
	}()

	pid, err := this.GetInt("product_id")
	if err != nil {
		res["code"] = service.ErrParamDeletion
		res["message"] = "product_id is null"
		return
	}

	// 处理请求的url参数
	source   := this.GetString("src")      // 来源
	authCode := this.GetString("authcode") // 验证码
	secTime  := this.GetString("time")     // 时间
	nance    := this.GetString("nance")    // 随机数
	if source == "" || authCode == "" || secTime == "" || nance == "" {
		res["code"] = service.ErrParamDeletion
		res["message"] = "param is deletion"
		return
	}

	// 获取cookie信息内容
	//userIdStr := this.Ctx.GetCookie("userid")
	//userid, err := strconv.Atoi(userIdStr)
	//userauthsign := this.Ctx.GetCookie("userauthsign")
	//if err != nil {
	//	res["code"] = service.ErrParamTypeError
	//	res["message"] = fmt.Sprintf("invalid cookie: userid[%s]", userIdStr)
	//	return
	//}
	//if userid == 0 || userauthsign == "" {
	//	res["code"] = service.ErrCookieParamDeletion
	//	res["message"] = "cookie info is deletion"
	//	return
	//}

	// 临时获取userId的方式
	userid, err := this.GetInt("userid")
	if err != nil {
		res["code"] = service.ErrParamDeletion
		res["message"] = "userid is null"
		return
	}
	userauthsign := "xxxxxxxxx"



	// 获取并处理客户端ip地址
	addr := this.Ctx.Request.RemoteAddr
	if len(addr) > 0 && strings.Contains(addr, ":") {
		splitArr := strings.Split(addr, ":")
		addr = splitArr[0]
	} else {
		addr = ""
	}

	referer := this.Ctx.Request.Referer()

	secRequest := &service.SecRequest{
		UserId:        userid,
		UserAuthSign:  userauthsign,
		ProductId:     pid,
		Source:        source,
		AuthCode:      authCode,
		SecTime:       secTime,
		Nance:         nance,
		AccessTime:    time.Now(),
		ClientAddr:    addr,
		ClientReferer: referer,
		ResultChan:    make(chan *service.SecResponse),
	}
	service.SeckillConfig.ReqMapLock.Lock()
	service.SeckillConfig.SecRequestMap[userid] = secRequest
	service.SeckillConfig.ReqMapLock.Unlock()

	_, code, err := service.SeckillProcess(secRequest)
	if err != nil {
		logs.Error("service deal failed! error: %v", err)
		res["code"] = code
		res["message"] = err.Error()
		return
	}

	// 阻塞在这里，等待Layer层返回对用户请求的处理结果
	// Layer层会把处理的结果放入redis中，所以需要开go协程去从redis中读取响应结果数据
	var result *service.SecResponse
	// 用户客户端等待的超时时间，如果在该时间内 无法从redis中读取到请求的处理结果 则直接返回相关错误信息
	ticker := time.NewTicker(time.Second * time.Duration(service.SeckillConfig.MaxRequestWaitTimeout))
	select {
	case result = <-secRequest.ResultChan:
	case <-ticker.C:
	}

	// 在指定的超时时间内 没有从redis中读取到请求的处理结果
	if result == nil {
		service.SeckillConfig.ReqMapLock.Lock()
		delete(service.SeckillConfig.SecRequestMap, userid)
		service.SeckillConfig.ReqMapLock.Unlock()

		logs.Error("got response from response_channel is timeout.............................")
		res["code"] = service.ErrServiceBusy
		res["message"] = "service busy"
		return
	}

	bytes, err := json.Marshal(result)
	if err != nil {
		logs.Error("json Marshal for response failed! response: %v, error: %v", result, err)
		res["code"] = service.ErrServiceBusy
		res["message"] = err.Error()
		return
	}
	res["data"] = string(bytes)
}

func (this *SeckillController) SecInfo() {
	res := make(map[string]interface{})
	res["code"] = service.ErrRequestSuccess
	res["message"] = "success"

	defer func() {
		this.Data["json"] = res
		this.ServeJSON()
	}()

	pid, err := this.GetInt("product_id")
	if err != nil {
		// 获取秒杀商品列表
		list, code, err := service.SecInfoList()
		if err != nil {
			logs.Error("service deal failed! error: %v", err)
			res["code"] = code
			res["message"] = err.Error()
			return
		}
		res["data"] = list
	} else {
		// 获取指定的秒杀商品的状态信息
		data, code, err := service.SecInfo(pid)
		if err != nil {
			logs.Error("service deal failed! error: %v", err)
			res["code"] = code
			res["message"] = err.Error()
			return
		}
		res["data"] = data
	}
}


package service

import (
	"fmt"
	"sync"
)

var (
	secLimitMgr = &SecLimitMgr{
		UserLimitMap: make(map[int]*SecLimit, 1000),
		IpLimitMap: make(map[string]*SecLimit, 1000),
	}
)

type SecLimitMgr struct {
	UserLimitMap map[int]*SecLimit
	IpLimitMap map[string]*SecLimit
	lock sync.Mutex
}

// 检测用户的访问频率
// 判断是否是恶意请求
func antiSpam(req *SecRequest) (err error) {
	secLimitMgr.lock.Lock()
	defer secLimitMgr.lock.Unlock()

	userLimit, ok := secLimitMgr.UserLimitMap[req.UserId]
	if !ok {
		secLimit := &SecLimit{
			count:   1,
			curTime: req.AccessTime.Unix(),
		}
		secLimitMgr.UserLimitMap[req.UserId] = secLimit
		return
	}

	userLimit.Count(req.AccessTime.Unix())
	count := userLimit.Check(req.AccessTime.Unix())

	// 用户每秒请求大于等于指定的次数则判定为恶意请求
	if count >= SeckillConfig.UserAccessLimitPerSecond {
		return fmt.Errorf("invalid request")
	}


	ipLimit, ok := secLimitMgr.IpLimitMap[req.ClientAddr]
	if !ok {
		secLimit := &SecLimit{
			count:   1,
			curTime: req.AccessTime.Unix(),
		}
		secLimitMgr.IpLimitMap[req.ClientAddr] = secLimit
		return
	}

	ipLimit.Count(req.AccessTime.Unix())
	count = ipLimit.Check(req.AccessTime.Unix())

	// 同一个客户端ip每秒请求大于等于指定的次数则判定为恶意请求
	if count >= SeckillConfig.IpAccessLimitPerSecond {
		return fmt.Errorf("invalid request")
	}

	return
}

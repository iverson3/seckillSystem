package service

import (
	"github.com/astaxie/beego/logs"
	"github.com/garyburd/redigo/redis"
	"strconv"
	"time"
)

var (
	SeckillConfig *SeckillConf
)

func InitService(serviceConf *SeckillConf) (err error) {
	SeckillConfig = serviceConf

	err = loadBlackList()
	if err != nil {
		logs.Error("load black list from redis failed")
		return
	}

	err = initProxy2LayerRedis()
	if err != nil {
		logs.Error("init proxy2layer redis failed")
		return
	}

	SeckillConfig.SecReqChan    = make(chan *SecRequest, SeckillConfig.SecReqChanSize)
	SeckillConfig.SecRequestMap = make(map[int]*SecRequest, 100000)
	initRedisProcessFunc()

	logs.Debug("init service success")
	return
}

func initRedisProcessFunc() {
	for i := 0; i < SeckillConfig.WriteProxy2LayerGoroutineNum; i++ {
		go writeToRedis()
	}
	for i := 0; i < SeckillConfig.ReadLayer2ProxyGoroutineNum; i++ {
		go readFromRedis()
	}
}

func initProxy2LayerRedis() (err error) {
	SeckillConfig.Proxy2LayerRedisPool = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", SeckillConfig.Redis.RedisAddr, redis.DialPassword(SeckillConfig.Redis.RedisPassword))
		},
		MaxIdle:     SeckillConfig.Redis.RedisMaxIdle,
		MaxActive:   SeckillConfig.Redis.RedisMaxActive,
		IdleTimeout: time.Duration(SeckillConfig.Redis.RedisIdleTimeout) * time.Second,
	}

	conn := SeckillConfig.Proxy2LayerRedisPool.Get()
	defer conn.Close()
	_, err = conn.Do("ping")
	if err != nil {
		logs.Error("ping redis failed! error: %v", err)
		return
	}
	return
}

// 加载ip黑名单和id黑名单
func loadBlackList() (err error) {
	err = initBlackRedis()
	if err != nil {
		logs.Error("init black list redis failed! error: %v", err)
		return
	}

	conn := SeckillConfig.BlackRedisPool.Get()
	defer conn.Close()

	// 获取user_id黑名单列表
	idList, err := redis.Strings(conn.Do("hgetall", "seckill_id_black_list"))
	if err != nil {
		logs.Error("redis do hgetall(seckill_id_black_list) failed! error: %v", err)
		return
	}
	for _, v := range idList {
		id, err := strconv.Atoi(v)
		if err != nil {
			logs.Warn("invalid user id: %s", v)
			continue
		}
		SeckillConfig.IdBlackList[id] = true
	}

	// 获取客户端ip黑名单列表
	ipList, err := redis.Strings(conn.Do("hgetall", "seckill_ip_black_list"))
	if err != nil {
		logs.Error("redis do hgetall(seckill_ip_black_list) failed! error: %v", err)
		return
	}
	for _, ip := range ipList {
		SeckillConfig.IpBlackList[ip] = true
	}

	go SyncIdBlackList()
	go SyncIpBlackList()

	return
}

// 定时的同步user id黑名单列表
func SyncIdBlackList() {
	conn := SeckillConfig.BlackRedisPool.Get()
	defer conn.Close()

	var idList []int
	lastTime := time.Now().Unix()
	for {
		// 阻塞式的从redis List中拿取一个元素
		id, err := redis.Int(conn.Do("BLPOP", "seckill_id_black_list_inc", time.Second))
		if err != nil {
			//logs.Warn("invalid user id: %s", id)
			continue
		}

		now := time.Now().Unix()
		idList = append(idList, id)
		if len(idList) >= 100 || now - lastTime >= 30 {
			SeckillConfig.BlackRwLock.Lock()
			for _, id := range idList {
				SeckillConfig.IdBlackList[id] = true
			}
			SeckillConfig.BlackRwLock.Unlock()
			logs.Info("sync id black list from redis success, idList: %v", idList)

			// 清空List，重设时间
			idList = idList[0:0]
			lastTime = now
		}
	}
}

// 定时的同步客户端IP黑名单列表
func SyncIpBlackList() {
	conn := SeckillConfig.BlackRedisPool.Get()
	defer conn.Close()

	var ipList []string
	lastTime := time.Now().Unix()
	for {
		// 阻塞式的从redis List中拿取一个元素
		ip, err := redis.String(conn.Do("BLPOP", "seckill_ip_black_list_inc", time.Second))
		if err != nil {
			continue
		}

		now := time.Now().Unix()
		ipList = append(ipList, ip)
		if len(ipList) >= 100 || now - lastTime >= 30 {
			SeckillConfig.BlackRwLock.Lock()
			for _, ip := range ipList {
				SeckillConfig.IpBlackList[ip] = true
			}
			SeckillConfig.BlackRwLock.Unlock()
			logs.Info("sync ip black list from redis success, ipList: %v", ipList)

			// 清空List，重设时间
			ipList = ipList[0:0]
			lastTime = now
		}
	}
}

func initBlackRedis() error {
	SeckillConfig.BlackRedisPool = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", SeckillConfig.BlackListRedis.RedisAddr, redis.DialPassword(SeckillConfig.BlackListRedis.RedisPassword))
		},
		MaxIdle:     SeckillConfig.BlackListRedis.RedisMaxIdle,
		MaxActive:   SeckillConfig.BlackListRedis.RedisMaxActive,
		IdleTimeout: time.Duration(SeckillConfig.BlackListRedis.RedisIdleTimeout) * time.Second,
	}

	conn := SeckillConfig.BlackRedisPool.Get()
	defer conn.Close()
	_, err := conn.Do("ping")
	if err != nil {
		logs.Error("ping redis failed! error: %v", err)
		return err
	}
	return nil
}

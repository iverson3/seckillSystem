package main

import (
	"fmt"
	"github.com/astaxie/beego"
	"strings"
	"seckillsystem/secProxy/service"
)

var (
	seckillConf = &service.SeckillConf{
		SecActivityListMap: make(map[int]*service.SecActivityConf),
	}
)

func initConfig() error {
	redisAddr := beego.AppConfig.String("redis_addr")
	etcdAddr := beego.AppConfig.String("etcd_addr")

	if len(redisAddr) == 0 || len(etcdAddr) == 0 {
		return fmt.Errorf("init config failed! redis[%s] or etcd[%s] config is null", redisAddr, etcdAddr)
	}

	pwd := beego.AppConfig.String("redis_password")
	if len(pwd) == 0 {
		return fmt.Errorf("init config failed! password[%s] config is null", pwd)
	}
	idle, err := beego.AppConfig.Int("redis_max_idle")
	if err != nil {
		return fmt.Errorf("init config failed! redis_max_idle[%d] config is null", idle)
	}
	active, err := beego.AppConfig.Int("redis_max_active")
	if err != nil {
		return fmt.Errorf("init config failed! redis_max_active[%d] config is null", active)
	}
	timeout, err := beego.AppConfig.Int("redis_idle_timeout")
	if err != nil {
		return fmt.Errorf("init config failed! redis_idle_timeout[%d] config is null", timeout)
	}
	proxy2layerQueueKey := beego.AppConfig.String("redis_proxy2layer_queue_key")
	if len(proxy2layerQueueKey) == 0 {
		return fmt.Errorf("init config failed! redis_proxy2layer_queue_key[%s] config is null", proxy2layerQueueKey)
	}
	layer2proxyQueueKey := beego.AppConfig.String("redis_layer2proxy_queue_key")
	if len(layer2proxyQueueKey) == 0 {
		return fmt.Errorf("init config failed! redis_layer2proxy_queue_key[%s] config is null", layer2proxyQueueKey)
	}


	etcdTimeout, err := beego.AppConfig.Int("etcd_timeout")
	if err != nil {
		return fmt.Errorf("init config failed! etcd_timeout[%d] config is null", etcdTimeout)
	}
	etcdSecKeyPrefix  := beego.AppConfig.String("etcd_sec_key_prefix")
	etcdSecActivityKey := beego.AppConfig.String("etcd_sec_activity_key")
	if len(etcdSecKeyPrefix) == 0 {
		return fmt.Errorf("init config failed! etcd_sec_key_prefix[%s] config is null", etcdSecKeyPrefix)
	}
	if len(etcdSecActivityKey) == 0 {
		return fmt.Errorf("init config failed! etcd_sec_activity_key[%s] config is null", etcdSecActivityKey)
	}
	if !strings.HasSuffix(etcdSecKeyPrefix, "/") {
		etcdSecKeyPrefix = etcdSecKeyPrefix + "/"
	}

	logsPath  := beego.AppConfig.String("logs_path")
	logsLevel := beego.AppConfig.String("logs_level")
	if len(logsPath) == 0 || len(logsLevel) == 0 {
		return fmt.Errorf("init config failed! logs_path[%s] or logs_level[%s] config is null", logsPath, logsLevel)
	}

	secretKey := beego.AppConfig.String("cookie_secret_key")
	if len(secretKey) != 32 {
		return fmt.Errorf("init config failed! cookie_secret_key[%s] config is error", secretKey)
	}

	userAccessLimit, err := beego.AppConfig.Int("user_access_limit_per_second")
	if err != nil {
		return fmt.Errorf("init config failed! user_access_limit_per_second[%d] config is null", userAccessLimit)
	}
	ipAccessLimit, err := beego.AppConfig.Int("ip_access_limit_per_second")
	if err != nil {
		return fmt.Errorf("init config failed! ip_access_limit_per_second[%d] config is null", ipAccessLimit)
	}

	var whiteList []string
	refererStr := beego.AppConfig.String("referer_white_list")
	if len(refererStr) > 0 {
		if strings.Contains(refererStr, ",") {
			whiteList = strings.Split(refererStr, ",")
		} else {
			whiteList = append(whiteList, refererStr)
		}
	}

	writeNum, err := beego.AppConfig.Int("write_proxy2layer_goroutine_num")
	if err != nil {
		return fmt.Errorf("init config failed! write_proxy2layer_goroutine_num[%d] config is null", writeNum)
	}
	readNum, err := beego.AppConfig.Int("read_layer2proxy_goroutine_num")
	if err != nil {
		return fmt.Errorf("init config failed! read_layer2proxy_goroutine_num[%d] config is null", readNum)
	}
	requestTimeOut, err := beego.AppConfig.Int("max_request_wait_timeout")
	if err != nil {
		return fmt.Errorf("init config failed! max_request_wait_timeout[%d] config is null", requestTimeOut)
	}

	seckillConf.Redis.RedisAddr = redisAddr
	seckillConf.Redis.RedisPassword = pwd
	seckillConf.Redis.RedisMaxIdle = idle
	seckillConf.Redis.RedisMaxActive = active
	seckillConf.Redis.RedisIdleTimeout = timeout
	seckillConf.Redis.RedisProxy2LayerQueueKey = proxy2layerQueueKey
	seckillConf.Redis.RedisLayer2ProxyQueueKey = layer2proxyQueueKey

	seckillConf.BlackListRedis.RedisAddr = redisAddr
	seckillConf.BlackListRedis.RedisPassword = pwd
	seckillConf.BlackListRedis.RedisMaxIdle = idle
	seckillConf.BlackListRedis.RedisMaxActive = active
	seckillConf.BlackListRedis.RedisIdleTimeout = timeout

	seckillConf.Etcd.EtcdAddr = etcdAddr
	seckillConf.Etcd.EtcdTimeout = etcdTimeout
	seckillConf.Etcd.EtcdSecKeyPrefix = etcdSecKeyPrefix
	seckillConf.Etcd.EtcdSecActivityKey = fmt.Sprintf("%s%s", etcdSecKeyPrefix, etcdSecActivityKey)

	seckillConf.Log.LogPath = logsPath
	seckillConf.Log.LogLevel = logsLevel

	seckillConf.CookieSecretKey = secretKey
	seckillConf.UserAccessLimitPerSecond = userAccessLimit
	seckillConf.IpAccessLimitPerSecond = ipAccessLimit
	seckillConf.RefererWhiteList = whiteList

	seckillConf.WriteProxy2LayerGoroutineNum = writeNum
	seckillConf.ReadLayer2ProxyGoroutineNum = readNum
	seckillConf.MaxRequestWaitTimeout = requestTimeOut
	return nil
}

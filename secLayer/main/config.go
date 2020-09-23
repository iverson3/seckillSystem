package main

import (
	"fmt"
	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/logs"
	"seckillsystem/secLayer/service"
	"strings"
)

var (
	AppConfig *service.SecLayerConf
)

func initConfig(confType, filename string) (err error) {
	conf, err := config.NewConfig(confType, filename)
	if err != nil {
		logs.Error("config.NewConfig failed! error: %v", err)
		return
	}

	AppConfig = &service.SecLayerConf{}
	AppConfig.Log.LogLevel = conf.String("logs::logs_level")
	if len(AppConfig.Log.LogLevel) == 0 {
		AppConfig.Log.LogLevel = "debug"
	}
	AppConfig.Log.LogPath = conf.String("logs::logs_path")
	if len(AppConfig.Log.LogPath) == 0 {
		AppConfig.Log.LogPath = "./logs/default.log"
	}


	AppConfig.Etcd.EtcdAddr = conf.String("etcd::etcd_addr")
	if len(AppConfig.Etcd.EtcdAddr) == 0 {
		logs.Error("read etcd::etcd_addr from config-file failed")
		return fmt.Errorf("read etcd::etcd_addr from config-file failed")
	}
	AppConfig.Etcd.EtcdTimeout, err = conf.Int("etcd::etcd_timeout")
	if err != nil {
		logs.Error("read etcd::etcd_timeout from config-file failed, error: %v", err)
		return fmt.Errorf("read etcd::etcd_timeout from config-file failed, error: %v", err)
	}
	etcdSecKeyPrefix := conf.String("etcd::etcd_sec_key_prefix")
	if len(etcdSecKeyPrefix) == 0 {
		logs.Error("read etcd::etcd_sec_key_prefix from config-file failed")
		return fmt.Errorf("read etcd::etcd_sec_key_prefix from config-file failed")
	}
	if !strings.HasSuffix(etcdSecKeyPrefix, "/") {
		etcdSecKeyPrefix = etcdSecKeyPrefix + "/"
	}
	AppConfig.Etcd.EtcdSecKeyPrefix = etcdSecKeyPrefix
	etcdSecActivityKey := conf.String("etcd::etcd_sec_activity_key")
	if len(etcdSecActivityKey) == 0 {
		logs.Error("read etcd::etcd_sec_activity_key from config-file failed")
		return fmt.Errorf("read etcd::etcd_sec_activity_key from config-file failed")
	}
	AppConfig.Etcd.EtcdSecActivityKey = fmt.Sprintf("%s%s", etcdSecKeyPrefix, etcdSecActivityKey)
	etcdSecBlackListKey := conf.String("etcd::etcd_sec_black_list_key")
	if len(etcdSecBlackListKey) == 0 {
		logs.Error("read etcd::etcd_sec_black_list_key from config-file failed")
		return fmt.Errorf("read etcd::etcd_sec_black_list_key from config-file failed")
	}
	AppConfig.Etcd.EtcdSecBlackListKey = fmt.Sprintf("%s%s", etcdSecKeyPrefix, etcdSecBlackListKey)


	AppConfig.Proxy2LayerRedis.RedisAddr = conf.String("redis::redis_addr")
	if len(AppConfig.Proxy2LayerRedis.RedisAddr) == 0 {
		logs.Error("read redis::redis_addr from config-file failed")
		return fmt.Errorf("read redis::redis_addr from config-file failed")
	}
	AppConfig.Proxy2LayerRedis.RedisPassword = conf.String("redis::redis_password")
	if len(AppConfig.Proxy2LayerRedis.RedisPassword) == 0 {
		logs.Error("read redis::redis_password from config-file failed")
		return fmt.Errorf("read redis::redis_password from config-file failed")
	}
	AppConfig.Proxy2LayerRedis.RedisMaxIdle, err = conf.Int("redis::redis_max_idle")
	if err != nil {
		logs.Error("read redis::redis_max_idle from config-file failed, error: %v", err)
		return fmt.Errorf("read redis::redis_max_idle from config-file failed, error: %v", err)
	}
	AppConfig.Proxy2LayerRedis.RedisMaxActive, err = conf.Int("redis::redis_max_active")
	if err != nil {
		logs.Error("read redis::redis_max_active from config-file failed, error: %v", err)
		return fmt.Errorf("read redis::redis_max_active from config-file failed, error: %v", err)
	}
	AppConfig.Proxy2LayerRedis.RedisIdleTimeout, err = conf.Int("redis::redis_idle_timeout")
	if err != nil {
		logs.Error("read redis::redis_idle_timeout from config-file failed, error: %v", err)
		return fmt.Errorf("read redis::redis_idle_timeout from config-file failed, error: %v", err)
	}
	AppConfig.Proxy2LayerRedis.RedisProxy2LayerQueueKey = conf.String("redis::redis_proxy2layer_queue_key")
	if len(AppConfig.Proxy2LayerRedis.RedisAddr) == 0 {
		logs.Error("read redis::redis_proxy2layer_queue_key from config-file failed")
		return fmt.Errorf("read redis::redis_proxy2layer_queue_key from config-file failed")
	}
	AppConfig.Proxy2LayerRedis.RedisLayer2ProxyQueueKey = conf.String("redis::redis_layer2proxy_queue_key")
	if len(AppConfig.Proxy2LayerRedis.RedisAddr) == 0 {
		logs.Error("read redis::redis_layer2proxy_queue_key from config-file failed")
		return fmt.Errorf("read redis::redis_layer2proxy_queue_key from config-file failed")
	}

	AppConfig.Layer2ProxyRedis.RedisAddr = AppConfig.Proxy2LayerRedis.RedisAddr
	AppConfig.Layer2ProxyRedis.RedisPassword = AppConfig.Proxy2LayerRedis.RedisPassword
	AppConfig.Layer2ProxyRedis.RedisMaxIdle = AppConfig.Proxy2LayerRedis.RedisMaxIdle
	AppConfig.Layer2ProxyRedis.RedisIdleTimeout = AppConfig.Proxy2LayerRedis.RedisIdleTimeout
	AppConfig.Layer2ProxyRedis.RedisMaxActive = AppConfig.Proxy2LayerRedis.RedisMaxActive
	AppConfig.Layer2ProxyRedis.RedisProxy2LayerQueueKey = AppConfig.Proxy2LayerRedis.RedisProxy2LayerQueueKey
	AppConfig.Layer2ProxyRedis.RedisLayer2ProxyQueueKey = AppConfig.Proxy2LayerRedis.RedisLayer2ProxyQueueKey

	AppConfig.HandleUserGoroutineNum, err = conf.Int("service::handle_user_goroutine_num")
	if err != nil {
		logs.Error("read service::handle_user_goroutine_num from config-file failed, error: %v", err)
		return fmt.Errorf("read service::handle_user_goroutine_num from config-file failed, error: %v", err)
	}
	AppConfig.WriteProxy2LayerGoroutineNum, err = conf.Int("service::write_proxy2layer_goroutine_num")
	if err != nil {
		logs.Error("read service::write_proxy2layer_goroutine_num from config-file failed, error: %v", err)
		return fmt.Errorf("read service::write_proxy2layer_goroutine_num from config-file failed, error: %v", err)
	}
	AppConfig.ReadLayer2ProxyGoroutineNum, err = conf.Int("service::read_layer2proxy_goroutine_num")
	if err != nil {
		logs.Error("read service::read_layer2proxy_goroutine_num from config-file failed, error: %v", err)
		return fmt.Errorf("read service::read_layer2proxy_goroutine_num from config-file failed, error: %v", err)
	}

	AppConfig.Read2HandleChanSize, err = conf.Int("service::read2handle_chan_size")
	if err != nil {
		logs.Error("read service::read2handle_chan_size from config-file failed, error: %v", err)
		return fmt.Errorf("read service::read2handle_chan_size from config-file failed, error: %v", err)
	}
	AppConfig.Handle2WriteChanSize, err = conf.Int("service::handle2write_chan_size")
	if err != nil {
		logs.Error("read service::handle2write_chan_size from config-file failed, error: %v", err)
		return fmt.Errorf("read service::handle2write_chan_size from config-file failed, error: %v", err)
	}
	AppConfig.MaxRequestWaitTimeout, err = conf.Int("service::max_request_wait_timeout")
	if err != nil {
		logs.Error("read service::max_request_wait_timeout from config-file failed, error: %v", err)
		return fmt.Errorf("read service::max_request_wait_timeout from config-file failed, error: %v", err)
	}
	AppConfig.SendToHandleChanTimeout, err = conf.Int("service::send_to_handle_chan_timeout")
	if err != nil {
		logs.Error("read service::send_to_handle_chan_timeout from config-file failed, error: %v", err)
		return fmt.Errorf("read service::send_to_handle_chan_timeout from config-file failed, error: %v", err)
	}
	AppConfig.SendToWriteChanTimeout, err = conf.Int("service::send_to_write_chan_timeout")
	if err != nil {
		logs.Error("read service::send_to_write_chan_timeout from config-file failed, error: %v", err)
		return fmt.Errorf("read service::send_to_write_chan_timeout from config-file failed, error: %v", err)
	}

	AppConfig.SeckillTokenPasswd = conf.String("service::seckill_token_passwd")
	if len(AppConfig.SeckillTokenPasswd) != 64 {
		logs.Error("read service::seckill_token_passwd from config-file failed! error: length is wrong")
		return fmt.Errorf("read service::seckill_token_passwd from config-file failed! error: length is wrong")
	}

	return
}

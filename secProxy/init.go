package main

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"github.com/garyburd/redigo/redis"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"golang.org/x/net/context"
	"seckillsystem/secProxy/service"
	"time"
)

var (
	etcdClient *clientv3.Client
)

func initLogger() (err error) {
	config := make(map[string]interface{})
	config["filename"] = seckillConf.Log.LogPath
	config["level"]    = convertLogLevel(seckillConf.Log.LogLevel)

	bytes, err := json.Marshal(config)
	if err != nil {
		return
	}
	// 日志记录调用的文件名和文件行号
	logs.EnableFuncCallDepth(true)
	// 自定义log日志的记录方式
	return logs.SetLogger(logs.AdapterFile, string(bytes))
}
func convertLogLevel(level string) int {
	switch level {
	case "debug":
		return logs.LevelDebug
	case "warn":
		return logs.LevelWarn
	case "info":
		return logs.LevelInfo
	case "trace":
		return logs.LevelTrace
	}
	return logs.LevelDebug
}

func initRedis() error {
	redisPool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", seckillConf.Redis.RedisAddr, redis.DialPassword(seckillConf.Redis.RedisPassword))
		},
		MaxIdle:         seckillConf.Redis.RedisMaxIdle,
		MaxActive:       seckillConf.Redis.RedisMaxActive,
		IdleTimeout:     time.Duration(seckillConf.Redis.RedisIdleTimeout) * time.Second,
	}

	conn := redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("ping")
	if err != nil {
		logs.Error("ping redis failed! error: %v", err)
		return err
	}

	return nil
}

func initEtcd() (err error) {
	etcdClient, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{seckillConf.Etcd.EtcdAddr},
		DialTimeout: time.Duration(seckillConf.Etcd.EtcdTimeout) * time.Second,
	})
	if err != nil {
		logs.Error("connect etcd failed! error: %v", err)
		return err
	} else {
		logs.Info("connect etcd success!")
	}
	return
}

// 加载秒杀商品相关的配置信息
func loadSeckillConf() (err error) {
	// 从etcd服务读取秒杀活动信息
	response, err := etcdClient.Get(context.Background(), seckillConf.Etcd.EtcdSecActivityKey)
	if err != nil {
		return
	}
	logs.Debug("got activity from etcd success! activity list: %v", response.Kvs)

	var secActivityList []service.SecActivityConf
	for _, v := range response.Kvs {
		err = json.Unmarshal(v.Value, &secActivityList)
		if err != nil {
			return
		}
	}
	logs.Debug("activity list from etcd: %v", secActivityList)

	updateSecActivityInfoList(secActivityList)
	return
}

func initSeckill() (err error) {
	err = initLogger()
	if err != nil {
		logs.Error("init logger failed! error: %v", err)
		return
	}

	//err = initRedis()
	//if err != nil {
	//	logs.Error("init redis failed! error: %v", err)
	//	return
	//}

	err = initEtcd()
	if err != nil {
		logs.Error("init etcd failed! error: %v", err)
		return
	}

	err = loadSeckillConf()
	if err != nil {
		logs.Error("load seckill config failed! error: %v", err)
		return
	}

	err = service.InitService(seckillConf)
	if err != nil {
		logs.Error("init service failed! error: %v", err)
		return
	}
	go watchSecActivityChange(seckillConf.Etcd.EtcdSecActivityKey)

	logs.Info("init seckill success!")
	return
}

// 监听etcd中秒杀活动数据的变化
func watchSecActivityChange(key string) {
	for {
		watchChan := etcdClient.Watch(context.Background(), key)

		var secActivityList []service.SecActivityConf
		var getConfSucc = true
		for watchResp := range watchChan {
			for _, ev := range watchResp.Events {
				if ev.Type == mvccpb.DELETE {
					logs.Warn("key[%s]'s config deleted.", key)
					continue
				}

				if ev.Type == mvccpb.PUT && string(ev.Kv.Key) == key {
					err := json.Unmarshal(ev.Kv.Value, &secActivityList)
					if err != nil {
						logs.Error("key[%s],json Unmarshal[%s] failed! error: %v", key, string(ev.Kv.Value), err)
						getConfSucc = false
						continue
					}
				}
				logs.Debug("get config from etcd, %s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			}

			if getConfSucc {
				logs.Debug("get config from etcd success! config: %v", secActivityList)
				updateSecActivityInfoList(secActivityList)
			}
		}
	}
}

func updateSecActivityInfoList(secActivityList []service.SecActivityConf) {
	tmpActivityListMap := make(map[int]*service.SecActivityConf)
	for _, v := range secActivityList {
		activity := v  // 解决bug
		tmpActivityListMap[v.ActivityId] = &activity
	}

	seckillConf.RwLock.Lock()
	seckillConf.SecActivityListMap = tmpActivityListMap
	seckillConf.RwLock.Unlock()
}
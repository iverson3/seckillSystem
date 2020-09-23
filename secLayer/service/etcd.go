package service

import (
	"context"
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"go.etcd.io/etcd/clientv3"
	"time"
)

func initEtcd(conf *SecLayerConf) (err error) {
	secLayerContext.EtcdClient, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{conf.Etcd.EtcdAddr},
		DialTimeout: time.Duration(conf.Etcd.EtcdTimeout) * time.Second,
	})
	if err != nil {
		logs.Error("connect etcd failed! error: %v", err)
		return
	}
	return
}

// 从etcd服务读取秒杀活动数据
func loadActivityFromEtcd(conf *SecLayerConf) (err error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*10)
	response, err := secLayerContext.EtcdClient.Get(ctx, conf.Etcd.EtcdSecActivityKey)
	if err != nil {
		logs.Error("get [%s] from etcd failed! error: %v", conf.Etcd.EtcdSecActivityKey, err)
		return
	}
	cancelFunc()
	logs.Debug("got activity from etcd success! activity list: %v", response.Kvs)

	var secActivityList []SecActivityConf
	for _, v := range response.Kvs {
		err = json.Unmarshal(v.Value, &secActivityList)
		if err != nil {
			logs.Error("json Unmarshal seckill activity list failed! error: %v", err)
			return
		}
	}
	logs.Debug("activity list from etcd: %v", secActivityList)

	updateSecActivityList(conf, secActivityList)
	return
}

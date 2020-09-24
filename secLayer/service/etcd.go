package service

import (
	"github.com/astaxie/beego/logs"
	"go.etcd.io/etcd/clientv3"
	"time"
)

func initEtcd() (err error) {
	secLayerContext.EtcdClient, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{secLayerContext.SecLayerConfig.Etcd.EtcdAddr},
		DialTimeout: time.Duration(secLayerContext.SecLayerConfig.Etcd.EtcdTimeout) * time.Second,
	})
	if err != nil {
		logs.Error("connect etcd failed! error: %v", err)
	}
	return
}

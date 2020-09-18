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

// 从etcd服务读取秒杀商品数据
func loadProductFromEtcd(conf *SecLayerConf) (err error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*10)
	response, err := secLayerContext.EtcdClient.Get(ctx, conf.Etcd.EtcdSecProductKey)
	if err != nil {
		logs.Error("get [%s] from etcd failed! error: %v", conf.Etcd.EtcdSecProductKey, err)
		return
	}
	cancelFunc()

	logs.Debug("load product from etcd success")

	logs.Debug("response from etcd is [%v]", response.Kvs)

	var secProductInfoList []SecProductInfoConf
	for k, v := range response.Kvs {
		logs.Debug("key[%v] value[%v]", k, v)

		err = json.Unmarshal(v.Value, &secProductInfoList)
		if err != nil {
			logs.Error("json Unmarshal seckill product info failed! error: %v", err)
			return
		}

		logs.Debug("seckill product info is [%v]", secProductInfoList)
	}

	logs.Debug("product info from etcd: %v", secProductInfoList)

	updateSecProductInfoList(conf, secProductInfoList)
	return
}

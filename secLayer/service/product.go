package service

import (
	"context"
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"go.etcd.io/etcd/mvcc/mvccpb"
)

// 监听商品信息的变化
func watchSecProductChange(conf *SecLayerConf) {
	key := conf.Etcd.EtcdSecProductKey
	for {
		watchChan := secLayerContext.EtcdClient.Watch(context.Background(), key)

		var secProductInfoList []SecProductInfoConf
		var getConfSucc = true
		for watchResp := range watchChan {
			for _, ev := range watchResp.Events {
				if ev.Type == mvccpb.DELETE {
					logs.Warn("key[%s]'s config deleted.", key)
					continue
				}

				if ev.Type == mvccpb.PUT && string(ev.Kv.Key) == key {
					err := json.Unmarshal(ev.Kv.Value, &secProductInfoList)
					if err != nil {
						logs.Error("key[%s],json Unmarshal[%s] failed! error: %v", key, string(ev.Kv.Value), err)
						getConfSucc = false
						continue
					}
				}
				logs.Debug("get config from etcd, %s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			}

			if getConfSucc {
				logs.Debug("get config from etcd success! config: %v", secProductInfoList)
				updateSecProductInfoList(conf, secProductInfoList)
			}
		}
	}
}

func updateSecProductInfoList(conf *SecLayerConf, productList []SecProductInfoConf) {
	tmpProductListMap := make(map[int]*SecProductInfoConf)
	for _, v := range productList {
		product := v  // 解决bug
		product.SecSoldLimit = &SecLimit{}
		tmpProductListMap[v.ProductId] = &product
	}

	secLayerContext.SecProductRwLock.Lock()
	conf.SecProductInfo = tmpProductListMap
	secLayerContext.SecProductRwLock.Unlock()
}

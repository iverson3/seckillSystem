package service

import (
	"context"
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"go.etcd.io/etcd/mvcc/mvccpb"
)

// 监听etcd中秒杀活动数据的变化
func watchSecActivityChange(conf *SecLayerConf) {
	key := conf.Etcd.EtcdSecActivityKey
	for {
		watchChan := secLayerContext.EtcdClient.Watch(context.Background(), key)

		var secActivityList []SecActivityConf
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
				updateSecActivityList(conf, secActivityList)
			}
		}
	}
}

func updateSecActivityList(conf *SecLayerConf, activityList []SecActivityConf) {
	tmpActivityListMap := make(map[int]*SecActivityConf)
	for _, v := range activityList {
		activity := v  // 解决bug
		activity.SecSoldLimit = &SecLimit{}
		tmpActivityListMap[v.ActivityId] = &activity
	}

	secLayerContext.SecActivityRwLock.Lock()
	conf.SecActivityListMap = tmpActivityListMap
	secLayerContext.SecActivityRwLock.Unlock()
}

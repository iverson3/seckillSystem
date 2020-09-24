package service

import (
	"context"
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"time"
)

// 从etcd服务读取秒杀活动数据
func loadActivityFromEtcd() (err error) {
	secActivityList, err := getDataFromEtcd()
	if err != nil {
		return
	}

	updateSecActivityList(secActivityList)
	return
}

func getDataFromEtcd() (secActivityList []*SecActivityConf, err error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second * 10)
	response, err := secLayerContext.EtcdClient.Get(ctx, secLayerContext.SecLayerConfig.Etcd.EtcdSecActivityKey)
	cancelFunc()
	if err != nil {
		logs.Error("get [%s] from etcd failed! error: %v", secLayerContext.SecLayerConfig.Etcd.EtcdSecActivityKey, err)
		return
	}
	logs.Debug("got activity from etcd success! activity list: %v", response.Kvs)

	for _, v := range response.Kvs {
		err = json.Unmarshal(v.Value, &secActivityList)
		if err != nil {
			logs.Error("json Unmarshal seckill activity list failed! error: %v", err)
			return
		}
	}
	logs.Debug("activity list from etcd: %v", secActivityList)
	return
}

// 监听etcd中秒杀活动数据的变化
func watchSecActivityChange() {
	key := secLayerContext.SecLayerConfig.Etcd.EtcdSecActivityKey
	for {
		watchChan := secLayerContext.EtcdClient.Watch(context.Background(), key)

		var secActivityList []*SecActivityConf
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
				updateSecActivityList(secActivityList)
			}
		}
	}
}

// 将状态有变化的活动商品同步到etcd中
func syncActivityChangeToEtcd(activity *SecActivityConf) {
	var err error
	defer func() {
		if err != nil {
			logs.Error("sync activity's change data to etcd failed! activity: %v; error: %v", activity, err)
		}
	}()
	secActivityList, err := getDataFromEtcd()
	if err != nil {
		return
	}

	// 找到对应的活动，并进行数据修改
	for _, v := range secActivityList {
		if v.ActivityId == activity.ActivityId {
			v.Status = activity.Status
			break
		}
	}

	data, err := json.Marshal(secActivityList)
	if err != nil {
		return
	}

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second * 2)
	_, err = secLayerContext.EtcdClient.Put(ctx, secLayerContext.SecLayerConfig.Etcd.EtcdSecActivityKey, string(data))
	cancelFunc()
	if err != nil {
		return
	}
	logs.Error("sync activity's change data to etcd success! activity: %v", activity)
}

func updateSecActivityList(activityList []*SecActivityConf) {
	tmpActivityListMap := make(map[int]*SecActivityConf)
	for _, v := range activityList {
		activity := *v  // 解决bug
		activity.SecSoldLimit = &SecLimit{}
		tmpActivityListMap[v.ActivityId] = &activity
	}

	secLayerContext.SecActivityRwLock.Lock()
	secLayerContext.SecLayerConfig.SecActivityListMap = tmpActivityListMap
	secLayerContext.SecActivityRwLock.Unlock()
}

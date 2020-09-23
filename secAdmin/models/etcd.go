package models

import (
	"context"
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"go.etcd.io/etcd/clientv3"
	"time"
)

var (
	EtcdClient *clientv3.Client
	EtcdConf EtcdConfig
)

type EtcdConfig struct {
	Addr        string
	Timeout     int
	KeyPrefix   string
	ActivityKey string
}

type EtcdModel struct {

}

func NewEtcdModel() *EtcdModel {
	return &EtcdModel{}
}

func SetEtcd(etcd *clientv3.Client, etcdConf EtcdConfig) {
	EtcdClient = etcd
	EtcdConf = etcdConf
}

func (this *EtcdModel) SyncActivityToEtcd(activity *Activity) (err error) {
	SeckillConfs, err := getSecActivityList()
	if err != nil {
		return
	}

	SeckillConfs = append(SeckillConfs,
		SecActivityConf{
			ActivityId: activity.Id,
			ProductId: activity.ProductId,
			Total:     activity.Total,
			Left:      activity.Total,
			Status:    activity.Status,
			StartTime: activity.StartTime,
			EndTime:   activity.EndTime,
			BuyRate:   activity.BuyRate,
			UserMaxBuyLimit: activity.BuyLimit,
			MaxSoldLimit: activity.SoldLimitSecond,
		})

	logs.Error("will put to etcd: %v", SeckillConfs)

	data, err := json.Marshal(SeckillConfs)
	if err != nil {
		return
	}

	logs.Error("will put to etcd: %s", string(data))

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	_, err = EtcdClient.Put(ctx, EtcdConf.ActivityKey, string(data))
	cancel()
	if err != nil {
		logs.Error("sync seckill activity to etcd failed! error: %v", err)
		return
	}
	return
}

// 从etcd服务读取秒杀活动数据
func getSecActivityList() (secActivityList []SecActivityConf, err error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second * 10)
	response, err := EtcdClient.Get(ctx, EtcdConf.ActivityKey)
	cancelFunc()
	if err != nil {
		logs.Error("get [%s] from etcd failed! error: %v", EtcdConf.ActivityKey, err)
		return
	}
	logs.Debug("got activity list from etcd success！ activity list: %v", response.Kvs)

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



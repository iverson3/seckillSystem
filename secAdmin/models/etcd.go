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
	Addr string
	Timeout int
	KeyPrefix string
	ProductKey string
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
			ProductId: activity.ProductId,
			Total:     activity.Total,
			Left:      activity.Total,
			Status:    activity.Status,
			StartTime: int(activity.StartTime),
			EndTime:   int(activity.EndTime),
			BuyRate:   activity.BuyRate,
			UserMaxBuyLimit: activity.BuyLimit,
			MaxSoldLimit: activity.SoldLimitSecond,
		})

	data, err := json.Marshal(SeckillConfs)
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	_, err = EtcdClient.Put(ctx, EtcdConf.ProductKey, string(data))
	cancel()
	if err != nil {
		logs.Error("sync seckill activity to etcd failed! error: %v", err)
		return
	}
	return
}

// 从etcd服务读取秒杀商品数据
func getSecActivityList() (secActivityList []SecActivityConf, err error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second * 10)
	response, err := EtcdClient.Get(ctx, EtcdConf.ProductKey)
	cancelFunc()
	if err != nil {
		logs.Error("get [%s] from etcd failed! error: %v", EtcdConf.ProductKey, err)
		return
	}
	logs.Debug("load activity list from etcd success")
	logs.Debug("response from etcd is [%v]", response.Kvs)

	for k, v := range response.Kvs {
		logs.Debug("key[%v] value[%v]", k, v)

		err = json.Unmarshal(v.Value, &secActivityList)
		if err != nil {
			logs.Error("json Unmarshal seckill activity list failed! error: %v", err)
			return
		}

		logs.Debug("seckill activity list is [%v]", secActivityList)
	}

	logs.Debug("activity list from etcd: %v", secActivityList)
	return
}



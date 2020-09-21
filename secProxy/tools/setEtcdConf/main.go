package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"go.etcd.io/etcd/clientv3"
	"time"
)

const EtcdKey = "/oldboy/backend/seckill/product"

// 临时把这个结构体拷贝到这使用
type SeckillInfoConf struct {
	ProductId int
	Total int
	Left int
	Status int
	StartTime int
	EndTime int
	BuyRate float64     // 秒杀成功的概率 (用户到达秒杀系统逻辑层 能够抢到该商品的概率)
	UserMaxBuyLimit int // 对于当前商品，每个用户最多可以购买的数量
	MaxSoldLimit int    // 商品每秒的秒杀数量限制
}

// 向etcd服务中添加几个测试数据
func SetSeckillConfToEtcd() {
	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"47.107.149.234:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		logs.Error("connect etcd failed! error: %v", err)
		return
	}
	defer etcdClient.Close()

	var SeckillConfs []SeckillInfoConf
	SeckillConfs = append(SeckillConfs,
		SeckillInfoConf{
			ProductId: 1028,
			Total:     100000,
			Left:      100000,
			Status:    0,
			StartTime: 1600653677,
			EndTime:   1600689677,
			BuyRate:   0.8,
			UserMaxBuyLimit: 1,
			MaxSoldLimit: 1000,
		})

	SeckillConfs = append(SeckillConfs,
		SeckillInfoConf{
			ProductId: 1029,
			Total:     200000,
			Left:      200000,
			Status:    0,
			StartTime: 1600653677,
			EndTime:   1600689677,
			BuyRate:   0.5,
			UserMaxBuyLimit: 2,
			MaxSoldLimit: 1000,
		})

	data, err := json.Marshal(SeckillConfs)
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	_, err = etcdClient.Put(ctx, EtcdKey, string(data))
	cancel()
	if err != nil {
		fmt.Println("etcd put failed! error: ", err)
		return
	}

	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	resp, err := etcdClient.Get(ctx, EtcdKey)
	cancel()
	if err != nil {
		fmt.Println("etcd get failed! error: ", err)
		return
	}

	for _, ev := range resp.Kvs {
		fmt.Printf("%s : %s \n", ev.Key, ev.Value)
	}
}

func main()  {
	SetSeckillConfToEtcd()
}

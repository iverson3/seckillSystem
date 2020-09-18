package service

import (
	"github.com/astaxie/beego/logs"
)

func InitSecLayer(conf *SecLayerConf) (err error) {
	err = initRedis(conf)
	if err != nil {
		logs.Error("init redis failed! error: %v", err)
		return
	}
	logs.Debug("init redis succ")

	err = initEtcd(conf)
	if err != nil {
		logs.Error("init etcd failed! error: %v", err)
		return
	}
	logs.Debug("init etcd succ")

	err = loadProductFromEtcd(conf)
	if err != nil {
		logs.Error("load product from etcd failed! error: %v", err)
		return
	}
	logs.Debug("load product succ")

	secLayerContext.SecLayerConfig = conf
	secLayerContext.Read2HandleChan  = make(chan *SecRequest, conf.Read2HandleChanSize)
	secLayerContext.Handle2WriteChan = make(chan *SecResponse, conf.Handle2WriteChanSize)
	secLayerContext.UserBuyHistoryMap = make(map[int]*UserBuyHistory, 100000)
	secLayerContext.ProductCountManager = NewProductCountMgr()

	go watchSecProductChange(conf)

	return
}

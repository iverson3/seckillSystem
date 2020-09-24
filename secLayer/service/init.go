package service

import (
	"github.com/astaxie/beego/logs"
)

func InitSecLayer(conf *SecLayerConf) (err error) {
	secLayerContext.SecLayerConfig = conf
	secLayerContext.Read2HandleChan  = make(chan *SecRequest, conf.Read2HandleChanSize)
	secLayerContext.Handle2WriteChan = make(chan *SecResponse, conf.Handle2WriteChanSize)
	secLayerContext.UserBuyHistoryMap = make(map[int]*UserBuyHistory, 100000)
	secLayerContext.ProductCountManager = NewProductCountMgr()

	err = initRedis()
	if err != nil {
		logs.Error("init redis failed! error: %v", err)
		return
	}
	logs.Debug("init redis succ")

	err = initEtcd()
	if err != nil {
		logs.Error("init etcd failed! error: %v", err)
		return
	}
	logs.Debug("init etcd succ")

	err = loadActivityFromEtcd()
	if err != nil {
		logs.Error("load activity list from etcd failed! error: %v", err)
		return
	}
	logs.Debug("load activity list success!")

	go watchSecActivityChange()
	return
}

package main

import (
	"github.com/astaxie/beego/logs"
	"seckillsystem/secLayer/service"
)

func main() {
	err := initConfig("ini", "./conf/seclayer.conf")
	if err != nil {
		logs.Error("init config failed! error: %v", err)
		panic(err)
		return
	}
	logs.Debug("app config: %v", AppConfig)

	err = initLogger()
	if err != nil {
		logs.Error("init logger failed! error: %v", err)
		panic(err)
		return
	}
	logs.Debug("init log succ")

	err = service.InitSecLayer(AppConfig)
	if err != nil {
		logs.Error("init seckill failed! error: %v", err)
		panic(err)
		return
	}

	logs.Info("service begin run")
	service.Run()
}

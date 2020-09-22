package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	_ "seckillsystem/secAdmin/routers"
)

func main() {
	err := Init()
	if err != nil {
		logs.Error("init failed! error: %v", err)
		return
	}

	beego.Run()
}


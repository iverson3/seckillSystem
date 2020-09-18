package main

import (
	"github.com/astaxie/beego"
	_ "seckillsystem/secProxy/router"
)

func main()  {
	err := initConfig()
	if err != nil {
		panic(err)
	}

	err = initSeckill()
	if err != nil {
		panic(err)
	}

	beego.Run()
}

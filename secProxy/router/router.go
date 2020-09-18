package router

import (
	"github.com/astaxie/beego"
	"seckillsystem/secProxy/controller"
)

func init()  {
	beego.Router("/seckill", &controller.SeckillController{}, "*:SecKill")
	beego.Router("/secinfo", &controller.SeckillController{}, "*:SecInfo")
}

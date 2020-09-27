package routers

import (
	"seckillsystem/secAdmin/controllers"
	"github.com/astaxie/beego"
)

func init() {
    //beego.Router("/", &controllers.MainController{})

	beego.Router("/", &controllers.ProductController{}, "GET:ListProduct")
	beego.Router("/product", &controllers.ProductController{}, "GET:ListProduct")
	beego.Router("/product/create", &controllers.ProductController{}, "*:CreateProduct")

	beego.Router("/activity", &controllers.ActivityController{}, "GET:ListActivity")
	beego.Router("/activity/create", &controllers.ActivityController{}, "*:CreateActivity")
	beego.Router("/activity/update/status", &controllers.ActivityController{}, "POST:UpdateActivityStatus")

	beego.Router("/activity/delete", &controllers.ActivityController{}, "POST:DeleteActivity")
	beego.Router("/activity/product/left", &controllers.ActivityController{}, "POST:GetActivityProductLeft")
}

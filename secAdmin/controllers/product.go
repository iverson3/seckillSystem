package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"seckillsystem/secAdmin/models"
)

type ProductController struct {
	beego.Controller
}

func (this *ProductController) ListProduct() {
	productModel := models.NewProductModel()

	list, err := productModel.GetProductList()
	if err != nil {
		logs.Warn("get product list failed! error: %v", err)
		return
	}

	logs.Info("product list: %v", list)

	this.Data["product_list"] = list
	this.Layout = "layout/layout.html"
	this.TplName = "product/list.html"
}

func (this *ProductController) CreateProduct() {
	if this.Ctx.Input.IsPost() {
		var err error
		defer func() {
			if err != nil {
				this.Data["ErrorTitle"] = "创建商品出错"
				this.Data["Error"] = err.Error()
				this.Layout = "layout/layout.html"
				this.TplName = "error/form_error.html"
			}
		}()

		name := this.GetString("product_name")
		total, err := this.GetInt("total")
		if name == "" || err != nil || total == 0 {
			logs.Warn("form validate failed!")
			//this.Redirect("/product/create", 401)
			return
		}

		status, err := this.GetInt("status")
		if err != nil {
			logs.Warn("form validate failed!")
			//this.Redirect("/product/create", 401)
			return
		}

		product := &models.Product{
			ProductName: name,
			Total:       total,
			Status:      status,
		}
		model := models.NewProductModel()
		err = model.CreateNewProduct(product)
		if err != nil {
			return
		}

		//this.Ctx.Redirect(200, "/product")
		this.Redirect("/product", 302)
		return
	}

	this.Layout = "layout/layout.html"
	this.TplName = "product/create.html"
}

package models

import (
	"github.com/astaxie/beego/logs"
)

type ProductModel struct {

}

type Product struct {
	ProductId int `db:"id"`
	ProductName string `db:"name"`
	Total int `db:"total"`
	Status int `db:"status"`
}

func NewProductModel() *ProductModel {
	productModel := &ProductModel{}
	return productModel
}

func (this *ProductModel) GetProductList() (list []*Product, err error) {
	sql := "select * from product"
	err = Db.Select(&list, sql)
	if err != nil {
		logs.Error("select product_list from mysql failed! sql: %v, error: %v", sql, err)
		return
	}
	return
}

func (this *ProductModel) CreateNewProduct(p *Product) (err error) {
	sql := "insert into product(name,total,status) values(?,?,?)"
	_, err = Db.Exec(sql, p.ProductName, p.Total, p.Status)
	if err != nil {
		logs.Error("insert data into product failed! sql: %v, error: %v", sql, err)
		return
	}
	return
}

func (this *ProductModel) ProductExists(pid int) (exist bool, err error, product *Product) {
	var productList []*Product
	sql := "select * from product where id = ?"
	err = Db.Select(&productList, sql, pid)
	if err != nil {
		logs.Error("select product_list from mysql failed! sql: %v, error: %v", sql, err)
		return
	}
	if len(productList) == 0 {
		return
	}

	return true, nil, productList[0]
}
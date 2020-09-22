package main

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/jmoiron/sqlx"
	"go.etcd.io/etcd/clientv3"
	"seckillsystem/secAdmin/models"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

var (
	Db *sqlx.DB
	EtcdClient *clientv3.Client
)

func Init() (err error) {
	err = initConfig()
	if err != nil {
		return
	}
	err = initDb()
	if err != nil {
		return
	}
	err = initEtcd()
	if err != nil {
		return
	}

	models.SetDb(Db)
	models.SetEtcd(EtcdClient, AppConf.etcd)
	return
}

func initDb() (err error) {
	dns := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", AppConf.mysql.UserName, AppConf.mysql.PassWd, AppConf.mysql.Host, AppConf.mysql.Port, AppConf.mysql.Database)
	Db, err = sqlx.Open("mysql", dns)
	if err != nil {
		logs.Error("open mysql failed! error: %v", err)
		return
	}
	logs.Info("connect to mysql success!")
	return
}

func initEtcd() (err error) {
	EtcdClient, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{AppConf.etcd.Addr},
		DialTimeout: time.Duration(AppConf.etcd.Timeout) * time.Second,
	})
	if err != nil {
		logs.Error("connect to etcd failed! error: %v", err)
		return
	}
	logs.Info("connect to etcd success!")
	return
}
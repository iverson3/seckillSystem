package main

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/garyburd/redigo/redis"
	"github.com/jmoiron/sqlx"
	"go.etcd.io/etcd/clientv3"
	"seckillsystem/secAdmin/models"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

var (
	Db *sqlx.DB
	EtcdClient *clientv3.Client
	RedisPool *redis.Pool
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
	err = initRedis()
	if err != nil {
		return
	}

	models.SetDb(Db)
	models.SetEtcd(EtcdClient, AppConf.etcd)
	models.SetRedis(RedisPool, AppConf.redis)

	go models.WatchEtcdActivityChange()
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

func initRedis() (err error) {
	RedisPool = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", AppConf.redis.Addr, redis.DialPassword(AppConf.redis.PassWd))
		},
		MaxIdle:         AppConf.redis.MaxIdle,
		MaxActive:       AppConf.redis.MaxActive,
		IdleTimeout:     time.Duration(AppConf.redis.IdleTimeout) * time.Second,
	}

	conn := RedisPool.Get()
	defer conn.Close()
	_, err = conn.Do("ping")
	if err != nil {
		logs.Error("ping redis failed! error: %v", err)
		return
	}
	return
}

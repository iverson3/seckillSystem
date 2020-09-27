package main

import (
	"fmt"
	"github.com/astaxie/beego"
	"seckillsystem/secAdmin/models"
	"strings"
)

var (
	AppConf Config
)

type Config struct {
	mysql models.MysqlConfig
	etcd models.EtcdConfig
	redis models.RedisConfig
}

func initConfig() (err error) {
	mysqlConf, err := initMysqlConf()
	if err != nil {
		return
	}

	etcdConf, err := initEtcdConf()
	if err != nil {
		return
	}

	redisConf, err := initRedisConf()
	if err != nil {
		return
	}

	AppConf = Config{
		mysql: mysqlConf,
		etcd: etcdConf,
		redis: redisConf,
	}
	return
}

func initMysqlConf() (models.MysqlConfig, error) {
	username  := beego.AppConfig.String("mysql_username")
	passwd    := beego.AppConfig.String("mysql_passwd")
	host      := beego.AppConfig.String("mysql_host")
	port, err := beego.AppConfig.Int("mysql_port")
	dbName    := beego.AppConfig.String("mysql_database")

	if len(username) == 0 || len(passwd) == 0 || len(host) == 0 || err != nil || port == 0 || len(dbName) == 0 {
		return models.MysqlConfig{}, fmt.Errorf("init database config failed, some config field is null")
	}

	mysqlConf := models.MysqlConfig{
		UserName: username,
		PassWd:   passwd,
		Host:     host,
		Port:     port,
		Database: dbName,
	}
	return mysqlConf, nil
}

func initEtcdConf() (models.EtcdConfig, error) {
	addr         := beego.AppConfig.String("etcd_addr")
	keyPrefix    := beego.AppConfig.String("etcd_sec_key_prefix")
	activityKey  := beego.AppConfig.String("etcd_sec_activity_key")
	timeout, err := beego.AppConfig.Int("etcd_timeout")
	if len(addr) == 0 || len(keyPrefix) == 0 || len(activityKey) == 0 || err != nil || timeout == 0 {
		return models.EtcdConfig{}, fmt.Errorf("init etcd config failed, some config field is null")
	}
	if !strings.HasSuffix(keyPrefix, "/") {
		keyPrefix = keyPrefix + "/"
	}

	etcd := models.EtcdConfig{
		Addr:        addr,
		Timeout:     timeout,
		KeyPrefix:   keyPrefix,
		ActivityKey: fmt.Sprintf("%s%s", keyPrefix, activityKey),
	}
	return etcd, nil
}

func initRedisConf() (models.RedisConfig, error) {
	addr         := beego.AppConfig.String("redis_addr")
	passWd       := beego.AppConfig.String("redis_password")
	leftKey      := beego.AppConfig.String("redis_product_left_key")
	if len(addr) == 0 || len(passWd) == 0 || len(leftKey) == 0 {
		return models.RedisConfig{}, fmt.Errorf("init redis config failed, some config field is null")
	}

	maxIdle, err1     := beego.AppConfig.Int("redis_max_idle")
	maxActive, err2   := beego.AppConfig.Int("redis_max_active")
	idleTimeout, err3 := beego.AppConfig.Int("redis_idle_timeout")
	if err1 != nil || err2 != nil || err3 != nil {
		return models.RedisConfig{}, fmt.Errorf("init redis config failed, some config field is null")
	}

	redis := models.RedisConfig{
		Addr:           addr,
		PassWd:         passWd,
		MaxIdle:        maxIdle,
		MaxActive:      maxActive,
		IdleTimeout:    idleTimeout,
		ProductLeftKey: leftKey,
	}
	return redis, nil
}
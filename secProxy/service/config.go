package service

import (
	"github.com/garyburd/redigo/redis"
	"sync"
	"time"
)

var (
	ProductStatusNormal       = 0
	ProductStatusSoldOut      = 1
	ProductStatusForceSoldOut = 2
)

type RedisConf struct {
	RedisAddr string
	RedisPassword string
	RedisMaxIdle int
	RedisMaxActive int
	RedisIdleTimeout int
	RedisProxy2LayerQueueKey string
	RedisLayer2ProxyQueueKey string
}

type EtcdConf struct {
	EtcdAddr string
	EtcdTimeout int
	EtcdSecKeyPrefix string
	EtcdSecProductKey string
}

type LogConf struct {
	LogPath string
	LogLevel string
}

type SeckillConf struct {
	Redis RedisConf
	BlackListRedis RedisConf
	Etcd EtcdConf
	Log LogConf
	SecProductInfo map[int]*SecProductInfoConf
	RwLock sync.RWMutex  // 读写锁
	CookieSecretKey string
	UserAccessLimitPerSecond int
	IpAccessLimitPerSecond int
	RefererWhiteList []string

	IpBlackList map[string]bool
	IdBlackList map[int]bool
	BlackRwLock sync.RWMutex  // 读写锁

	BlackRedisPool *redis.Pool
	Proxy2LayerRedisPool *redis.Pool

	WriteProxy2LayerGoroutineNum int
	ReadLayer2ProxyGoroutineNum int
	MaxRequestWaitTimeout int

	SecReqChan chan *SecRequest
	SecReqChanSize int

	SecRequestMap map[int]*SecRequest
	ReqMapLock sync.RWMutex
}

// 秒杀商品的相关信息结构
type SecProductInfoConf struct {
	ProductId int
	Total int
	Left int
	Status int
	StartTime int64
	EndTime int64
}

// 秒杀请求的相关信息 (请求参数 ip地址 请求时间等)
type SecRequest struct {
	UserId int
	UserAuthSign string
	ProductId int
	Source string
	AuthCode string
	SecTime string
	Nance string
	AccessTime time.Time   // 请求到达服务器的时间 (用作检测恶意请求使用)
	ClientAddr string
	ClientReferer string

	//CloseNotify <-chan bool `json:"-"`       // json序列化的时候忽略该字段
	ResultChan chan *SecResponse `json:"-"`  // json序列化的时候忽略该字段
}

type SecResponse struct {
	UserId int
	ProductId int
	Token string
	TokenTime int64
	Code int
}
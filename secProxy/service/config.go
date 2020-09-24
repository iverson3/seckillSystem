package service

import (
	"github.com/garyburd/redigo/redis"
	"sync"
	"time"
)

const (
	ActivityStatusNormal  = 0    // 正常可用
	ActivityStatusDisable = 1    // 禁用
	ActivityStatusSoldOut = 2    // 售罄
	ActivityStatusExpire  = 3    // 过期或结束
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
	EtcdAddr           string
	EtcdTimeout        int
	EtcdSecKeyPrefix   string
	EtcdSecActivityKey string
}

type LogConf struct {
	LogPath string
	LogLevel string
}

type SeckillConf struct {
	Redis                    RedisConf
	BlackListRedis           RedisConf
	Etcd                     EtcdConf
	Log                      LogConf
	SecActivityListMap       map[int]*SecActivityConf
	RwLock                   sync.RWMutex  // 读写锁
	CookieSecretKey          string
	UserAccessLimitPerSecond int
	IpAccessLimitPerSecond   int
	RefererWhiteList         []string

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

// 秒杀活动的相关信息结构
type SecActivityConf struct {
	ActivityId int
	ProductId int
	Total int
	Left int
	Status int
	StartTime int64
	EndTime int64
	BuyRate float64     // 秒杀成功的概率 (用户到达秒杀系统逻辑层 能够抢到该商品的概率)
	UserMaxBuyLimit int // 对于当前商品，每个用户最多可以购买的数量
	MaxSoldLimit int    // 商品每秒的秒杀数量限制
}

// 秒杀请求的相关信息 (请求参数 ip地址 请求时间等)
type SecRequest struct {
	UserId        int
	UserAuthSign  string
	ActivityId    int
	Source        string
	AuthCode      string
	SecTime       string
	Nance         string
	AccessTime    time.Time   // 请求到达服务器的时间 (用作检测恶意请求使用)
	ClientAddr    string
	ClientReferer string

	//CloseNotify <-chan bool `json:"-"`       // json序列化的时候忽略该字段
	ResultChan chan *SecResponse `json:"-"`  // json序列化的时候忽略该字段
}

type SecResponse struct {
	UserId     int
	ActivityId int
	Token      string
	TokenTime  int64
	Code       int
}
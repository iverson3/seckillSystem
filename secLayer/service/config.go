package service

import (
	"github.com/garyburd/redigo/redis"
	"go.etcd.io/etcd/clientv3"
	"sync"
	"time"
)

var (
	secLayerContext = &SecLayerContext{}
)

type RedisConf struct {
	RedisAddr string
	RedisPassword string
	RedisMaxIdle int
	RedisMaxActive int
	RedisIdleTimeout int
	RedisProxy2LayerQueueKey string
	RedisLayer2ProxyQueueKey string
	ProductLeftKey string
}

type EtcdConf struct {
	EtcdAddr            string
	EtcdTimeout         int
	EtcdSecKeyPrefix    string
	EtcdSecActivityKey  string
	EtcdSecBlackListKey string
}

type LogConf struct {
	LogPath string
	LogLevel string
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

	SecSoldLimit *SecLimit `json:"-"`  // 限速控制器
}

type SecLayerConf struct {
	Proxy2LayerRedis RedisConf
	Layer2ProxyRedis RedisConf
	Etcd             EtcdConf
	Log              LogConf

	HandleUserGoroutineNum       int
	WriteProxy2LayerGoroutineNum int
	ReadLayer2ProxyGoroutineNum  int
	Read2HandleChanSize          int
	Handle2WriteChanSize         int
	MaxRequestWaitTimeout        int
	SendToHandleChanTimeout      int
	SendToWriteChanTimeout       int
	SecActivityListMap           map[int]*SecActivityConf
	SeckillTokenPasswd           string
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
}

// 对用户秒杀请求处理之后的结果
type SecResponse struct {
	UserId     int
	ActivityId int
	Token      string
	TokenTime  int64
	Code       int
}

type SecLayerContext struct {
	Proxy2LayerRedisPool *redis.Pool
	Layer2ProxyRedisPool *redis.Pool

	EtcdClient        *clientv3.Client
	SecActivityRwLock sync.RWMutex
	SecWaitGroup      sync.WaitGroup

	Read2HandleChan chan *SecRequest
	Handle2WriteChan chan *SecResponse

	UserBuyHistoryMap map[int]*UserBuyHistory // 记录所有用户的购买历史
	HistoryMapLock sync.RWMutex

	ProductCountManager *ProductCountMgr

	SecLayerConfig *SecLayerConf
}

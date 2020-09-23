package models

import (
	"github.com/astaxie/beego/logs"
	"time"
)

const (
	ActivityStatusNormal  = 0
	ActivityStatusDisable = 1
	ActivityStatusExpire  = 2
)

type Activity struct {
	Id int `db:"id" form:"-" json:"id"`
	Name string `db:"name" form:"name"`
	ProductId int `db:"pid" form:"pid"`
	StartTime int64 `db:"start_time" form:"-"`
	EndTime int64 `db:"end_time" form:"-"`
	BuyRate float64 `db:"buy_rate" form:"buy_rate"`
	Total int `db:"total" form:"total"`
	Status int `db:"status" form:"status"`
	BuyLimit int `db:"buy_limit" form:"buy_limit"`
	SoldLimitSecond int `db:"sold_limit_second" form:"sold_limit_second"`
	CreateTime int64 `db:"create_time" form:"-"`

	StartTimeStr string
	EndTimeStr string
	CreateTimeStr string
	StatusStr string
}

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

type ActivityModel struct {

}

func NewActivityModel() *ActivityModel {
	return &ActivityModel{}
}

func (this *ActivityModel) GetActivityList() (list []*Activity, err error) {
	sql := "select * from activity order by id desc"
	err = Db.Select(&list, sql)
	if err != nil {
		logs.Error("select activity_list from table failed! error: %v", err)
		return
	}

	for _, v := range list {
		start  := time.Unix(v.StartTime, 0)
		end    := time.Unix(v.EndTime, 0)
		create := time.Unix(v.CreateTime, 0)
		v.StartTimeStr  = start.Format("2006-01-02 15:04:05")
		v.EndTimeStr    = end.Format("2006-01-02 15:04:05")
		v.CreateTimeStr = create.Format("2006-01-02 15:04:05")

		now := time.Now().Unix()
		if now >= v.EndTime {
			v.Status = ActivityStatusExpire
			v.StatusStr = "已结束"
			continue
		}

		if v.Status == ActivityStatusDisable {
			v.StatusStr = "已禁用"
			continue
		}

		if now < v.StartTime {
			v.StatusStr = "未开始"
		} else {
			v.StatusStr = "活动进行中"
		}
	}
	return
}

func (this *ActivityModel) CreateNewActivity(a *Activity) (id int64, err error) {
	sql := "insert into activity(name,pid,start_time,end_time,buy_rate,total,status,buy_limit,sold_limit_second,create_time) values(?,?,?,?,?,?,?,?,?,?)"
	result, err := Db.Exec(sql, a.Name, a.ProductId, a.StartTime, a.EndTime, a.BuyRate, a.Total, a.Status, a.BuyLimit, a.SoldLimitSecond, a.CreateTime)
	if err != nil {
		logs.Error("insert data into activity failed! sql: %v, error: %v", sql, err)
		return
	}

	id, _ = result.LastInsertId()
	return
}

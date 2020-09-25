package models

import (
	"github.com/astaxie/beego/logs"
	"time"
)

const (
	ActivityStatusNormal  = 0    // 正常可用
	ActivityStatusDisable = 1    // 禁用
	ActivityStatusSoldOut = 2    // 售罄
	ActivityStatusExpire  = 3    // 过期或结束
)

type Activity struct {
	Id int `db:"id" form:"-" json:"id"`
	Name string `db:"name" form:"name"`
	ProductId int `db:"pid" form:"pid"`
	StartTime int64 `db:"start_time" form:"-"`
	EndTime int64 `db:"end_time" form:"-"`
	BuyRate float64 `db:"buy_rate" form:"buy_rate"`
	Total int `db:"total" form:"total"`
	Left int `db:"left" form:"-"`
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
			if v.Status != ActivityStatusExpire {
				v.Status = ActivityStatusExpire
				go this.UpdateActivityExpireStatusToDb(v)
				go SyncActivityStatusToEtcd(v)
			}
			v.StatusStr = "已结束"
			continue
		}

		if v.Status == ActivityStatusDisable {
			v.StatusStr = "已禁用"
			continue
		}
		if v.Status == ActivityStatusSoldOut {
			v.StatusStr = "商品已售罄"
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
	sql := "insert into activity(name,pid,start_time,end_time,buy_rate,total,`left`,status,buy_limit,sold_limit_second,create_time) values(?,?,?,?,?,?,?,?,?,?,?)"
	result, err := Db.Exec(sql, a.Name, a.ProductId, a.StartTime, a.EndTime, a.BuyRate, a.Total, a.Left, a.Status, a.BuyLimit, a.SoldLimitSecond, a.CreateTime)
	if err != nil {
		logs.Error("insert data into activity failed! sql: %v, error: %v", sql, err)
		return
	}

	id, _ = result.LastInsertId()
	return
}

// 将活动状态的变化更新到数据库
func (this *ActivityModel) UpdateActivityExpireStatusToDb(activity *Activity) (err error) {
	sql := "update activity set status=? where id=?"
	_, err = Db.Exec(sql, activity.Status, activity.Id)
	if err != nil {
		logs.Error("update column<status> for activity_table failed! activity_id: %d, status: %d, error: %v", activity.Id, activity.Status, err)
	} else {
		logs.Info("update column<status> for activity_table success! activity_id: %d, status: %d", activity.Id, activity.Status)
	}
	return
}

// 将etcd中获取到的活动数据更新到数据库 (主要更新left和status字段 即活动商品的剩余数量和活动的状态)
func updateActivityToDb(activityList *[]SecActivityConf) {
	now := time.Now().Unix()
	for _, v := range *activityList {
		if v.Status == ActivityStatusDisable || v.Status == ActivityStatusExpire || v.EndTime <= now {
			continue
		}
		sql := "update activity set `left`=?,status=? where id=?"

		_, err := Db.Exec(sql, v.Left, v.Status, v.ActivityId)
		if err != nil {
			logs.Error("update column<status> for activity_table failed! activity_id: %d, status: %d, error: %v", v.ActivityId, v.Status, err)
		} else {
			logs.Info("update column<status> for activity_table success! activity_id: %d, status: %d", v.ActivityId, v.Status)
		}
	}
}
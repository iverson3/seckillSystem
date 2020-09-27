package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"seckillsystem/secAdmin/models"
	"strings"
	"time"
)

type ActivityController struct {
	beego.Controller
}

func (this *ActivityController) ListActivity() {
	model := models.NewActivityModel()

	list, err := model.GetActivityList()
	if err != nil {
		logs.Warn("get activity list failed! error: %v", err)
	}
	logs.Info("activity list: %v", list)

	this.Data["activity_list"] = list
	this.Layout = "layout/layout.html"
	this.TplName = "activity/list.html"
}

func (this *ActivityController) UpdateActivityStatus() {
	result := make(map[string]interface{})
	defer func() {
		this.Data["json"] = result
		this.ServeJSON()
	}()
	result["code"] = 200
	result["msg"] = "success"

	id, err := this.GetInt("id")
	if err != nil {
		result["code"] = 501
		result["msg"] = "get params[id] failed"
		return
	}
	status, err := this.GetInt("status")
	if err != nil {
		result["code"] = 501
		result["msg"] = "get params[status] failed"
		return
	}

	activity := &models.Activity{
		Id:     id,
		Status: status,
		Left: -1,
	}
	model := models.NewActivityModel()
	err = model.UpdateActivityExpireStatusToDb(activity)
	if err != nil {
		result["code"] = 502
		result["msg"] = fmt.Sprintf("update activity status failed! error: %s", err.Error())
		return
	}
	go models.SyncActivityChangeToEtcd(activity)

	return
}

func (this *ActivityController) DeleteActivity() {
	result := make(map[string]interface{})
	defer func() {
		this.Data["json"] = result
		this.ServeJSON()
	}()
	result["code"] = 200
	result["msg"] = "success"

	activityId, err := this.GetInt("id")
	if err != nil {
		result["code"] = 501
		result["msg"] = "get params[id] failed"
		return
	}

	// 从数据库中删除活动
	model := models.NewActivityModel()
	err = model.DeleteActivityById(activityId)
	if err != nil {
		result["code"] = 502
		result["msg"] = fmt.Sprintf("delete activity from db failed! error: %v", err)
		return
	}

	// 从Etcd中删除活动
	etcdModel := models.NewEtcdModel()
	err = etcdModel.DeleteActivityFromEtcd(activityId)
	if err != nil {
		result["code"] = 502
		result["msg"] = fmt.Sprintf("delete activity from etcd failed! error: %v", err)
		return
	}
	return
}

// 获取活动商品的剩余数
func (this *ActivityController) GetActivityProductLeft() {
	result := make(map[string]interface{})
	defer func() {
		this.Data["json"] = result
		this.ServeJSON()
	}()
	result["code"] = 200
	result["msg"] = "success"

	activityId, err := this.GetInt("id")
	if err != nil {
		result["code"] = 501
		result["msg"] = "get params[id] failed"
		return
	}

	model := models.NewRedisModel()
	left, err := model.GetProductLeftNum(activityId)
	if err != nil {
		result["code"] = 502
		result["msg"] = fmt.Sprintf("get product left from redis failed! error: %v", err)
		return
	}

	activity := &models.Activity{
		Id:              activityId,
		Left:            left,
		Status:          -1,
	}
	activityModel := models.NewActivityModel()
	_ = activityModel.UpdateProductLeftNumToDb(activityId, left)

	go models.SyncActivityChangeToEtcd(activity)

	result["data"] = left
	return
}

func (this *ActivityController) CreateActivity() {
	if this.Ctx.Input.IsPost() {
		var err error
		defer func() {
			if err != nil {
				this.Data["ErrorTitle"] = "创建活动出错"
				this.Data["Error"] = err.Error()
				this.Layout = "layout/layout.html"
				this.TplName = "error/form_error.html"
			}
		}()

		activity := &models.Activity{}
		err = this.ParseForm(activity)
		if err != nil {
			return
		}
		activity.CreateTime = time.Now().Unix()
		activity.Left       = activity.Total

		startTime, err := dealTimeField(this.GetString("start_time"))
		if err != nil {
			return
		}
		activity.StartTime = startTime

		endTime, err := dealTimeField(this.GetString("end_time"))
		if err != nil {
			return
		}
		activity.EndTime = endTime

		err = validateForm(activity)
		if err != nil {
			return
		}

		model := models.NewActivityModel()
		newActivityId, err := model.CreateNewActivity(activity)
		if err != nil {
			return
		}
		// 得到新插入数据库记录的主键Id
		activity.Id = int(newActivityId)

		etcdModel := models.NewEtcdModel()
		err = etcdModel.SyncNewActivityToEtcd(activity)
		if err != nil {
			return
		}

		this.Redirect("/activity", 302)
		return
	}

	model := models.NewProductModel()
	list, err := model.GetProductList()
	if err != nil {
		logs.Error("get product list failed! error: %v", err)
		return
	}

	this.Data["product_list"] = list
	this.Layout = "layout/layout.html"
	this.TplName = "activity/create.html"
}

func dealTimeField(timeStr string) (timeUnix int64, err error) {
	// 时间转换的模板，golang里面只能是 "2006-01-02 15:04:05" （go的诞生时间）
	timeTemplate := "2006-01-02 15:04:05"

	splitTime := strings.Split(timeStr, "T")
	timeNewStr := splitTime[0] + " " + splitTime[1] + ":00"
	times, err := time.ParseInLocation(timeTemplate, timeNewStr, time.Local)
	if err != nil {
		return
	}
	timeUnix = times.Unix()
	return
}

func validateForm(activity *models.Activity) (err error) {
	if activity.Name == "" {
		err = fmt.Errorf("活动标题不能为空")
		return
	}
	if activity.ProductId == 0 {
		err = fmt.Errorf("活动商品Id不能为空")
		return
	}
	if activity.EndTime <= time.Now().Unix() {
		err = fmt.Errorf("活动的时间设置有误")
		return
	}
	if activity.StartTime >= activity.EndTime {
		err = fmt.Errorf("活动的结束时间必须大于开始时间")
		return
	}
	if !(activity.BuyRate > 0 && activity.BuyRate <= 1) {
		err = fmt.Errorf("秒杀概率必须大于0 且小于等于1")
		return
	}
	if activity.Total < 1 {
		err = fmt.Errorf("活动的商品数量设置有误")
		return
	}

	model := models.NewProductModel()
	exists, err, product := model.ProductExists(activity.ProductId)
	if err != nil {
		return
	}
	if exists == false {
		err = fmt.Errorf("活动的商品Id不存在")
		return
	}
	if activity.Total > product.Total {
		err = fmt.Errorf("活动商品的数量超过了该商品的总库存数量")
		return
	}
	return
}

package service

import "sync"

// 管理各个秒杀商品卖出的总数量
type ProductCountMgr struct {
	ProductCount map[int]int
	lock sync.RWMutex
}

// 创建一个管理活动商品的数量的计数器
func NewProductCountMgr() *ProductCountMgr {
	return &ProductCountMgr{
		ProductCount: make(map[int]int, 128),
	}
}

func (this *ProductCountMgr) Count(activityId int) int {
	this.lock.RLock()
	defer this.lock.RUnlock()

	return this.ProductCount[activityId]
}

func (this *ProductCountMgr) Increment(activityId, count int) {
	this.lock.Lock()
	defer this.lock.Unlock()

	sum, ok := this.ProductCount[activityId]
	if !ok {
		this.ProductCount[activityId] = count
	} else {
		this.ProductCount[activityId] = sum + count
	}
}
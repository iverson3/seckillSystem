package service

import "sync"

// 记录当前用户所有抢到的秒杀商品的数量
type UserBuyHistory struct {
	history map[int]int
	lock sync.RWMutex
}

func (this *UserBuyHistory) Increment(activityId, count int) {
	this.lock.Lock()
	defer this.lock.Unlock()

	num, ok := this.history[activityId]
	if !ok {
		this.history[activityId] = count
	} else {
		this.history[activityId] = num + count
	}
}

func (this *UserBuyHistory) GetProductBuyCount(activityId int) (count int) {
	this.lock.RLock()
	defer this.lock.RUnlock()

	count, ok := this.history[activityId]
	if !ok {
		count = 0
	}
	return
}


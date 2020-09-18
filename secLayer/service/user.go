package service

import "sync"

// 记录当前用户所有抢到的秒杀商品的数量
type UserBuyHistory struct {
	history map[int]int
	lock sync.RWMutex
}

func (this *UserBuyHistory) Increment(pid, count int) {
	this.lock.Lock()
	defer this.lock.Unlock()

	num, ok := this.history[pid]
	if !ok {
		this.history[pid] = count
	} else {
		this.history[pid] = num + count
	}
}

func (this *UserBuyHistory) GetProductBuyCount(pid int) (count int) {
	this.lock.RLock()
	defer this.lock.RUnlock()

	count, ok := this.history[pid]
	if !ok {
		count = 0
	}
	return
}


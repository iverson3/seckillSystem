package service

import "sync"

// 管理各个秒杀商品卖出的总数量
type ProductCountMgr struct {
	ProductCount map[int]int
	lock sync.RWMutex
}

func NewProductCountMgr() *ProductCountMgr {
	return &ProductCountMgr{
		ProductCount: make(map[int]int, 128),
	}
}

func (this *ProductCountMgr) Count(pid int) int {
	this.lock.RLock()
	defer this.lock.RUnlock()

	return this.ProductCount[pid]
}

func (this *ProductCountMgr) Increment(pid, count int) {
	this.lock.Lock()
	defer this.lock.Unlock()

	sum, ok := this.ProductCount[pid]
	if !ok {
		this.ProductCount[pid] = count
	} else {
		this.ProductCount[pid] = sum + count
	}
}
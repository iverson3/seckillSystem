package main

import (
	"fmt"
	"sync"
	"time"
)

// DCL   Double Check Lock

// 假设这个结构体对象在程序运行中只能有一个 即单例模式
type DbManager struct {
	driver string
	pool []int
}

var DbMgr *DbManager
var lock sync.Mutex

func main() {
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go worker(i, &wg)
	}

	wg.Wait()
}

func worker(n int, wg *sync.WaitGroup) {
	instance := getInstance(n)
	if instance == nil {
		fmt.Println("instance is nil")
	} else {
		fmt.Printf("instance's driver is %s \n", instance.driver)
	}

	wg.Done()
}

func getInstance(no int) *DbManager {
	// Double Check Lock
	// 才能保证多线程的情况下，只会产生一个实例
	if DbMgr == nil {
		lock.Lock()
		if DbMgr == nil {
			time.Sleep(1 * time.Millisecond)
			DbMgr = &DbManager{
				driver: fmt.Sprintf("mysql-%d", no),
				pool:   nil,
			}
		}
		lock.Unlock()
	}
	return DbMgr
}

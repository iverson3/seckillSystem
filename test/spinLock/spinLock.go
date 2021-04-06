package main

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

// 自旋锁

type PublicData struct {
	Count int
}

var LockEmpty = true
var myLock sync.Locker
var publicData *PublicData

var sm sync.Mutex     // 互斥锁
var smRW sync.RWMutex  // 读写锁

func main() {
	myLock = NewSpinLock()
	publicData = &PublicData{
		Count: 0,
	}
	loopTimes := 100000

	var wg sync.WaitGroup
	start := time.Now().UnixNano() / int64(time.Millisecond)
	for i := 0; i < 10000; i++ {
		//go Worker(loopTimes, i, &sm, &wg)
		//go WorkerWithLock(i)
		go WorkerWithLock2(loopTimes, i, &wg)
		wg.Add(1)
	}
	wg.Wait()
	end := time.Now().UnixNano() / int64(time.Millisecond)

	fmt.Printf("run time1: %d \n", end - start)
	fmt.Println(publicData.Count)
}

func Worker(loopTimes int, no int, sm *sync.Mutex, wg *sync.WaitGroup) {
	sm.Lock()
	for i := 0; i < loopTimes; i++ {
		publicData.Count++
	}
	sm.Unlock()
	wg.Done()
}

func WorkerWithLock(no int) {
	lock := Lock()
	if !lock {
		fmt.Printf("worker[%d] get lock failed! \n", no)
		return
	}
	fmt.Printf("worker[%d] get lock successed! \n", no)

	// 获取锁之后 操作共享数据
	fmt.Printf("worker[%d] prepare to operation public data! \n", no)
	for i := 0; i < 1000; i++ {
		publicData.Count++
	}
	fmt.Printf("worker[%d] operation public data is over! \n", no)

	UnLock()
	fmt.Printf("worker[%d] unlock successed! \n", no)
}

func WorkerWithLock2(loopTimes int, no int, wg *sync.WaitGroup) {
	myLock.Lock()
	// 获取锁之后 操作共享数据
	for i := 0; i < loopTimes; i++ {
		publicData.Count++
	}
	myLock.Unlock()
	wg.Done()
}


// 自己实现 (高并发下 会出现问题)
// 因为对标识变量的判断和对其的赋值 不是原子操作，而是两步操作；在高并发下 势必会有问题出现
func Lock() bool {
	//times := 0
	for !TryLock() {
		//time.Sleep(1 * time.Millisecond)
		//times++
		//if times >= 100 {
		//	return false
		//}
	}
	return true
}
func TryLock() bool {
	// 这里在判断和操作LockEmpty变量时，高并发下还是会出现问题
	// 必须确保下面的判断和赋值是原子操作 （可以使用atomic的原子操作对标识变量进行判断和修改 保证操作的原子性）
	if LockEmpty {
		LockEmpty = false
		return true
	}
	return false
}
func UnLock()  {
	LockEmpty = true
}


// 使用atomic实现
type spinLock uint32
func (sl *spinLock) Lock() {
	// 使用atomic的原子操作对标识变量进行判断和修改 保证操作的原子性
	for !atomic.CompareAndSwapUint32((*uint32)(sl), 0, 1) {
		runtime.Gosched()
	}
}
func (sl *spinLock) Unlock() {
	atomic.StoreUint32((*uint32)(sl), 0)
}
func NewSpinLock() sync.Locker {
	var lock spinLock
	return &lock
}
package main

import (
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"sync"
	"time"
)

// zookeeper分布式锁

func main() {
	hosts := []string{"47.107.149.234:2181", "47.107.149.234:2182", "47.107.149.234:2183"}
	zkConn, _, err := zk.Connect(hosts, 5 * time.Second)
	if err != nil {
		panic(err)
	}
	defer zkConn.Close()

	lockPath := "/test/zklock"
	mylock := zk.NewLock(zkConn, lockPath, zk.WorldACL(zk.PermAll))

	var wg sync.WaitGroup
	for i := 1; i < 10; i++ {
		wg.Add(1)
		go worker(&wg, i, mylock)
	}

	// unix.epollCreate1()

	wg.Wait()
}

func worker(wg *sync.WaitGroup, no int, mylock *zk.Lock) {
	defer wg.Done()

	// 加锁前的校验逻辑
	time.Sleep(1 * time.Second)

	// 加锁
	err := mylock.Lock()
	if err != nil {
		fmt.Printf("worker[%d] lock failed! \n", no)
		return
	}
	defer func() {
		// 释放锁
		err = mylock.Unlock()
		if err != nil {
			fmt.Printf("worker[%d] unlock failed! \n", no)
		} else {
			fmt.Printf("worker[%d] unlock successfully! \n", no)
		}
	}()

	fmt.Printf("worker[%d] lock successfully! \n", no)

	// lock success, do work
	fmt.Printf("worker[%d] start to work! \n", no)

	time.Sleep(2 * time.Second)

	// work is over
	fmt.Printf("worker[%d] end for work! \n", no)
}


//func lockExist(conn *zk.Conn, lockPath string) (err error, exists bool) {
//	exists, _, err = conn.Exists(lockPath)
//	if err != nil {
//		return
//	}
//	return nil, exists
//}

//func getLock(conn *zk.Conn, lockPath string) (err error, lock bool) {
//	acls := zk.WorldACL(zk.PermAll)
//	// 创建临时节点
//	_, err = conn.Create(lockPath, []byte(""), 0, acls)
//	if err != nil {
//		return
//	}
//	return nil, true
//}
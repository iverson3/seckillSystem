package main

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

// 实现一个简单的双向链表
type Node struct {
	no int
	status int
	locked int

	Next *Node // 指向下一个节点的地址
	Pre *Node  // 指向上一个节点的地址
}
type LinkList struct {
	atomicSign uint32
	headNode *Node // 头节点
	tailNode *Node // 尾节点
}

// 获取链表中节点的数量
func (this *LinkList) Length() int {
	cur := this.headNode
	sum := 0
	for cur != nil {
		sum++
		cur = cur.Next
	}
	return sum
}
func (this *LinkList) getNodeByNo(key int) (node *Node) {
	cur := this.headNode
	for cur != nil {
		if cur.no == key {
			node = cur
			break
		}
		cur = cur.Next
	}
	return
}
// 向链表尾部追加节点
func (this *LinkList) Append(key int) *Node {
	node := &Node{no: key}

	// 确保中间这段代码是原子性的
	for !atomic.CompareAndSwapUint32(&(this.atomicSign), 0, 1) {
		runtime.Gosched()
	}

	// 判断头节点为nil (即链表为空)
	if this.headNode == nil {
		node.locked = 1
		this.headNode = node
	} else {
		this.tailNode.Next = node
		if this.tailNode.locked == 1 && this.tailNode.status == 1 {
			this.tailNode.Next.locked = 1
		}
	}
	// 更新尾节点
	this.tailNode = node

	atomic.StoreUint32(&(this.atomicSign), 0)
	return node
}

// 向链表头部插入节点 (暂时不用)
//func (this *LinkList) Add(key int) *Node {
//	node := &Node{no: key}
//	node.Next = this.headNode
//	this.headNode = node
//
//	next := this.headNode.Next
//	if next != nil {
//		next.Pre = node
//	}
//	return node
//}



type spinLockWithQueue struct {
	queueList LinkList
}
func (this *spinLockWithQueue) Lock() {
	key := GetGoroutineId()
	node := this.queueList.Append(key)
	// 在本地变量上自旋
	for node.locked == 0 {
		runtime.Gosched()
	}
}
func (this *spinLockWithQueue) Unlock() {
	key := GetGoroutineId()
	node := this.queueList.getNodeByNo(key)
	if node != nil {
		node.status = 1
	}
	if node != nil && node.Next != nil {
		// 通知下一个节点结束自旋 可获取锁了
		node.Next.locked = 1
	}
}
func NewSpinLockWithQueue() *spinLockWithQueue {
	var lock spinLockWithQueue
	lock.queueList = LinkList{headNode: nil}
	return &lock
}
// 获取当前goroutine的协程id
func GetGoroutineId() int {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("panic recover:panic info:%v", err)
		}
	}()

	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	//fmt.Println(string(buf[:n]))
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		panic(fmt.Sprintf("cannot get goroutine id: %v", err))
	}
	return id
}



var myLock *spinLockWithQueue
type PublicData struct {
	Count int
}
var publicData *PublicData

func main() {
	publicData = &PublicData{Count: 0}
	myLock = NewSpinLockWithQueue()

	for i := 0; i < 10000; i++ {
		go Worker(i + 1)
	}
	time.Sleep(5 * time.Second)

	// 遍历链表
	//cur := myLock.queueList.headNode
	//for cur != nil {
	//	fmt.Printf("node no: %d; node locked: %d \n", cur.no, cur.locked)
	//	cur = cur.Next
	//}

	fmt.Println(publicData.Count)
	fmt.Println(myLock.queueList.Length())
}
func Worker(no int) {
	myLock.Lock()
	fmt.Printf("worker[%d] get lock successed! \n", no)

	// 获取自旋锁之后 操作共享数据
	fmt.Printf("worker[%d] prepare to operation public data! \n", no)
	for i := 0; i < 1000; i++ {
		publicData.Count++
	}
	fmt.Printf("worker[%d] operation public data is over! \n", no)

	myLock.Unlock()
	fmt.Printf("worker[%d] unlock successed! \n", no)
}
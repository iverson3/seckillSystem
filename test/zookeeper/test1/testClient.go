package main

import (
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"time"
)

func callback(event zk.Event) {
	fmt.Println(">>>>>>>>>>>>>>>>>>>")
	fmt.Println("path:", event.Path)
	fmt.Println("type:", event.Type.String())
	fmt.Println("state:", event.State.String())
	fmt.Println("<<<<<<<<<<<<<<<<<<<")
}

func main() {
	// 添加事件回调通知配置
	//options := zk.WithEventCallback(callback)

	hosts := []string{"47.107.149.234:2181"}
	conn, _, err := zk.Connect(hosts, 5 * time.Second)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// 参数
	path := "/test/t3"
	data := "hello,world"
	var flags int32 = 0
	acls := zk.WorldACL(zk.PermAll)  // permission

	// Create
	created, err := conn.Create(path, []byte(data), flags, acls)
	if err != nil {
		fmt.Printf("create path[%s] error: %v", path, err)
		return
	}
	fmt.Printf("create reasult: %s \n", created)

	// Exists
	exists, stat, err := conn.Exists(path)
	if err != nil {
		fmt.Printf("exist path[%s] error: %v", path, err)
		return
	}
	fmt.Printf("exist reasult: %v \n", exists)

	// Update
	newData := "new World!"
	stat, err = conn.Set(path, []byte(newData), stat.Version)
	if err != nil {
		fmt.Printf("update path[%s] error: %v", path, err)
		return
	}
	fmt.Printf("update path[%s] successfully \n", path)

	// Get
	res, stat, err := conn.Get(path)
	if err != nil {
		fmt.Printf("get path[%s] error: %v", path, err)
		return
	}
	fmt.Printf("get data: %s \n", res)

	// Delete
	err = conn.Delete(path, stat.Version)
	if err != nil {
		fmt.Printf("delete path[%s] error: %v", path, err)
		return
	}
	fmt.Printf("delete path[%s] successfully \n", path)
}

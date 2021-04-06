package main

import (
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"time"
)

func main() {
	hosts := []string{"47.107.149.234:2181", "47.107.149.234:2182", "47.107.149.234:2183"}
	conn, _, err := zk.Connect(hosts, 5 * time.Second)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// 参数
	rootPath := "/test/t5"
	data := "v111"
	var flags int32 = 0
	acls := zk.WorldACL(zk.PermAll)  // permission

	// 添加根节点
	exists, _, err := conn.Exists(rootPath)
	if err != nil {
		return
	}
	if !exists {
		created, err := conn.Create(rootPath, []byte(data), flags, acls)
		if err != nil {
			fmt.Printf("create path[%s] error: %v", rootPath, err)
			return
		}
		fmt.Printf("create rootPath reasult: %s \n", created)
	}

	// 添加子节点
	cur_time := time.Now().Unix()
	childPath := fmt.Sprintf("%s/child_%d", rootPath, cur_time)
	created, err := conn.Create(childPath, []byte("child_111111"), zk.FlagEphemeral, acls)
	if err != nil {
		fmt.Printf("create path[%s] error: %v", childPath, err)
		return
	}
	fmt.Printf("create childPath reasult: %s \n", created)

	// 监听指定节点
	children, stat, eventsChan, err := conn.ChildrenW(rootPath)
	if err != nil {
		fmt.Printf("watch path[%s] error: %v", rootPath, err)
		return
	}
	fmt.Printf("path[%s] stat: %v", rootPath, stat)

	// 遍历指定节点的所有子节点
	fmt.Printf("root_path[%s] child_count: %d\n", rootPath, len(children))
	for idx, ch := range children {
		fmt.Printf("%d, %s \n", idx, ch)
	}

	// 获取发生的监听事件
	select {
	case event := <- eventsChan:
		fmt.Println("path:", event.Path)
		fmt.Println("type:", event.Type.String())
		fmt.Println("state:", event.State.String())

		switch event.Type {
		case zk.EventNodeCreated:
			fmt.Printf("has new node[%s] be created \n", event.Path)
		case zk.EventNodeDataChanged:
			fmt.Printf("has node[%d] data be changed \n", event.Path)
		case zk.EventNodeDeleted:
			fmt.Printf("has node[%s] be deteted \n", event.Path)
		case zk.EventNodeChildrenChanged:
			fmt.Printf("Children node be changed, root node[%s] \n", event.Path)
		default:
			fmt.Printf("unknown event \n")
		}
	}
}

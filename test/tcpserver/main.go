package main

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:9999")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			logs.Error("accept failed! error: ", err)
			continue
		}

		go worker(conn)
	}
}

func worker(conn net.Conn) {
	defer conn.Close()

	var buffer []byte
	_, err := conn.Read(buffer)
	if err != nil {
		logs.Warn("read from client failed! error: ", err)
		return
	}

	logs.Info("got data from client: ", string(buffer))

	response := fmt.Sprintf("response from server: %s", string(buffer))
	_, err = conn.Write([]byte(response))
	if err != nil {
		logs.Warn("send to client failed! error: ", err)
	}

	logs.Info("send data to client successfully!")
}

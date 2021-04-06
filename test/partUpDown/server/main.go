package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"seckillsystem/test/partUpDown/common"
	"strings"
)

func main() {
	listener, err := net.Listen("tcp", ":9999")
	if err != nil {
		log.Fatal("listen failed! error: ", err)
		return
	}
	defer listener.Close()

	go WaitExitSignal()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("accept connect failed! error: ", err)
			break
		}

		log.Println("a new client connect to server")
		go ProcessConn(conn)
	}
}

func ProcessConn(conn net.Conn) {
	defer conn.Close()

	tf := &common.Transfer{
		Conn: conn,
	}

	closeConn := false
	for {
		log.Println("waiting for message from client...")
		mess, err := tf.ReadPkg()
		if err != nil {
			log.Println("read data from client failed! error: ", err)

			// 判断错误是否是客户端关闭了连接
			if err == io.EOF || strings.Contains(err.Error(), "close") {
				closeConn = true
				log.Println("Client has closed the connection")
			} else {
				continue
			}
		}

		if closeConn {
			log.Println("Server will close the connection")
			break
		}


		fmt.Printf("data from client: %v\n", mess)
	}
}

func WaitExitSignal() {
	ch := make(chan os.Signal)

	signal.Notify(ch, os.Interrupt, os.Kill)
	<-ch

	log.Println("got exit signal!")

	// 退出前的处理工作

	log.Println("process will exit!")
	os.Exit(1)
}

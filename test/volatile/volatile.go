package main

import (
	"fmt"
	"time"
)

var running bool = true

func main() {

	go func() {
		fmt.Println("start")
		for running {
		}
		fmt.Println("end")
	}()

	time.Sleep(time.Second)

	running = false

	time.Sleep(time.Second)
}

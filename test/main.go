package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

// 模拟大量用户同时请求秒杀系统

func main() {
	// 从命令行获取参数  go run main.go -n 1000 -c 100
	requestNum    := flag.Int("n", 1, "request num")
	concurrentNum := flag.Int("c", 1, "request concurrent num")
	flag.Parse()

	wg := &sync.WaitGroup{}
	chanInt := make(chan int, 100000)
	go Tongji(chanInt)

	// 每个并发需要发送的请求数
	requestsPerConcurrent := *requestNum / *concurrentNum
	activityId := 0
	start := time.Now()
	for i := 1; i <= *concurrentNum; i++ {
		userid := 1000 + i

		if i % 2 == 0 {
			activityId = 15
		} else {
			activityId = 15
		}
		go mockUserRequest(wg, userid, activityId, requestsPerConcurrent, chanInt)
		wg.Add(1)
	}

	wg.Wait()
	end := time.Now()
	seconds := end.Sub(start).Seconds()

	time.Sleep(3 * time.Second)
	close(chanInt)

	fmt.Printf("=========== total use seconds: %v ================= \n", seconds)
}

// 模拟用户请求
func mockUserRequest(wg *sync.WaitGroup, userid int, activityId int, requestNum int, in chan int) {
	defer wg.Done()
	for i := 1; i <= requestNum; i++ {
		url := fmt.Sprintf("http://localhost:8083/seckill?activity_id=%d&src=111&authcode=222&time=333&nance=444&userid=%d", activityId, userid + i * 10000)

		needContinue := false
		var data []byte
		func() {
			resp, err := http.Get(url)
			if err != nil || resp.StatusCode != http.StatusOK {
				fmt.Printf("request error: %v \n", err)
				in <- 4
				needContinue = true
				return
			}
			defer resp.Body.Close()

			data, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Printf("read data from request error: %v \n", err)
				in <- 4
				needContinue = true
				return
			}
		}()

		if needContinue {
			continue
		}

		//fmt.Printf("data: %v \n", string(data))
		if strings.Contains(string(data), "\"data\"") {
			if strings.Contains(string(data), "\\\"Code\\\":1000") {
				//fmt.Printf("request success; seckill success!\n")
				in <- 1
			} else {
				//fmt.Printf("request success; seckill failed!\n")
				in <- 2
			}
		} else {
			//fmt.Printf("user request failed!\n")
			in <- 3
		}
	}
}

func Tongji(out chan int) {
	seckillSucc := 0
	requestSuccess := 0
	requestFail := 0
	requestError := 0
	Sum := 0
	for res := range out {
		Sum++
		switch res {
		case 1:
			seckillSucc++
		case 2:
			requestSuccess++
		case 3:
			requestFail++
		case 4:
			requestError++
		}
	}

	fmt.Printf("Total Request: %d \n seckill success: %d \n seckill fail: %d \n failed Request: %d \n error Request: %d \n",
		Sum, seckillSucc, requestSuccess, requestFail, requestError)
}

package main

import (
	"fmt"
	"github.com/tealeg/xlsx"
	"runtime"
	"sync"
	"time"
)

var names []string
var scores []int
var arrLen int
var resCount int

var sheet *xlsx.Sheet
var err error
var wg sync.WaitGroup
var lock sync.Mutex

func main() {
	scores = []int{22, 22, 26, 26, 30, 30, 31, 31, 32, 33, 34, 36, 37, 37, 38, 38, 39, 40, 40, 40, 40,
		43, 43, 43, 43, 44, 44, 44, 45, 45, 46, 46, 48, 51, 52, 53, 53, 53, 53, 53, 53, 53, 54, 54, 54, 54}
	names = []string{"刘", "培", "刘", "芬", "范", "邵", "谢", "付", "葛", "明", "陈", "焱", "陈", "董", "焱",
		"戴", "范", "李", "涛", "明", "戴", "舒", "丹", "潘", "林", "邵", "伟", "潘", "葛", "秦", "张",
		"舒", "涛", "郑", "郑", "周", "董", "丹", "候", "爽", "秦", "胡", "付", "胡", "伟", "培"}

	arrLen = len(names)
	var reNameList []string
	var res [16]int
	var sum int

	//file := xlsx.NewFile()
	//sheet, err = file.AddSheet("Sheet1")
	//if err != nil {
	//	panic(err)
	//	return
	//}

	// 设置表头
	//row := sheet.AddRow()
	//row.SetHeightCM(1)
	//for i := 0; i < 16; i++ {
	//	cell := row.AddCell()
	//	cell.Value = "姓名:人数"
	//}

	go func() {
		for {
			fmt.Println(runtime.NumGoroutine())
			time.Sleep(100 * time.Millisecond)
			//runtime.GC()

		}
	}()

	start := time.Now().Unix()
	worker(0, 0, reNameList, res, sum, false)
	wg.Wait()
	end := time.Now().Unix()

	fmt.Println("结果总数: ", resCount)
	fmt.Printf("总耗时: %v 秒\n", end - start)

	//err = file.Save("result.xlsx")
	//if err != nil {
	//	panic(err)
	//} else {
	//	fmt.Println("save success")
	//}
}

func worker(n int, m int, reNameList []string, res [16]int, sum int, needDone bool) {
	for i := m; i < arrLen - 15 + n; i++ {
		if n - 1 + arrLen - i < 16 || sum >= 664 {
			break
		}
		
		if !in_arr(reNameList, names[i]){
			reNameList = append(reNameList, names[i])
			res[n] = i
			if n < 15 {
				if n == 0 {
					wg.Add(1)
					go worker(n + 1, i + 1, reNameList, res, sum + scores[i], true)
				} else {
					worker(n + 1, i + 1, reNameList, res, sum + scores[i], false)
				}
			} else {
				if sum + scores[i] == 664 {
					lock.Lock()
					resCount++
					lock.Unlock()
					//printArr(res)
					//outputExcel(res)
				}
			}
		}
	}
	if needDone {
		fmt.Println("goroutine over")
		wg.Done()
	}
}

func outputExcel(arr [16]int) {
	row := sheet.AddRow()
	row.SetHeightCM(1) //设置每行的高度
	for _, v := range arr {
		cell := row.AddCell()
		cell.Value = fmt.Sprintf("%s:%d", names[v], scores[v])
	}
}

func printArr(arr [16]int) {
	for _, v := range arr {
		fmt.Printf("%s:%d  ", names[v], scores[v])
	}
	fmt.Println()
}

func in_arr(arr []string, str string) bool {
	for _, v := range arr {
		if v == str {
			return true
		}
	}
	return false
}

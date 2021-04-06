package main

import "fmt"

var names []string
var scores []int

var total = 664
var totalLevel = 16

var resultCount = 0


func main() {
	scores = []int{22, 22, 26, 26, 30, 30, 31, 31, 32, 33, 34, 36, 37, 37, 38, 38, 39, 40, 40, 40, 40,
		43, 43, 43, 43, 44, 44, 44, 45, 45, 46, 46, 48, 51, 52, 53, 53, 53, 53, 53, 53, 53, 54, 54, 54, 54}

	names = []string{"刘", "培", "刘", "芬", "范", "邵", "谢", "付", "葛", "明", "陈", "焱", "陈", "董", "焱",
		"戴", "范", "李", "涛", "明", "戴", "舒", "丹", "潘", "林", "邵", "伟", "潘", "葛", "秦", "张",
		"舒", "涛", "郑", "郑", "周", "董", "丹", "候", "爽", "秦", "胡", "付", "胡", "伟", "培"}

	var hasScanIndex []int
	var hasScanName []string
	var hasScanScore []int
	worker(0, 0, 0, hasScanIndex, hasScanName, hasScanScore)
	fmt.Println(resultCount)
}

func worker(level int, n int, sum int, hasScanIndex []int, hasScanName []string, hasScanScore []int) {
	if sum > total {
		return
	}
	if (len(hasScanName) + len(names) - n) < totalLevel {
		return
	}
	if level == totalLevel || len(hasScanName) == totalLevel {
		if sum == total {
			resultCount++
			fmt.Printf("===== %d \n", sum)

			printIntArr(hasScanIndex)
			printStrArr(hasScanName)
			printIntArr(hasScanScore)
		}
		return
	} else {
		for i := len(names) - 1; i >= level; i-- {
			if is_in(hasScanIndex, i) {
				continue
			}
			if is_in_str(hasScanName, names[i]) {
				continue
			}
			if (level == totalLevel - 1) && (sum + scores[i] < total) {
				break
			}
			goWorker(i, level, sum, hasScanIndex, hasScanName, hasScanScore)
		}
	}
}

func goWorker(i int, level int, sum int, hasScanIndex []int, hasScanName []string, hasScanScore []int) {
	sum = sum + scores[i]
	hasScanIndex = append(hasScanIndex, i)
	hasScanName  = append(hasScanName, names[i])
	hasScanScore = append(hasScanScore, scores[i])

	worker(level + 1, i, sum, hasScanIndex, hasScanName, hasScanScore)
}

func is_in(arr []int, num int) bool {
	for _, v := range arr {
		if v == num {
			return true
		}
	}
	return false
}
func is_in_str(arr []string, str string) bool {
	for _, v := range arr {
		if v == str {
			return true
		}
	}
	return false
}

func printStrArr(arr []string) {
	for _, v := range arr {
		fmt.Printf(" %s ", v)
	}
	fmt.Println()
}
func printIntArr(arr []int) {
	for _, v := range arr {
		fmt.Printf(" %d ", v)
	}
	fmt.Println()
}
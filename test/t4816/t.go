package main
//
//import (
//	"fmt"
//	"strconv"
//)
//
//var names []string
//var scores []int
//
//var total = 48
//var totalLevel = 16
//var totalSum = 664
//
//func main() {
//	scores = []int{22, 22, 26, 26, 30, 30, 31, 31, 31, 32, 33, 34, 36, 37, 37, 38, 38, 39, 40, 40, 40, 40,
//		43, 43, 43, 43, 44, 44, 44, 45, 45, 46, 46, 48, 51, 52, 53, 53, 53, 53, 53, 53, 53, 53, 54, 54, 54, 54}
//
//	names = []string{"刘", "培", "刘", "芬", "范", "邵", "谢", "谢", "付", "葛", "明", "陈", "焱", "陈", "董", "焱",
//		"戴", "范", "李", "涛", "明", "戴", "舒", "丹", "潘", "林", "邵", "伟", "潘", "葛", "秦", "张",
//		"舒", "涛", "郑", "郑", "周", "董", "丹", "候", "周", "爽", "秦", "胡", "付", "胡", "伟", "培"}
//
//
//	var num uint64 = 0
//
//	var min uint64 = 1
//	var max uint64 = 1
//	for i := 0; i < totalLevel; i++ {
//		min = min * 2
//	}
//	max = min << 32
//	max--
//	min--
//	max = max << 48
//
//	maxStr := convertToBin(max)
//	maxStr = maxStr[0:total]
//	max = uint64(Str2DEC(maxStr))
//
//	//fmt.Println(min, max)
//	//fmt.Println(len(convertToBin(min)))
//	//fmt.Println(convertToBin(min))           // 低位 16个1
//	//fmt.Println(len(convertToBin(max)))
//	//fmt.Println(convertToBin(max))           // 高位 16个1
//
//
//	num = max
//	for num > min {
//
//		num--
//
//
//		lens := len(convertToBin(num))
//
//		breakLoop := false
//		sum := 0
//		tmpNum := num
//		var scanIndex []int
//		var scanNames []string
//		for i := 0; i < lens; i++ {
//			res := tmpNum & 0x01
//			if res == 1 {
//				if is_in_str(scanNames, names[i]) {
//					breakLoop = true
//					break
//				}
//				sum = sum + scores[i]
//				scanIndex = append(scanIndex, i)
//				scanNames = append(scanNames, names[i])
//			}
//			if len(scanIndex) > totalLevel {
//				breakLoop = true
//				break
//			}
//			if len(scanIndex) + lens - i < totalLevel {
//				breakLoop = true
//				break
//			}
//
//			tmpNum = tmpNum >> 1
//		}
//		if breakLoop {
//			continue
//		}
//		if len(scanIndex) == totalLevel && sum == totalSum {
//			fmt.Println(lens)
//			fmt.Println(convertToBin(num))
//
//			printIntArr(scanIndex)
//		}
//
//	}
//	fmt.Println("over")
//}
//
////
////
////func worker()  {
////	while(i!= 0) {
////		c+= i & 0x01;
////		i>>= 1;
////	}
////}
//
//// 将十进制数字转化为二进制字符串
//func convertToBin(num uint64) string {
//	s := ""
//	if num == 0 {
//		return "0"
//	}
//	// num /= 2 每次循环的时候 都将num除以2  再把结果赋值给 num
//	for ;num > 0; num /= 2 {
//		lsb := num % 2
//		// strconv.Itoa() 将数字强制性转化为字符串
//		s = strconv.Itoa(int(lsb)) + s
//	}
//	return s
//}
//
//func Str2DEC(s string) (num int) {
//	l := len(s)
//	for i := l - 1; i >= 0; i-- {
//		num += (int(s[l-i-1]) & 0xf) << uint8(i)
//	}
//	return
//}
//
//func is_in(arr []int, num int) bool {
//	for _, v := range arr {
//		if v == num {
//			return true
//		}
//	}
//	return false
//}
//func is_in_str(arr []string, str string) bool {
//	for _, v := range arr {
//		if v == str {
//			return true
//		}
//	}
//	return false
//}
//
//func printStrArr(arr []string) {
//	for _, v := range arr {
//		fmt.Printf(" %s ", v)
//	}
//	fmt.Println()
//}
//func printIntArr(arr []int) {
//	for _, v := range arr {
//		fmt.Printf(" %d ", v)
//	}
//	fmt.Println()
//}
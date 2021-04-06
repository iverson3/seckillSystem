package main

import "fmt"

var arr [5]int

func main() {
	arr[0] = 1
	arr[1] = 2
	arr[2] = 3
	arr[3] = 4
	arr[4] = 5

	s1 := arr[:2]

	fmt.Println(len(s1))
	fmt.Println(cap(s1))
	fmt.Println(&s1[0])
	fmt.Println(s1[1])

	s1 = append(s1, 987654)
	worker(s1)
	s1 = append(s1, 1000000000)
	s1 = append(s1, 2000000000)
	s1 = append(s1, 3000000000)

	fmt.Println(len(s1))
	fmt.Println(cap(s1))
	fmt.Println(&s1[0])
	fmt.Println(s1[1])

	fmt.Println("============================")

	var s2 []int
	s2 = append(s2, 111)
	fmt.Println(&s2[0])
	worker2(s2)
	fmt.Println(&s2[0])
}

func worker2(s []int) {
	fmt.Println(&s[0])
	s = append(s, 222)
	fmt.Println(&s[0])
}

func worker(s []int) {
	s[1] = 123456

	s = append(s, 666)
	s = append(s, 777)
	s = append(s, 888)
	s = append(s, 999)

	fmt.Printf("func slice len: %d \n", len(s))
	fmt.Printf("func slice cap: %d \n", cap(s))
	fmt.Printf("func slice addr: %v \n", &s[0])
	fmt.Printf("func slice s[1]: %d \n", s[1])
}
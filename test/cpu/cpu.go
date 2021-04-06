package main

import (
	"fmt"
	"sync"
)

var i, j, x, y int

func main() {
	n := 0

	for {
		n++
		i, j = 0, 0
		x, y = 0, 0

		var wg sync.WaitGroup
		wg.Add(2)

		go func(wg *sync.WaitGroup) {
			x = 1
			y = 2
			i = x
			wg.Done()
		}(&wg)
		go func(wg *sync.WaitGroup) {
			y = 1
			x = 2
			j = y
			wg.Done()
		}(&wg)

		//go func(wg *sync.WaitGroup) {
		//	//fmt.Printf("x === %d \n", x)
		//	y = 1
		//	j = x
		//	wg.Done()
		//}(&wg)

		wg.Wait()
		//fmt.Println("-----------------")

		if i == 0 || j == 0 {
			fmt.Printf("i = %d, x = %d \n", i, x)
			fmt.Printf("i = 0, n = %d break \n", n)
			break
		}

		//if i == 0 && j == 0 {
		//	fmt.Printf("i = %d, j = %d, x = %d, y = %d \n", i, j, x, y)
		//	fmt.Printf("i = 0, j = 0, n = %d break \n", n)
		//	break
		//}
		if n % 10000 == 0 {
			fmt.Printf("n ====================================== %d \n", n)
		}
		if n > 10000000 {
			fmt.Printf("n = %d, break \n", n)
			break
		}
	}
}

package main

import (
	"fmt"
	//"golang.org/x/sys/cpu"

	//"golang.org/x/sys/cpu"
	"sync"
	"time"
	_ "unsafe"
)

type T struct {
	i1, i2, i3, i4, i5, i6, i7 int64
	i int64
}

var s []T

func main() {
	//fmt.Println(unsafe.Sizeof(i))
	//cpu.CacheLinePad{}

	s := make([]T, 2)

	t1 := T{}
	t2 := T{}
	s = append(s, t1)
	s = append(s, t2)

	var wg sync.WaitGroup
	wg.Add(2)

	start := time.Now().UnixNano()
	go func(wg *sync.WaitGroup) {
		for i := 0; i < 10000000; i++ {
			s[0].i = int64(i)
		}
		wg.Done()
	}(&wg)

	go func(wg *sync.WaitGroup) {
		for j := 0; j < 10000000; j++ {
			s[1].i = int64(j)
		}
		wg.Done()
	}(&wg)

	wg.Wait()
	end := time.Now().UnixNano()

	fmt.Printf("time: %v", end - start)
}

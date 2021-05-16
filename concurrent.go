package main

import (
	"fmt"
	"runtime"
	"sync"
)

func RaceT1() {
	var count int64 = 0
	fmt.Println("go data race")
	cpuNums := runtime.NumCPU()
	// 使用waitGroup阻塞等待8个goroutiue完成任务
	var wg sync.WaitGroup
	wg.Add(cpuNums)
	for i := 0; i < cpuNums; i++ {
		// 开启8个goroutiue
		go func(nums int) {
			// 计数减 -1
			defer wg.Done()
			//每一个goroutine对计数加1000
			for j := 0; j < nums; j++ {
				count = count + 1
			}
		}(1000)
	}
	// 等待goroutine返回
	wg.Wait()
	fmt.Printf("RaceT1 count %d\n", count)
}

func RaceT2() {
	var count int64 = 0
	fmt.Println("go data race")
	// 互斥量
	var mu sync.Mutex
	proccess := runtime.GOMAXPROCS(0)
	cpuNums := runtime.NumCPU()
	fmt.Printf("proccess: %d,cpu nums:%d \n", proccess, cpuNums)
	// 使用waitGroup阻塞等待8个goroutiue完成任务
	var wg sync.WaitGroup
	wg.Add(cpuNums)
	for i := 0; i < cpuNums; i++ {
		// 开启8个goroutiue
		go func(nums int) {
			// 计数减 -1
			defer wg.Done()
			//每一个goroutine对计数加100
			for j := 0; j < nums; j++ {
				mu.Lock()
				count = count + 1
				mu.Unlock()
			}
		}(1000)
	}
	// 等待goroutine返回
	wg.Wait()
	fmt.Printf("RaceT2 count %d\n", count)
}

// 封装 使用匿名 mutex
type Counter struct {
	Count int64
	sync.Mutex
}

func RaceT3() {
	var c Counter = Counter{}
	var cpuNums int = runtime.NumCPU()
	fmt.Printf("cpu nums:%d \n", cpuNums)
	var wg sync.WaitGroup
	wg.Add(cpuNums)
	for i := 0; i < cpuNums; i++ {
		go func(nums int) {
			// wg 计数减一
			defer wg.Done()
			for j := 0; j < nums; j++ {
				c.Lock()
				c.Count++
				c.Unlock()
			}
		}(1000)
	}
	// 阻塞等待所有的子goroutine返回
	wg.Wait()
	fmt.Printf("RaceT3 count %d\n", c.Count)

}
func main() {
	// 数据竞争
	RaceT1()
	// 使用互斥量解决临界区 data race 问题
	RaceT2()
	RaceT3()
}

/*
运行时 发现 data race
PS D:\awesomeProject> go run -race .\concurrent.go
go data race
==================
WARNING: DATA RACE
Read at 0x000000692e90 by goroutine 8:
  main.main.func1()
      D:/awesomeProject/concurrent.go:29 +0x84

Previous write at 0x000000692e90 by goroutine 7:
  main.main.func1()
      D:/awesomeProject/concurrent.go:29 +0xa4

Goroutine 8 (running) created at:
  main.main()
      D:/awesomeProject/concurrent.go:23 +0x146

编译时插入指令 go race Detector工具能够检测出data race
go tool compile -race -S concurrent.go
*/

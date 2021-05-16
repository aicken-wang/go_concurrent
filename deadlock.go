package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	// 死锁
	var mTom sync.Mutex
	var mJack sync.Mutex
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		mTom.Lock()
		defer mTom.Unlock()
		// sleep
		time.Sleep(time.Millisecond*1000)
		mJack.Lock()
		defer mJack.Unlock()
	}()
	go func() {
		defer wg.Done()
		mJack.Lock()
		defer mJack.Unlock()
		time.Sleep(time.Millisecond*200)
		mTom.Lock()
		defer mTom.Unlock()
	}()
	wg.Wait()
	fmt.Println("执行完毕")
}

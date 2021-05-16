package main

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

// TestCond cond示例

func TestCond(nums int){
	c := sync.NewCond(&sync.Mutex{})
	var ready int
	for i := 0; i < nums; i++ {
		go func(idx int){
			time.Sleep(time.Duration(rand.Int63n(10)) * time.Second)
			// locker
			c.L.Lock()
			ready++
			c.L.Unlock()
			log.Printf("go routine %d 准备就绪\n",idx)
			// 广播唤醒所有的等待者
			c.Broadcast()
		}(i)
	}
	// 使用前加锁
	c.L.Lock()
	for ready != nums {
		//条件不满足 wait
		c.Wait()
		log.Println("唤醒一次，但条件不满足！")
	}
	c.L.Unlock()
	// 就绪
	log.Printf("ready: %d,条件满足所有go routine就绪", ready)
}
func main(){
	TestCond(100)
}

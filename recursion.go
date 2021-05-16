package main

import (
	"fmt"
	"github.com/petermattis/goid"
	"sync"
	"sync/atomic"
)
/*
实现一个可重的Mutex
1、标准包sync的mutex是不可重入的锁
2、设计一个可重入的锁，可以解决代码中或递归函数中重复加锁导致阻塞/死锁问题，
若要设计一个可重入的锁，必须记录goroutine的锁状态即只有上锁的goroutine才能解锁，通过goroutine id来标志
*/
// RecursionMutex 包装mutex实现一个可重入锁
type RecursionMutex struct{
	sync.Mutex // 处理临界区互斥
	owner int64 // 当前持有锁的goroutine id
	recursion int32 // 当前goroutine重入次数
}
// 实现Lock
func (m *RecursionMutex)Lock() {
	gid := goid.Get()
	fmt.Println("gid ",gid)
	// 根据 m.owner来判断当前goroutine是否为重入
	if atomic.LoadInt64(&m.owner) == gid {
		// recursion ++
		m.recursion++
		fmt.Printf("gid:%d return.\n",gid)
		return
	}
	// 第一次进入临界区 m.Lock() 加锁
	m.Mutex.Lock()
	// 记录当前 goroutine id
	atomic.StoreInt64(&m.owner,gid)
	fmt.Printf("gid:%d init lock.\n",gid)
	m.recursion = 1
}
// 实现 Unlock()
func(m *RecursionMutex)Unlock(){
	// 获取当前goroutine id
	gid := goid.Get()
	if atomic.LoadInt64(&m.owner) != gid{
		panic(fmt.Sprintf("wrong this Recursion owner:(%d),cur gid:(%d)",m.owner,gid))
	}
	// 释放一次计数并判断判断当前 recursion 是否为最后一次
	m.recursion--
	if m.recursion != 0 {
		// 不是最后一次直接返回
		fmt.Printf("gid:%d not last unlock.\n",gid)
		return
	}
	// 释放临界区资源,将当前recursionMutex的 owner记为 -1 标记失效
	atomic.StoreInt64(&m.owner,-1)
	m.Mutex.Unlock()
	fmt.Printf("gid:%d destory mutex.\n",gid)
}
func main() {
	// go get -u github.com/petermattis/goid
	var wg sync.WaitGroup
	rmu := RecursionMutex{}
	wg.Add(5)
	for i:=0; i < 5; i++{
		go func() {
			defer wg.Done()
			rmu.Lock()
			rmu.Unlock()
			gid := goid.Get()
			fmt.Println("goroutine gid ",gid)
		}()
	}
	wg.Wait()
	fmt.Println("执行完毕")
}

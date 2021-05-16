package main

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)
// mutex 定义常量
const (
	mutexLocked =  1 << iota // 加锁标识位置 第1个bit位
	mutexWoken // 唤醒标识位置 第2个bit位
	mutexStarving // 锁饥饿标识位置 第3个bit位
	mutexWaiterShift = iota // 标识waiter的起始bit位置 1<<(32-3) - 1
)
// Mutex 扩展一个Mutex结构，其他语言是继承，go使用组合
type Mutex struct {
	sync.Mutex
}

// TryLock 尝试获取锁
func (m *Mutex) TryLock() bool{
 	// 如果成功抢到锁 return true
	if atomic.CompareAndSwapInt32((*int32)(unsafe.Pointer(&m.Mutex)),0,mutexLocked) {
		return true
	}
	// 若处于唤醒，持锁，饥饿 这三个状态之一，这次请求就不参与竞争，直接返回false
	old := atomic.LoadInt32((*int32)(unsafe.Pointer(&m.Mutex)))
	if old&(mutexLocked|mutexWoken|mutexStarving) != 0 {
		return false
	}
	// 尝试在竞争的状态下请求锁,不一定能请求成功
	new := old|mutexLocked
	return atomic.CompareAndSwapInt32((*int32)(unsafe.Pointer(&m.Mutex)),old,new)

}
// Count 获取持有锁和等待获取锁的个数
func(m *Mutex)Count() int{
	/*
	 TODO:
	 为什么不封装为一个func是应该多个goroutine是并发操作这个mutex的状态，
	 直接调用就是获取当前时间锁的状态。所以不要为了简化代码而封装为一个函数调用
	 */

	// 获取state字段的值
	v:= atomic.LoadInt32((*int32)(unsafe.Pointer(&m.Mutex)))
	//fmt.Printf("v:%d\n",v)
	// 获取等待着的数量
	waiter := v >> mutexWaiterShift
	//fmt.Printf("v:%d\n",v)
	// 加上锁持有者的数量 0 或 1
	v = waiter + (v & mutexLocked)
	//fmt.Printf("v:%d\n",v)
	return int(v)
}
// IsLocked 是否被持有
func (m *Mutex) IsLocked() bool {
	// 获取锁的状态
	state := atomic.LoadInt32((*int32)(unsafe.Pointer(&m.Mutex)))
	// 是否持有锁
	return state&mutexLocked == mutexLocked
}
// IsWoken 是否有等待者被唤醒
func(m *Mutex) IsWoken() bool {
	// 获取锁的状态
	state := atomic.LoadInt32((*int32)(unsafe.Pointer(&m.Mutex)))
	return state&mutexWoken == mutexWoken
}
// IsStarving 锁是否处于饥饿状态
func(m *Mutex) IsStarving()bool {
	state := atomic.LoadInt32((*int32)(unsafe.Pointer(&m.Mutex)))
	return state&mutexStarving == mutexStarving
}
// Test1000 测试
func Test1000(){
	var mu Mutex
	for i:=0;i<1000;i++{
		go func() {
			mu.Lock()
			// 时间要设置大一点
			time.Sleep(time.Millisecond*2000) // sleep 2s
			mu.Unlock()
		}()
	}
	time.Sleep(time.Second) // sleep 1s
	// print
	fmt.Printf("wrong log waiter:%d, isLocked:%t, woken:%t,starving:%t \n",mu.Count(),mu.IsLocked(),mu.IsWoken(),mu.IsStarving())
}
// TestWG5 测试waiter
func TestWG5(){
	const num int = 10
	var m Mutex
	var wg sync.WaitGroup
	wg.Add(num)
	for i:=0;i < num; i++{
		go func(){
			defer wg.Done()
			fmt.Println("当前持有和等待锁的数：",m.Count())
			m.Lock() // 阻塞
			//下面的代码执行不了

		}()
	}
	wg.Wait()
	fmt.Println(m.Count())
}
// try Test Demo
func try(){
	var mu Mutex
	// 启动一个goroutine持有锁,随机sleep
	go func() {
		mu.Lock()
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		mu.Unlock()
	}()
	time.Sleep(50 * time.Millisecond)
	ok := mu.TryLock() // 尝试获取锁
	if ok {
		//获取锁成功
		fmt.Println("I get the mutex lock")
		// 释放
		mu.Unlock()
		return
	}
	fmt.Println("can't get the mutex lock")
}
// Test5Try Test
func Test5Try(){
	for i:=0;i<5;i++ {
		try()
	}
}
func main() {
	Test5Try()
	Test1000()
	TestWG5()
}

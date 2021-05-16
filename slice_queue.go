package main

import (
	"fmt"
	"sync"
)

// SliceQueue 基于切片的队列元素
type SliceQueue struct {
	data[] interface{} // any
	mu sync.Mutex // safe data race
}
// NewSliceQueue 构造
func NewSliceQueue(n int) (q *SliceQueue) {
	return &SliceQueue{data:make([]interface{},0,n)}
}
// EnQueue 入队
func (q *SliceQueue) EnQueue(v interface{}) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.data = append(q.data,v)
}
// DeQueue 出队
func (q *SliceQueue) DeQueue() interface{} {
	q.mu.Lock()
	defer q.mu.Unlock()

	if 0 == len(q.data) {
		return nil
	}
	v := q.data[0]
	q.data = q.data[1:]
	return v
}
func(q *SliceQueue) Print(){
	q.mu.Lock()
	defer q.mu.Unlock()
	for i,v:= range q.data {
		fmt.Printf("index:%d, value:%v \n",i,v)
	}
	fmt.Println()
}
func main() {
	queue:= NewSliceQueue(10)
	queue.EnQueue(1)
	queue.EnQueue(2)
	queue.Print()
	queue.DeQueue()
	queue.Print()
}

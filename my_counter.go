package main

import (
	"fmt"
	"sync"
)

type MyCounter struct{
	sync.Mutex
	Count int64
}
func main() {
	c := MyCounter{Count: 0}
	c.Lock()
	defer c.Unlock()
	c.Count++
	Incr(c)
	fmt.Println("hello mutex")
}
func Incr(c MyCounter) {
	c.Lock()
	c.Count++
	c.Unlock()
}

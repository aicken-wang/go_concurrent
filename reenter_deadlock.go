package main

import (
	"fmt"
	"sync"
)

func foo(l sync.Locker) {
	l.Lock()
	defer l.Unlock()
	bar(l)
	fmt.Println("I am foo.")

}
func bar(l sync.Locker){
	l.Lock()
	defer l.Unlock()
	fmt.Println("I am bar。")
}

func main() {
	m:= &sync.Mutex{}
	foo(m)
}

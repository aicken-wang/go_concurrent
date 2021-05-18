package main

import (
	"fmt"
	"time"
)
const NUM =4
func taskSchedule(){
	chanArr := make([]chan int, NUM)
	for i:=0; i<NUM; i++{
		chanArr[i] = make(chan int,1)
	}
	chanArr[0] <- 1
	for i:=0;i<NUM;i++ {
		idx := (i+1)%NUM
		go func(cur,next chan int,idx int) {
			for  {
				<-cur
				time.Sleep(time.Second)
				fmt.Println(idx+1)
				next <- 1
			}
		}(chanArr[i],chanArr[idx],i)
	}
	select {
		//阻塞
	}
}
func main(){
	fmt.Printf("start nums  = %d task\n",NUM)
	taskSchedule()
}

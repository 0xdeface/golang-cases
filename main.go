package main

import (
	"fmt"
	"time"
)

func main() {
	var start = make(chan struct{})
	for i := 0; i < 10; i++ {
		go func() {
			<-start
			fmt.Println("start")
		}()
	}
	close(start)
	time.Sleep(10 * time.Second)
}

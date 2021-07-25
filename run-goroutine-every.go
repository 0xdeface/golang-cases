package main

import (
	"fmt"
	"time"
)

func main() {
	everyThreeSecond := time.NewTicker(3 * time.Second)
	for { // if you remove for it's work once like a timer
		select {
		case <-everyThreeSecond.C:
			fmt.Println("there is every 3 seconds after run")
		}
	}
}

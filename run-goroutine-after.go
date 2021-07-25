package main

import (
	"fmt"
	"time"
)

func main() {
	afterTenSecond := time.NewTimer(10 * time.Second)
	select {
		case <- afterTenSecond.C:
			fmt.Println("there is ten seconds after run")
	}
}

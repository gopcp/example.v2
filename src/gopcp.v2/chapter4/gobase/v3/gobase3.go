package main

import (
	"fmt"
)

func main() {

	sync := make(chan int)
	name := "Eric"
	go func() {
		fmt.Printf("Hello, %s!\n", name)
		sync <- 5
	}()
	name = "Harry"
	<-sync
}

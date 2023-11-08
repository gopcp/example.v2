package main

import (
	"fmt"
	"time"
)

func main() {
	sendingInterval := time.Second
	receptionInterval := time.Second * 2
	intChan := make(chan int)
	go func() {
		for i := 1; i <= 5; i++ {
			intChan <- i
			fmt.Println("send[", i, "]:   ", time.Now())
			time.Sleep(sendingInterval)
		}
		close(intChan)
	}()
Loop:
	for {
		select {
		case v, ok := <-intChan:

			if !ok {
				break Loop
			}
			fmt.Println("receive[", v, "]:", time.Now())
		}
		time.Sleep(receptionInterval)
	}
	fmt.Println("End.")
}

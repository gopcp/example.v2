package main

import (
	"time"
)

func main() {
	go println("Go! Goroutine!")
	time.Sleep(time.Millisecond)
	// 前一条语句也可用runtime.Gosched()替换。
}

package main

import (
	"fmt"
	"time"
)

var strChan = make(chan string, 3)

var goChan = make(chan string)
var go2Chan = make(chan string)
var ptr bool = true

func main() {
	syncChan1 := make(chan struct{}, 1)
	syncChan2 := make(chan struct{}, 2)
	go func() { // 用于演示接收操作。
		i := 0
		<-syncChan1
		fmt.Println("Received a sync signal and wait a second... [receiver]")
		//time.Sleep(time.Second)
		for {
			//time.Sleep(time.Second)
			if elem, ok := <-strChan; ok {
				fmt.Println("Received:", elem, "[receiver]")
				i++
				if i == 3 {
					i = 0
					goChan <- "GO"
					<-go2Chan
				}

			} else {
				break
			}
		}
		fmt.Println("Stopped. [receiver]")
		syncChan2 <- struct{}{}
	}()
	go func() { // 用于演示发送操作。
		i := 0
		for _, elem := range []string{"a", "b", "c", "d", "e", "f", "g", "h"} {
			time.Sleep(time.Second)
			strChan <- elem
			fmt.Println("Sent:", elem, "[sender]")
			i++
			if i == 3 {
				if ptr {
					syncChan1 <- struct{}{}
					ptr = false
				}
				fmt.Println("Sent a sync signal. [sender]")
				i = 0

				<-goChan
				go2Chan <- "GO"
			}
		}
		fmt.Println("Wait 2 seconds... [sender]")
		//time.Sleep(time.Second * 2)
		close(strChan)
		syncChan2 <- struct{}{}
	}()
	<-syncChan2
	<-syncChan2
}

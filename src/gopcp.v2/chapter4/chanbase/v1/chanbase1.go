package main

import (
	"fmt"
)

var strChan = make(chan string, 3)

var goChan = make(chan string)
var flag = make(chan string, 1)
var ptr bool = false

func main() {
	syncChan1 := make(chan struct{}, 1)
	syncChan2 := make(chan struct{}, 2)
	letterSlice := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i"}
	flag <- "S"
	go func() { // 用于演示接收操作。
		<-syncChan1
		i := 0
		count := 0
		fmt.Println("Receiver Begin... [receiver]")
		for {

			fValue := <-flag
			//time.Sleep(time.Second)
			if fValue == "R" {
				if elem, ok := <-strChan; ok {
					fmt.Println("Received:", elem, "[receiver]")
					i++
					count++
					if count == len(letterSlice) {
						break
					} else if i == 3 {
						i = 0
						flag <- "S"
						goChan <- "GO"
					} else {
						flag <- "R"
					}

				} else {
					break
				}
			} else {
				flag <- "S"
			}
		}
		fmt.Println("Stopped. [receiver]")
		syncChan2 <- struct{}{}
	}()
	go func() { // 用于演示发送操作。
		i := 0

		fmt.Println("Send Begin... [sender]")
		for index, elem := range letterSlice {
			//time.Sleep(time.Second)

			fValue := <-flag

			if fValue == "S" {
				strChan <- elem
				fmt.Println("Sent:", elem, "[sender]")
				i++

				if index == (len(letterSlice) - 1) {
					flag <- "R"
				} else if i == 3 {
					fmt.Println("Sent a sync signal. [sender]")
					i = 0
					if !ptr {
						syncChan1 <- struct{}{}
						ptr = true
					}
					flag <- "R"
					<-goChan
				} else {
					flag <- "S"
				}
			} else {
				flag <- "R"
			}
		}
		close(strChan)
		syncChan2 <- struct{}{}
		fmt.Println("Stopped. [sender]")
	}()
	<-syncChan2
	<-syncChan2
}

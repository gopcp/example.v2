package main

import (
	"fmt"
)

// Counter 代表计数器的类型。
type Counter struct {
	pcount int
}

var mapChan = make(chan map[string]*Counter, 1)

func main() {
	waitChan := make(chan struct{}, 2)
	next := make(chan struct{})

	go func() { // 用于演示接收操作。
		for {

			if elem, ok := <-mapChan; ok {

				stop, ok := <-next
				if !ok {
					break
				}

				counter := elem["pcount"]
				counter.pcount++

				next <- stop

			} else {
				break
			}

		}
		fmt.Println("Stopped. [receiver]")
		waitChan <- struct{}{}
	}()
	go func() { // 用于演示发送操作。
		countMap := map[string]*Counter{
			"pcount": &Counter{},
		}
		for i := 0; i < 5; i++ {

			stop := <-next

			mapChan <- countMap
			fmt.Printf("The count map: %v. [sender]\n", countMap)

			if i == 4 {
				close(next)
				break
			}
			next <- stop
		}
		close(mapChan)
		waitChan <- struct{}{}
	}()

	next <- struct{}{}

	<-waitChan
	<-waitChan
}

func (counter *Counter) String() string {

	return fmt.Sprintf("{count: %d}", counter.pcount)

}

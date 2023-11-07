package main

func main() {

	chanSync := make(chan int)
	go func() {
		println("Go! Goroutine!")
		<-chanSync
	}()

	chanSync <- 9
}

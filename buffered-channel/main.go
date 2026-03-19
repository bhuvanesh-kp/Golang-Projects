package main

import (
	"fmt"
	"time"
)

func ListenAndSend(ch <-chan int) {
	for {
		i := <-ch
		fmt.Println("Message from channel: ", i)

		time.Sleep(2 * time.Second)
	}
}

func main() {
	ch := make(chan int, 90)

	go ListenAndSend(ch)

	for i:=0;i<100;i++{
		fmt.Println("Message send from the channel -> ")
		ch <- i
	}

	fmt.Println("channels are done with job")
	close(ch)
}
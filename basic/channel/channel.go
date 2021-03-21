package main

import (
	"fmt"
	"time"
)

func main() {
	// channelDeadlock() // this deadlock
	avoidChannelDeadlock()

	// channelDeadlock2() // this deadlock
	avoidChannelDeadlock2()

	// bufferecChannel doesn't deadlock
	bufferedChannel()

	channelSelectBlocking()

	channelSelectNonBlocking()
}

// this deadlock
func channelDeadlock() {

	// initialize the channel without size
	ch := make(chan bool)

	ch <- true

	fmt.Println(<-ch)
}

func avoidChannelDeadlock() {
	ch := make(chan bool)

	// send only
	go func(ch chan<- bool) {
		time.Sleep(1 * time.Second)
		ch <- true
	}(ch)

	fmt.Printf("avoidChannelDeadlock, %v\n", <-ch)
}

// this also deadlock
func channelDeadlock2() {

	// initialize the channel with buffer size 0
	ch := make(chan bool, 0)

	ch <- true

	fmt.Println(<-ch)
}

func avoidChannelDeadlock2() {
	ch := make(chan bool, 0)

	// send only
	go func(ch chan<- bool) {
		time.Sleep(1 * time.Second)
		ch <- true
	}(ch)

	fmt.Printf("avoidChannelDeadlock2, %v\n", <-ch)
}

func bufferedChannel() {

	// initialize the channel with buffer size 1, buffered channel
	ch := make(chan bool, 1)

	ch <- true

	fmt.Println(<-ch)
}

// blocking
// without default case, it is blocking and waiting for a channel to receive
func channelSelectBlocking() {
	quit := make(chan bool)
	go func() {
		time.Sleep(3 * time.Second)
		quit <- true
	}()

	ch := make(chan bool)
	go func() {
		for i := 0; i < 10; i++ {
			ch <- true
		}
	}()

	for {
		select {
		case b := <-ch:
			fmt.Printf("channelSelectBlocking, %v\n", b)
		case q := <-quit:
			fmt.Printf("channelSelectBlocking quit, %v\n", q)
			return
		}
	}
}

// non-blocking
// with default case, it is non-blocking
func channelSelectNonBlocking() {
	quit := make(chan bool)
	go func() {
		time.Sleep(5 * time.Second)
		quit <- true
	}()

	ch := make(chan bool)
	go func() {
		time.Sleep(2 * time.Second)
		for i := 0; i < 10; i++ {
			ch <- true
			time.Sleep(500 * time.Millisecond)
		}
	}()

	for {
		select {
		case b := <-ch:
			fmt.Printf("channelSelectNonBlocking, %v\n", b)
		case q := <-quit:
			fmt.Printf("channelSelectNonBlocking quit, %v\n", q)
			return
		default:
			fmt.Printf("channelSelectNonBlocking, default case\n")
			time.Sleep(100 * time.Millisecond)
		}
	}
}

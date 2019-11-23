package main

import (
	"fmt"
	"time"
)

func main() {
	s := make([]byte, 1000)

	copy(s, "hello world")
	ch := make(chan bool)

	go func(sp []byte) {
		fmt.Println("goroutine1", string(sp))
		time.Sleep(10 * time.Millisecond)
		fmt.Println("goroutine2", string(sp))
		ch <- true
	}(s[0:len("hello world")])

	time.Sleep(50 * time.Microsecond)
	copy(s, "world hello")

	<-ch
}

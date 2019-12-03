package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var mu sync.Mutex
	go func() {
		mu.Lock()
		time.Sleep(3 * time.Second)

		fmt.Println("Unlock 2 start...")
		mu.Unlock()
		fmt.Println("Unlock 2 end...")
	}()
	time.Sleep(time.Second)

	fmt.Println("Unlock 1 start...")
	mu.Unlock()
	fmt.Println("Unlock 1 end...")

	fmt.Println("waiting...")
	select {}
}

package main

import (
	"fmt"
	"os"
	"runtime"
	"time"
)

const (
	MB = 1024 * 1024
)

func main() {
	blocks := make([][MB]byte, 0)
	fmt.Println("Child pid is", os.Getpid())

	for {
		blocks = append(blocks, [MB]byte{})
		printMemUsage()
		time.Sleep(time.Second)
	}
}

func printMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc %d MiB\tSys %d MiB\n", m.Alloc/MB, m.Sys/MB)
}

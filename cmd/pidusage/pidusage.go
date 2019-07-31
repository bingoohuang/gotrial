package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/struCoder/pidusage"
)

var (
	interval int64
)

func init() {
	flag.Int64Var(&interval, "d", 15, "ticker durations(seconds)")

	flag.Parse()
}

func main() {
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	quit := make(chan bool)

	OnInterrupt(func() {
		quit <- true
	})

	go func() {
		for {
			select {
			case <-ticker.C:
				sysInfo, _ := pidusage.GetStat(os.Getpid())

				log.Printf("cpu: %.2f, mem: %.2f\n", sysInfo.CPU, sysInfo.Memory/1024/1024)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	select {}

}

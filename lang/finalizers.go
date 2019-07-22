package main

import (
	"flag"
	"fmt"
	"math/rand"
	"runtime"
	"runtime/debug"
	"strconv"
	"time"
)

/*
Go: Finalizers
https://medium.com/@blanchon.vincent/go-finalizers-786df8e17687

➜  gistgo go run finalizers.go
Allocation: 0.121944 Mb, Number of allocation: 152
Allocation: 31.139055 Mb, Number of allocation: 2390094
Allocation: 110.090820 Mb, Number of allocation: 4472773
Allocation: 0.136576 Mb, Number of allocation: 198
➜  gistgo go run finalizers.go  -c
Allocation: 0.122190 Mb, Number of allocation: 155
Allocation: 18.161726 Mb, Number of allocation: 1390097
Allocation: 0.126434 Mb, Number of allocation: 177
Allocation: 0.127913 Mb, Number of allocation: 181

理解要点：
When the garbage collector finds an unreachable block with an associated finalizer,
it clears the association and runs finalizer(obj) in a separate goroutine.
This makes obj reachable again, but now without an associated finalizer.
Assuming that SetFinalizer is not called again, the next time the garbage
collector sees that obj is unreachable, it will free obj.

两阶段清扫：
1. 发现finalizer，加入独立协程，解除关联finalizer
2. GC再次运行时 ，真正释放对象

*/

type Foo struct {
	a int
}

var commentFinalizer *bool

func init() {
	commentFinalizer = flag.Bool("c", false, "commented finalizer")
	flag.Parse()
}

func main() {
	debug.SetGCPercent(-1)

	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)

	fmt.Printf("Allocation: %f Mb, Number of allocation: %d\n", float32(ms.HeapAlloc)/float32(1024*1204), ms.HeapObjects)

	for i := 0; i < 1000000; i++ {
		f := NewFoo(i)
		_ = fmt.Sprintf("%d", f.a)
	}

	runtime.ReadMemStats(&ms)
	fmt.Printf("Allocation: %f Mb, Number of allocation: %d\n", float32(ms.HeapAlloc)/float32(1024*1204), ms.HeapObjects)

	runtime.GC()
	time.Sleep(time.Second)

	runtime.ReadMemStats(&ms)
	fmt.Printf("Allocation: %f Mb, Number of allocation: %d\n", float32(ms.HeapAlloc)/float32(1024*1204), ms.HeapObjects)

	runtime.GC()
	time.Sleep(time.Second)

	runtime.ReadMemStats(&ms)
	fmt.Printf("Allocation: %f Mb, Number of allocation: %d\n", float32(ms.HeapAlloc)/float32(1024*1204), ms.HeapObjects)
}

//go:noinline
func NewFoo(i int) *Foo {
	f := &Foo{a: rand.Intn(50)}

	if !*commentFinalizer {
		runtime.SetFinalizer(f, func(f *Foo) {
			_ = fmt.Sprintf("foo " + strconv.Itoa(i) + " has been garbage collected")
		})
	}

	return f
}

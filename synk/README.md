# Go 标准库源码学习（一）详解短小精悍的 Once

读了文章[Go 标准库源码学习（一）详解短小精悍的 Once](https://mp.weixin.qq.com/s/Lsm-BMdKCKNQjRndNCLwLw)，写了一个不用atomic的对比版本。


```go
package synk

import (
	"sync"
)

// Once2 is an object that will perform exactly one action.
type Once2 struct {
	done uint32
	m    sync.Mutex
}

func (o *Once2) Do(f func()) {
	if o.done == 0 {
		o.doSlow(f)
	}
}

func (o *Once2) doSlow(f func()) {
	o.m.Lock()
	defer o.m.Unlock()
	if o.done == 0 {
		defer func() {
			o.done = 1
		}()
		f()
	}
}
```

把测试用例也拷贝过来，跑测试用例，通过：

```bash
$ go test -run .
PASS
ok  	github.com/bingoohuang/golang-trial/synk	0.005s
```

跑性能测试

```bash
$ go test -bench=.
goos: darwin
goarch: amd64
pkg: github.com/bingoohuang/golang-trial/synk
BenchmarkOnce2-12    	1000000000	         0.165 ns/op
BenchmarkOnce-12     	1000000000	         0.168 ns/op
PASS
ok  	github.com/bingoohuang/golang-trial/synk	0.379s
```

竟然不用atomic的性能更高，哈哈哈。

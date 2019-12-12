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

竟然不用atomic的性能更高，哈哈哈，但是可能存在问题，就是done=1，对于其它goroutine不可见，导致每次都进入到互斥锁。

详细解释[Go语言并发编程03 - 并发的内存模型](https://chai2010.cn/post/2018/go-concurrency-03/)


不同Goroutine之间: 不满足顺序一致性!
如果我们将初始化msg和done的代码放到另一个Goroutine中，情况就完成不一样了！下面的并发代码将是错误的：

```go
var msg string
var done bool = false
func main() {
    go func() {
        msg = "hello, world"
        done = true
    }()
    for {
        if done {
            println(msg); break
        }
        println("retry...")
    }
}
```

运行时，大概有几种错误类型：一是main函数无法看到被修改后的done，因此main的for循环无法正常结束；二是main函数虽然看到了done被修改为true，但是msg依然没有初始化，这将导致错误的输出。

出现上述错误的原因是因为，Go语言的内存模型明确说明不同Goroutine之间不满足顺序一致性！同时编译器为了优化代码，进行初始化的Goroutine可能调整msg和done的执行顺序。main函数并不能从done状态的变化推导msg的初始化状态。

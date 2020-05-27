package main

import (
	"fmt"
	"reflect"
	"runtime"
	"time"
)

// GO语言的修饰器编程 https://coolshell.cn/articles/17929.html
// https://gist.github.com/saelo/4190b75724adc06b1c5a
func getFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func Decorate(impl interface{}) interface{} {
	fn := reflect.ValueOf(impl)

	inner := func(in []reflect.Value) []reflect.Value {
		f := reflect.ValueOf(impl)
		fmt.Println("before")

		defer func(t time.Time) {
			fmt.Printf("--- Time Elapsed (%s): %v ---\n",
				getFunctionName(impl), time.Since(t))
		}(time.Now())
		ret := f.Call(in)
		fmt.Println("after")
		return ret
	}

	v := reflect.MakeFunc(fn.Type(), inner)

	return v.Interface()
}

func add(a, b int) int {
	return a + b
}

func main() {
	var add2 = Decorate(add).(func(a, b int) int)
	fmt.Println(add2(1, 2))
}

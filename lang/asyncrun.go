package main

import (
	"bytes"
	"runtime/debug"
	"fmt"
	"sync"
)

// https://medium.com/compass-true-north/concurrent-programming-in-go-de33441ace1c

type GenericFunction func() error

func RunAsyncAllowErrors(functions ...GenericFunction) []error {
	var wg sync.WaitGroup
	wg.Add(len(functions))
	defer wg.Wait()

	errors := make([]error, len(functions))
	f := func(j int) {
		defer wg.Done()
		// Clause for handling panic errors
		defer func() {
			if r := recover(); r != nil {
				// Skip 4 stack frames:
				// 1) debug.Stack()
				// 2) formatStack()
				// 3) this anonymous func
				// 4) runtime/panic
				errors[j] = fmt.Errorf( "panic in async function: %v\n%s", r, formatStack(4))
			}
		}()
		errors[j] = functions[j]()
	}

	for j := range functions {
		go f(j)
	}

	return errors
}

// Return formatted stack trace, skipping "skip" leading stack frames
func formatStack(skip int) string {
	lines := bytes.Split(bytes.TrimSpace(debug.Stack()), []byte("\n"))
	formatted := bytes.Join(lines[1+2*skip:], []byte("\n"))
	return string(formatted)
}


func makeRequest(url string) GenericFunction {
	return func() error {
		panic(url)
		return fmt.Errorf("%s", url)
	}
}

func  main() {
	funcs := []GenericFunction{
		makeRequest("http://abc.com"),
		makeRequest("http://def.com"),
		makeRequest("http://ghi.com"),
	}
	for _, err := range RunAsyncAllowErrors(funcs...) {
		fmt.Println(err)
	}
}
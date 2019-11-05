package main

import (
	"fmt"
	"plugin"
)

// Parser use to parse things.
type Parser interface {
	Parse([]byte) (meta map[string]string, data map[string]float64, err error)
}

func pa() {
	plug, err := plugin.Open("./plugins/car.so")
	if err != nil {
		panic(err)
	}
	car, err := plug.Lookup("Car")
	if err != nil {
		panic(err)
	}
	p, ok := car.(Parser)
	if ok {
		meta, data, err := p.Parse([]byte("a"))
		if err != nil {
			panic(err)
		}
		fmt.Printf("meta: %v, data: %v \n", meta, data)
	}
}
func pb() {
	plug, err := plugin.Open("./plugins/phone.so")
	if err != nil {
		panic(err)
	}
	phone, err := plug.Lookup("Phone")
	if err != nil {
		panic(err)
	}
	p, ok := phone.(Parser)
	if ok {
		meta, data, _ := p.Parse([]byte("a"))
		fmt.Printf("meta: %v, data: %v \n", meta, data)
	}
}
func main() {
	pa()
	pb()
}

package main

import (
	"fmt"
	"testing"
)

func TestSlice(t *testing.T) {
	slice0()
	fmt.Println()

	slice1()
	fmt.Println()

	slice2()
}

func slice0() {
	a := []int{1}
	a = append(a, 2)
	a = append(a, 3)

	b := append(a, 4)
	c := append(a, 5)

	fmt.Printf("a:%v\n", a)
	fmt.Printf("b:%v\n", b)
	fmt.Printf("c:%v\n", c)
}

func slice1() {
	a := []int{1}
	fmt.Printf("1a:%+v, cap:%d, pointer:%p\n", a, cap(a), &a)
	a = append(a, 2)
	fmt.Printf("2a:%+v, cap:%d, pointer:%p\n", a, cap(a), &a)
	a = append(a, 3)
	fmt.Printf("3a:%+v, cap:%d, pointer:%p\n", a, cap(a), &a)

	b := append(a, 4)
	fmt.Printf("4b:%+v, cap:%d, pointer:%p\n", b, cap(b), &b)
	c := append(a, 5)
	fmt.Printf("5c:%+v, cap:%d, pointer:%p\n", c, cap(c), &c)

	fmt.Printf("4b:%+v, cap:%d, pointer:%p\n", b, cap(b), &b)
}

func slice2() {
	a := []int{1}
	fmt.Printf("1a:%+v, cap:%d, pointer:%p\n", a, cap(a), &a)
	a = append(a, 2)
	fmt.Printf("2a:%+v, cap:%d, pointer:%p\n", a, cap(a), &a)
	a = append(a, 3)
	fmt.Printf("3a:%+v, cap:%d, pointer:%p\n", a, cap(a), &a)

	b := append(a, 4)
	fmt.Printf("4b:%+v, cap:%d, pointer:%p\n", b, cap(b), &b)
	c := append(a, 5, 6)
	fmt.Printf("5c:%+v, cap:%d, pointer:%p\n", c, cap(c), &c)

	fmt.Printf("4b:%+v, cap:%d, pointer:%p\n", b, cap(b), &b)
}

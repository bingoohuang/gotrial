package main

import "fmt"

type Book struct {
	pages int
}

func (b Book) Pages() int {
	return b.pages
}
func (b *Book) SetPages(pages int) {
	b.pages = pages
}
func main() {
	var book Book
	// 调用这两个隐式声明的函数。
	(*Book).SetPages(&book, 123)
	fmt.Println(Book.Pages(book))     // 123
	fmt.Println((*Book).Pages(&book)) // 123

	s := "hello"
	sayHello([]byte(s))
	fmt.Print(s)

}

func sayHello(s []byte) {
	s[0] = 'B'
}

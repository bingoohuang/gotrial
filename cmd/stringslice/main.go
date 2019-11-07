package main

import "fmt"

// 【Go】高效截取字符串的一些思考
// https://blog.thinkeridea.com/201910/go/efficient_string_truncation.html?hmsr=toutiao.io&utm_medium=toutiao.io&utm_source=toutiao.io
func main() {
	s := "abcdef"
	fmt.Println(s[1:4]) // bcd

	s = "Go 语言"
	fmt.Println(s[1:4]) // 乱码 o �

	rs := []rune(s)
	fmt.Println(string(rs[1:4]))

}

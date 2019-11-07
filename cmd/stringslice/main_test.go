package main

import (
	"testing"
	"unicode/utf8"

	"github.com/thinkeridea/go-extend/exunicode/exutf8"
)

var benchmarkSubString = "Go语言是Google开发的一种静态强类型、编译型、并发型，并具有垃圾回收功能的编程语言。为了方便搜索和识别，有时会将其称为Golang。"
var benchmarkSubStringLength = 20

func SubStrRunes(s string, length int) string {
	if utf8.RuneCountInString(s) > length {
		rs := []rune(s)
		return string(rs[:length])
	}

	return s
}

func BenchmarkSubStrRunes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SubStrRunes(benchmarkSubString, benchmarkSubStringLength)
	}
}

/*
$ go test -bench=.
goos: darwin
goarch: amd64
pkg: github.com/bingoohuang/golang-trial/cmd/stringslice
BenchmarkSubStrRunes-12    	  865286	      1223 ns/op
PASS
ok  	github.com/bingoohuang/golang-trial/cmd/stringslice	1.078s

对 69 个的字符串截取前 20 个字符需要大概 1.3 微秒，这极大的超出了我的心里预期，我发现因为类型转换带来了内存分配，
这产生了一个新的字符串，并且类型转换需要大量的计算。
*/

/*
救命稻草 - utf8.DecodeRuneInString
我想改善类型转换带来的额外运算和内存分配，我仔细的梳理了一遍 strings 包，发现并没有相关的工具，这时我想到了 utf8 包，
它提供了多字节计算相关的工具，实话说我对它并不熟悉，或者说没有主动（直接）使用过它，我查看了它所有的文档发现
utf8.DecodeRuneInString 函数可以转换单个字符，并给出字符占用字节的数量，我尝试了如此下的实验：
*/
func SubStrDecodeRuneInString(s string, length int) string {
	var size, n int
	for i := 0; i < length && n < len(s); i++ {
		_, size = utf8.DecodeRuneInString(s[n:])
		n += size
	}

	return s[:n]
}

func BenchmarkSubStrDecodeRuneInString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SubStrDecodeRuneInString(benchmarkSubString, benchmarkSubStringLength)
	}
}

/*
$ go test -bench=.
goos: darwin
goarch: amd64
pkg: github.com/bingoohuang/golang-trial/cmd/stringslice
BenchmarkSubStrRunes-12                 	  810282	      1240 ns/op
BenchmarkSubStrDecodeRuneInString-12    	14441078	        82.4 ns/op
PASS
ok  	github.com/bingoohuang/golang-trial/cmd/stringslice	2.303s

较 []rune 类型转换效率提升了 13倍，消除了内存分配，它的确令人激动和兴奋，我迫不及待的回复了 “hollowaykeanho”
告诉他我发现了一个更好的方法，并提供了相关的性能测试。

我有些小激动，兴奋的浏览着论坛里各种有趣的问题，在查看一个问题的帮助时 (忘记是哪个问题了-_-||) ，我惊奇的发现了另一个思路。

*/

/*
良药不一定苦 - range 字符串迭代
许多人似乎遗忘了 range 是按字符迭代的，并非字节。使用 range 迭代字符串时返回字符起始索引和对应的字符，
我立刻尝试利用这个特性编写了如下用例
*/
func SubStrRange(s string, length int) string {
	var n, i int
	for i = range s {
		if n == length {
			break
		}

		n++
	}

	return s[:i]
}

func BenchmarkSubStrRange(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SubStrRange(benchmarkSubString, benchmarkSubStringLength)
	}
}

/*

$ go test -bench=.
goos: darwin
goarch: amd64
pkg: github.com/bingoohuang/golang-trial/cmd/stringslice
BenchmarkSubStrRunes-12                 	  842479	      1229 ns/op
BenchmarkSubStrDecodeRuneInString-12    	14430904	        82.6 ns/op
BenchmarkSubStrRange-12                 	11979513	        99.2 ns/op
PASS
ok  	github.com/bingoohuang/golang-trial/cmd/stringslice	3.627s

它仅仅提升了13%(??)，但它足够的简单和易于理解，这似乎就是我苦苦寻找的那味良药。

如果你以为这就结束了，不、这对我来只是探索的开始。

终极时刻 - 自己造轮子
喝了 range 那碗甜的腻人的良药，我似乎冷静下来了，我需要造一个轮子，它需要更易用，更高效。

于是乎我仔细观察了两个优化方案，它们似乎都是为了查找截取指定长度字符的索引位置，如果我可以提供一个这样的方法，是否就可以提供用户一个简单的截取实现 s[:strIndex(20)] ，这个想法萌芽之后我就无法再度摆脱，我苦苦思索两天来如何来提供易于使用的接口。

之后我创造了 exutf8.RuneIndexInString 和 exutf8.RuneIndex 方法，分别用来计算字符串和字节切片中指定字符数量结束的索引位置。

我用 exutf8.RuneIndexInString 实现了一个字符串截取测试：
*/

func SubStrRuneSubString(s string, length int) string {
	return exutf8.RuneSubString(s, 0, length)
}

func BenchmarkSubStrRuneSubString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SubStrRuneSubString(benchmarkSubString, benchmarkSubStringLength)
	}
}

/*
$ go test -bench=.
goos: darwin
goarch: amd64
pkg: github.com/bingoohuang/golang-trial/cmd/stringslice
BenchmarkSubStrRunes-12                 	  819668	      1233 ns/op
BenchmarkSubStrDecodeRuneInString-12    	14313582	        83.8 ns/op
BenchmarkSubStrRange-12                 	11969416	       102 ns/op
BenchmarkSubStrRuneSubString-12         	16747983	        71.0 ns/op
PASS
ok  	github.com/bingoohuang/golang-trial/cmd/stringslice	4.901s

虽然相较 exutf8.RuneIndexInString 有所下降，但它提供了易于交互和使用的接口，我认为这应该是最实用的方案，
如果你追求极致仍然可以使用 exutf8.RuneIndexInString，它依然是最快的方案。
*/

package main

import (
	"flag"
	"fmt"
	"github.com/bingoohuang/golang-trial/randimg"
	"github.com/bingoohuang/gou/rand"
)

var (
	width  int
	height int
	picfmt string
	fixMib int64
	many   int
)

func init() {
	flag.IntVar(&width, "w", 640, "picture width")
	flag.IntVar(&height, "h", 320, "picture height")
	flag.Int64Var(&fixMib, "m", 0, "fixed size(MiB)")
	flag.IntVar(&many, "i", 1, "how many pictures to create")
	flag.StringVar(&picfmt, "f", "png", "picture format(png/jpg)")

	flag.Parse()
}

func main() {
	s := rand.Int()

	for i := 0; i < many; i++ {
		randText := fmt.Sprintf("%d", s+i)
		fileName := fmt.Sprintf("%d.%s", s, picfmt)
		randimg.GenerateRandomImageFile(width, height, randText, fileName, fixMib<<20)

		fmt.Println(fileName, "with width", width, "height", height, "randText", randText, "generated!")
	}
}

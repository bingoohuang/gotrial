package main

import (
	"flag"
	"fmt"

	humanize "github.com/dustin/go-humanize"

	"github.com/bingoohuang/golang-trial/randimg"
	"github.com/bingoohuang/gou/ran"
)

var (
	width     int
	height    int
	picfmt    string
	fixedSize uint64
	many      int
)

func init() {
	flag.IntVar(&width, "w", 640, "picture width")
	flag.IntVar(&height, "h", 320, "picture height")
	fixedSizeStr := flag.String("s", "", "fixed size(eg. 44kB, 17MB)")
	flag.IntVar(&many, "i", 1, "how many pictures to create")
	flag.StringVar(&picfmt, "f", "png", "picture format(png/jpg)")

	flag.Parse()

	if *fixedSizeStr == "" {
		fixedSize = 10 << 20 // 10MiB
	} else {
		var err error
		fixedSize, err = humanize.ParseBytes(*fixedSizeStr)
		if err != nil {
			panic("illegal fixed size " + err.Error())
		}
	}
}

func main() {
	s := ran.Int()

	for i := 0; i < many; i++ {
		randText := fmt.Sprintf("%d", s+i)
		fileName := fmt.Sprintf("%d.%s", s, picfmt)
		randimg.GenerateRandomImageFile(width, height, randText, fileName, int64(fixedSize))

		fmt.Println(fileName, "with width", width, "height", height, "randText", randText,
			"fixedSize", humanize.IBytes(fixedSize), "generated!")
	}
}

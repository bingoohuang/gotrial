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
	fastMode  bool
)

func init() {
	flag.IntVar(&width, "w", 640, "picture width")
	flag.IntVar(&height, "h", 320, "picture height")
	fixedSizeStr := flag.String("s", "", "fixed size(eg. 44kB, 17MB)")
	flag.IntVar(&many, "i", 1, "how many pictures to create")
	flag.StringVar(&picfmt, "f", "png", "picture format(png/jpg)")
	flag.BoolVar(&fastMode, "fast", false, "fast mode")

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
		fileName := fmt.Sprintf("%d.%s", s+i, picfmt)
		c := randimg.RandImageConfig{
			Width:      width,
			Height:     height,
			RandomText: randText,
			FileName:   fileName,
			FixedSize:  int64(fixedSize),
			FastMode:   fastMode,
		}
		c.GenerateFile()

		fmt.Println(fileName, "with ", width, "x", height, "randText", randText,
			"size", humanize.IBytes(fixedSize), "generated!")
	}
}

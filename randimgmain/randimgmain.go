package main

import (
	"flag"
	"math/rand"
	"strconv"
	"time"

	"github.com/bingoohuang/golang-trial/randimg"
)

var (
	width    int
	height   int
	filename string
	seq      int
	picfmt   string
	fixMib   int64
	many     int
)

func init() {
	flag.IntVar(&width, "w", 640, "picture width")
	flag.IntVar(&height, "h", 320, "picture height")
	flag.IntVar(&seq, "s", 0, "picture sequence number")
	flag.Int64Var(&fixMib, "m", 0, "fixed size(MiB)")
	flag.IntVar(&many, "i", 1, "how many pictures to create")
	flag.StringVar(&picfmt, "f", "png", "picture format(png/jpg)")

	flag.Parse()
}

func main() {
	s := uint64(seq)
	if s <= 0 {
		rand.Seed(time.Now().UnixNano())
		s = randimg.RandUint64()
	}

	for i := 0; i < many; i++ {
		randText := strconv.FormatUint(s, 10)
		s++
		randimg.GenerateRandomImageFile(width, height, randText, randText+"."+picfmt, fixMib<<20)
	}
}
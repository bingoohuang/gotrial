package main

import (
	"flag"
	"os"

	"github.com/averagesecurityguy/random"
	humanize "github.com/dustin/go-humanize"
	"github.com/sirupsen/logrus"
)

func main() {
	fileSizeStr := ""
	outFileName := ""

	flag.StringVar(&fileSizeStr, "s", "1MiB", "file size(eg. 1MiB, 1G)")
	flag.StringVar(&outFileName, "o", "out.txt", "out file name")

	flag.Parse()

	fileSize, err := humanize.ParseBytes(fileSizeStr)
	if err != nil {
		logrus.Panicf("unable to parse file size %s error %v", fileSizeStr, err)
	}

	outFile, err := os.OpenFile(outFileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		logrus.Panicf("os.OpenFile %s error %v", outFileName, err)
	}

	spaces := "\n \t\r\v\f"
	l := 0
	fileSizeInt := int(fileSize)

	for l < fileSizeInt {
		randLen, err := random.Uint64Range(3, 20)
		if err != nil {
			logrus.Panicf("random.Uint64Range error %v", err)
		}

		alphaNum, err := random.AlphaNum(randLen)
		if err != nil {
			logrus.Panicf("random.AlphaNum %d error %v", randLen, err)
		}

		_, _ = outFile.WriteString(alphaNum)

		randSpace, err := random.Uint64Range(0, uint64(len(spaces)))
		if err != nil {
			logrus.Panicf("random.Uint64Range error %v", err)
		}

		_, _ = outFile.WriteString(string(spaces[randSpace]))

		l += len(alphaNum) + 1
	}

	_ = outFile.Close()
}

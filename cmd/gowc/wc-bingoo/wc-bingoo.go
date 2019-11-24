package main

import (
	"io"
	"log"
	"os"
	"runtime"
)

func main() {
	if len(os.Args) < 2 {
		panic("no file path specified")
	}

	filePath := os.Args[1]
	fileStat, err := os.Stat(filePath)
	if err != nil {
		panic(err)
	}
	fileSize := int(fileStat.Size())

	counts := make(chan Count)
	numWorkers := runtime.NumCPU()
	workerSize := fileSize / numWorkers

	for i := 0; i < numWorkers; i++ {
		startSize := i * workerSize
		endSize := startSize + workerSize
		if endSize > fileSize {
			endSize = fileSize
		}

		go FileReaderCounter(filePath, counts, startSize-1, endSize)
	}

	totalCount := Count{}
	for i := 0; i < numWorkers; i++ {
		totalCount.Add(<-counts)
	}
	close(counts)

	println(totalCount.LineCount, totalCount.WordCount, fileSize, fileStat.Name())
}

func FileReaderCounter(filePath string, counts chan Count, startSize, endSize int) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	if startSize > 0 {
		if _, err := file.Seek(int64(startSize), io.SeekStart); err != nil {
			log.Fatal(err)
		}
	}

	const bufferSize = 16 * 1024
	buffer := make([]byte, bufferSize)
	totalCount := Count{}
	lastCharIsSpace := false

	countBytes := endSize - startSize
	for readBytes := 0; readBytes < countBytes; {
		bytes, err := file.Read(buffer)
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}

		bufferStart := 0
		if readBytes == 0 {
			bufferStart = 1
			lastCharIsSpace = IsSpace(buffer[0])
		}

		readBytes += bytes
		if readBytes > countBytes {
			bytes -= readBytes - countBytes
		}

		totalCount.Add(GetCount(lastCharIsSpace, buffer[bufferStart:bytes]))
		lastCharIsSpace = IsSpace(buffer[bytes-1])
	}

	counts <- totalCount
}

type Count struct {
	LineCount int
	WordCount int
}

func (c *Count) Add(count Count) {
	c.LineCount += count.LineCount
	c.WordCount += count.WordCount
}

func GetCount(prevCharIsSpace bool, buffer []byte) (count Count) {
	for _, b := range buffer {
		if IsSpace(b) {
			if b == '\n' {
				count.LineCount++
			}

			prevCharIsSpace = true
		} else if prevCharIsSpace {
			prevCharIsSpace = false
			count.WordCount++
		}
	}

	return
}

func IsSpace(b byte) bool {
	switch b {
	case ' ', '\t', '\r', '\v', '\f', '\n':
		return true
	default:
		return false
	}
}

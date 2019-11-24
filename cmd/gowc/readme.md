# Beating C with 70 Lines of Go

[Beating C with 70 Lines of Go](https://ajeetdsouza.github.io/blog/posts/beating-c-with-70-lines-of-go/), [source code](https://github.com/ajeetdsouza/blog-wc-go)

## Install gtime on OSX with Homebrew: `brew install gnu-time`

The time command on OSX has less features than gnu-time gtime on Linux. [Time command on OSX, Linux](https://gist.github.com/gregelin/9529716)

## My trial environment

1. go version go1.13.4 darwin/amd64
1. Hardware Overview `system_profiler SPHardwareDataType Hardware`:

    * Model Name: MacBook Pro
    * Model Identifier: MacBookPro15,1
    * Processor Name: 6-Core Intel Core i7
    * Processor Speed: 2.6 GHz
    * Number of Processors: 1
    * Total Number of Cores: 6
    * L2 Cache (per Core): 256 KB
    * L3 Cache: 9 MB
    * Hyper-Threading Technology: Enabled
    * Memory: 16 GB
    * Boot ROM Version: 1037.40.124.0.0 (iBridge: 17.16.11081.0.0,0)

## Create a 100M and 1G ascii file for test

1. `go install cmd/gowc/randascii/randascii.go`
1. `randascii -s 100MiB -o 100m.txt `
1. `randascii -s 1GiB -o 1g.txt `

## gtime for system wc

```bash
$ gtime -f "%es %MKB" wc 100m.txt 
 1455778 8736963 104857603 100m.txt
0.28s 1728KB
$ gtime -f "%es %MKB" wc 1g.txt   
 14910007 89478242 1073741824 1g.txt
2.75s 1728KB
```

|x|100m|1g|
|:---:|:---:|:---:|
|0 wc|0.28s 1728KB|2.75s 1728KB|
|1 naive|0.51s 1816KB|5.22s 1816KB|
|2 chunks|0.22s 1780KB|2.19s 1788KB|
|3 channel|0.06s 6948KB|0.48s 8500KB|
|4 mutex|0.04s 2436KB|0.32s 2432KB|
|5 bingoo|0.03s 2648KB|0.24s 2668KB|


## codes lines statistics from [tokei](https://github.com/XAMPPRocky/tokei)

```bash
$ tokei -f cmd/gowc/
-------------------------------------------------------------------------------------------------------
 Language                                    Files        Lines         Code     Comments       Blanks
-------------------------------------------------------------------------------------------------------
 Go                                              7          590          482            0          108
-------------------------------------------------------------------------------------------------------
 cmd/gowc/wc-naive/wc-naive.go                               56           47            0            9
 cmd/gowc/sliceshare/sliceshare.go                           25           19            0            6
 cmd/gowc/wc-bingoo/wc-bingoo.go                            124          102            0           22
 cmd/gowc/wc-channel/wc-channel.go                          112           95            0           17
 cmd/gowc/wc-mutex/wc-mutex.go                              128          105            0           23
 cmd/gowc/wc-chunks/wc-chunks.go                             86           70            0           16
 cmd/gowc/randascii/randascii.go                             59           44            0           15
-------------------------------------------------------------------------------------------------------
 Markdown                                        1           76           76            0            0
-------------------------------------------------------------------------------------------------------
 cmd/gowc/readme.md                                          76           76            0            0
-------------------------------------------------------------------------------------------------------
 Total                                           8          666          558            0          108
-------------------------------------------------------------------------------------------------------
```


## Beating C with 70 Lines of Go

Chris Penner’s recent article, [Beating C with 80 Lines of Haskell]https://chrispenner.ca/posts/wc), has generated quite some controversy over the Internet, and it has since turned into a game of trying to take on the venerable wc with different languages:

* [Ada](http://verisimilitudes.net/2019-11-11)
* [C](https://github.com/expr-fi/fastlwc/)
* [Common Lisp](http://verisimilitudes.net/2019-11-12)
* [Dyalog APL](https://ummaycoc.github.io/wc.apl/)
* [Futhark](https://futhark-lang.org/blog/2019-10-25-beating-c-with-futhark-on-gpu.html)
* [Haskell](https://chrispenner.ca/posts/wc)
* [Rust](https://medium.com/@martinmroz/beating-c-with-120-lines-of-rust-wc-a0db679fe920)

Today, we will be pitting Go against wc. Being a compiled language with excellent concurrency primitives, it should be trivial to achieve comparable performance to C.

While wc is also designed to read from stdin, handle non-ASCII text encodings, and parse command line flags ([manpage](https://ss64.com/osx/wc.html)), we will not be doing that here. Instead, like the articles mentioned above, we will focus on keeping our implementation as simple as possible.

The source code for this post can be found [here](https://github.com/ajeetdsouza/blog-wc-go).

Benchmarking & comparison
We will use the GNU [time](https://www.gnu.org/software/time/) utility to compare elapsed time and maximum resident set size.

$ /usr/bin/time -f "%es %MKB" wc test.txt
We will use the [same version of wc](https://opensource.apple.com/source/text_cmds/text_cmds-68/wc/wc.c.auto.html) as the original article, compiled with gcc 9.2.1 and -O3. For our own implementation, we will use go 1.13.4 (I did try gccgo too, but the results were not very promising). We will run all benchmarks with the following setup:

* Intel Core i5-6200U @ 2.30 GHz (2 physical cores, 4 threads)
* 4+4 GB RAM @ 2133 MHz
* 240 GB M.2 SSD
*  Fedora 31

For a fair comparison, all implementations will use a 16 KB buffer for reading input. The input will be two us-ascii encoded text files of 100 MB and 1 GB.

A naïve approach
Parsing arguments is easy, since we only require the file path:

```go
if len(os.Args) < 2 {
    panic("no file path specified")
}
filePath := os.Args[1]

file, err := os.Open(filePath)
if err != nil {
    panic(err)
}
defer file.Close()
```

We’re going to iterate through the text bytewise, keeping track of state. Fortunately, in this case, we require only 2 states:

* The previous byte was whitespace
* The previous byte was not whitespace

When going from a whitespace character to a non-whitespace character, we increment the word counter. This approach allows us to read directly from a byte stream, keeping memory consumption low.

```go
const bufferSize = 16 * 1024
reader := bufio.NewReaderSize(file, bufferSize)

lineCount := 0
wordCount := 0
byteCount := 0

prevByteIsSpace := true
for {
    b, err := reader.ReadByte()
    if err != nil {
        if err == io.EOF {
            break
        } else {
            panic(err)
        }
    }

    byteCount++

    switch b {
    case '\n':
        lineCount++
        prevByteIsSpace = true
    case ' ', '\t', '\r', '\v', '\f':
        prevByteIsSpace = true
    default:
        if prevByteIsSpace {
            wordCount++
            prevByteIsSpace = false
        }
    }
}
```

To display the result, we will use the native println() function - in my tests, importing the fmt library caused a ~400 KB increase in executable size!

```go
println(lineCount, wordCount, byteCount, file.Name())
```

Let’s run this:

versions|input size |	elapsed time |	max memory
---|:---:|:---:|:---:
wc	| 100 MB	| 0.58 s	| 2052 KB
wc-naive| 	100 MB	| 0.77 s| 	1416 kB
wc| 	1 GB	| 5.56 s| 	2036 KB
wc-naive	| 1 GB	| 7.69 s	| 1416 KB

The good news is that our first attempt has already landed us pretty close to C in terms of performance. In fact, we’re actually doing better in terms of memory usage!

## Splitting the input

While buffering I/O reads is critical to improving performance, calling ReadByte() and checking for errors in a loop introduces a lot of unnecessary overhead. We can avoid this by manually buffering our read calls, rather than relying on bufio.Reader.

To do this, we will split our input into buffered chunks that can be processed individually. Fortunately, to process a chunk, the only thing we need to know about the previous chunk (as we saw earlier) is if its last character was whitespace.

Let’s write a few utility functions:

```go
type Chunk struct {
    PrevCharIsSpace bool
    Buffer          []byte
}

type Count struct {
    LineCount int
    WordCount int
}

func GetCount(chunk Chunk) Count {
    count := Count{}

    prevCharIsSpace := chunk.PrevCharIsSpace
    for _, b := range chunk.Buffer {
        switch b {
        case '\n':
            count.LineCount++
            prevCharIsSpace = true
        case ' ', '\t', '\r', '\v', '\f':
            prevCharIsSpace = true
        default:
            if prevCharIsSpace {
                prevCharIsSpace = false
                count.WordCount++
            }
        }
    }

    return count
}

func IsSpace(b byte) bool {
    return b == ' ' || b == '\t' || b == '\n' || b == '\r' || b == '\v' || b == '\f'
}
```

Now, we can split the input into Chunks and feed them to the GetCount function.

```go
totalCount := Count{}
lastCharIsSpace := true

const bufferSize = 16 * 1024
buffer := make([]byte, bufferSize)

for {
    bytes, err := file.Read(buffer)
    if err != nil {
        if err == io.EOF {
            break
        } else {
            panic(err)
        }
    }

    count := GetCount(Chunk{lastCharIsSpace, buffer[:bytes]})
    lastCharIsSpace = IsSpace(buffer[bytes-1])

    totalCount.LineCount += count.LineCount
    totalCount.WordCount += count.WordCount
}
```

To obtain the byte count, we can make one system call to query the file size:

```go
fileStat, err := file.Stat()
if err != nil {
    panic(err)
}
byteCount := fileStat.Size()
```

Now that we’re done, let’s see how this performs:

versions| input size| 	elapsed time| 	max memory
---|:---:|:---:|:---:
wc| 	100 MB| 	0.58 s	| 2052 KB
wc-chunks	| 100 MB	| 0.34 s	| 1404 KB
wc| 	1 GB| 	5.56 s	| 2036 KB
wc-chunks	| 1 GB| 	3.31 s	| 1416 KB

Looks like we’ve blown past wc on both counts, and we haven’t even begun to parallelize our program yet. [tokei](https://github.com/XAMPPRocky/tokei) reports that this program is just 70 lines of code!

## Parallelization

Admittedly, a parallel wc is overkill, but let’s see how far we can go. The original article reads from the input file in parallel, and while it improved runtime, the author does admit that performance gains due to parallel reads might be limited to only certain kinds of storage, and would be detrimental elsewhere.

For our implementation, we want our code to be performant on all devices, so we will not be doing this. We will set up 2 channels, chunks and counts. Each worker will read and process data from chunks until the channel is closed, and then write the result to counts.

```go
func ChunkCounter(chunks <-chan Chunk, counts chan<- Count) {
    totalCount := Count{}
    for {
        chunk, ok := <-chunks
        if !ok {
            break
        }
        count := GetCount(chunk)
        totalCount.LineCount += count.LineCount
        totalCount.WordCount += count.WordCount
    }
    counts <- totalCount
}
```

We will spawn one worker per logical CPU core:

```go
numWorkers := runtime.NumCPU()

chunks := make(chan Chunk)
counts := make(chan Count)

for i := 0; i < numWorkers; i++ {
    go ChunkCounter(chunks, counts)
}
```

Now, we run in a loop, reading from the disk and assigning jobs to each worker:

```go
const bufferSize = 16 * 1024
lastCharIsSpace := true

for {
    buffer := make([]byte, bufferSize)
    bytes, err := file.Read(buffer)
    if err != nil {
        if err == io.EOF {
            break
        } else {
            panic(err)
        }
    }
    chunks <- Chunk{lastCharIsSpace, buffer[:bytes]}
    lastCharIsSpace = IsSpace(buffer[bytes-1])
}
close(chunks)
Once this is done, we can simply sum up the counts from each worker:

totalCount := Count{}
for i := 0; i < numWorkers; i++ {
    count := <-counts
    totalCount.LineCount += count.LineCount
    totalCount.WordCount += count.WordCount
}
close(counts)
```

Let’s run this and see how it compares to the previous results:

versions| input size| 	elapsed time| 	max memory
---|:---:|:---:|:---:
wc|	100 MB	|0.58 s|	2052 KB
wc-channel|	100 MB|	0.27 s	|6644 KB
wc|	1 GB	|5.56 s	|2036 KB
wc-channel	|1 GB	|2.22 s|	6752 KB

Our wc is now a lot faster, but there has been quite a regression in memory usage. In particular, notice how our input loop allocates memory at every iteration! Channels are a great abstraction over sharing memory, but for some use cases, simply not using channels can improve performance tremendously.

## Better parallelization

In this section, we will allow every worker to read from the file, and use sync.Mutex to ensure that reads don’t happen simultaneously. We can create a new struct to handle this for us:

```go
type FileReader struct {
    File            *os.File
    LastCharIsSpace bool
    mutex           sync.Mutex
}

func (fileReader *FileReader) ReadChunk(buffer []byte) (Chunk, error) {
    fileReader.mutex.Lock()
    defer fileReader.mutex.Unlock()

    bytes, err := fileReader.File.Read(buffer)
    if err != nil {
        return Chunk{}, err
    }

    chunk := Chunk{fileReader.LastCharIsSpace, buffer[:bytes]}
    fileReader.LastCharIsSpace = IsSpace(buffer[bytes-1])

    return chunk, nil
}
```

We then rewrite our worker function to read directly from the file:

```go
func FileReaderCounter(fileReader *FileReader, counts chan Count) {
    const bufferSize = 16 * 1024
    buffer := make([]byte, bufferSize)

    totalCount := Count{}

    for {
        chunk, err := fileReader.ReadChunk(buffer)
        if err != nil {
            if err == io.EOF {
                break
            } else {
                panic(err)
            }
        }
        count := GetCount(chunk)
        totalCount.LineCount += count.LineCount
        totalCount.WordCount += count.WordCount
    }

    counts <- totalCount
}
```

Like earlier, we can now spawn these workers, one per CPU core:

```go
fileReader := &FileReader{
    File:            file,
    LastCharIsSpace: true,
}
counts := make(chan Count)

for i := 0; i < numWorkers; i++ {
    go FileReaderCounter(fileReader, counts)
}

totalCount := Count{}
for i := 0; i < numWorkers; i++ {
    count := <-counts
    totalCount.LineCount += count.LineCount
    totalCount.WordCount += count.WordCount
}
close(counts)
```

Let’s see how this performs:

versions| input size| 	elapsed time| 	max memory
---|:---:|:---:|:---:
wc	|100 MB|	0.58 s	|2052 KB
wc-mutex	|100 MB	|0.12 s|	1580 KB
wc|	1 GB	|5.56 s	|2036 KB
wc-mutex	|1 GB	|1.21 s|	1576 KB

Our parallelled implementation runs at more than 4.5x the speed of wc, with lower memory consumption! This is pretty significant, especially if you consider that Go is a garbage collected language.

## Bingoo's version

Since we just read file, we can remove the mutex from the goroutines. Each goroutine starts to read from a separate position and deal a individual part statistics.

versions| input size| 	elapsed time| 	max memory
---|:---:|:---:|:---:
wc|100M|0.28s | 1728KB
wc-bingoo|100M|0.03s |2648KB
wc|1G|2.75s |1728KB
wc-bingoo|1G|0.24s |2668KB

we got 11x faster than the wc to deal 1G ascii file.

## Conclusion

While in no way does this article imply that Go > C, I hope it demonstrates that Go can be a viable alternative to C as a systems programming language.
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


## codes lines statistics from [tokei](https://github.com/XAMPPRocky/tokei)

```bash
$ tokei -f cmd/gowc/
---------------------------------------------------------------------------------------
 Language                    Files        Lines         Code     Comments       Blanks
---------------------------------------------------------------------------------------
 Go                              5          441          361            0           80
---------------------------------------------------------------------------------------
 cmd/gowc/wc-naive/wc-naive.go               56           47            0            9
 cmd/gowc/wc-channel/wc-channel.go          112           95            0           17
 cmd/gowc/wc-mutex/wc-mutex.go              128          105            0           23
 cmd/gowc/wc-chunks/wc-chunks.go             86           70            0           16
 cmd/gowc/randascii/randascii.go             59           44            0           15
---------------------------------------------------------------------------------------
```
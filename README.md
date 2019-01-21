# golang-trial
golang 试验田

## gossh 
test ssh in golang

> env GOOS=linux GOARCH=amd64 go build -o gossh.linux.bin

```bash
bingoobjca@bogon ~/g/s/g/b/g/gosshmain> ./gosshmain
Usage of ./gosshmain:
  -P string
    	password
  -f string
    	scripts file
  -h string
    	host
  -p int
    	port (default 22)
  -s string
    	scripts string
  -t string
    	prompt tip (default ">")
  -u string
    	user, default to current user
```

## randimg

create randomized image in jpg/png format.

```bash
bingoobjca@bogon ~/g/s/g/b/g/randimgmain> ./randimgmain -h
flag needs an argument: -h
Usage of ./randimgmain:
  -f string
    	picture format(png/jpg) (default "png")
  -h int
    	picture height (default 320)
  -i int
    	how many pictures to create (default 1)
  -m int
    	fixed size(MiB)
  -s int
    	picture sequence number
  -w int
    	picture width (default 640)
```

demo image:
![image](https://user-images.githubusercontent.com/1940588/51447936-afcb2180-1d5d-11e9-923d-9aef890bf51f.png)


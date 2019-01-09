# golang-trial
golang 试验田

## gossh 
test ssh in golang

> env GOOS=linux GOARCH=amd64 go build -o gossh.linux.bin

```bash
[~]$ ./gossh.linux.bin -help
Usage of ./gossh.linux.bin:
  -h string
        host
  -p int
        port (default 22)
  -s string
        scripts (default "uname -a")
  -u string
        user, default to current user
```

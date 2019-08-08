package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/sirupsen/logrus"
)

// http://www.songjiayang.com/posts/shi-yong-cgroup-shi-xian-nei-cun-kong-zhi
// https://wudaijun.com/2018/10/linux-cgroup/

const (
	cgroupRoot      = "/sys/fs/cgroup/memory"
	procsFile       = "cgroup.procs"
	memoryLimitFile = "memory.limit_in_bytes"
	swapLimitFile   = "memory.swappiness"
	MB              = 1024 * 1024
)

func main() {
	var rssLimit int
	var memoryPath string
	// 我们可以使用 mkdir /sys/fs/cgroup/memory/climits 来创建属于自己的内存组 climits
	flag.StringVar(&memoryPath, "path", "/climits", "a static path to use for memory cgroups.")
	flag.IntVar(&rssLimit, "memory", 10, "memory limit with MB.")
	flag.Parse()

	// set memory limit
	writeFile(filepath.Join(cgroupRoot, memoryPath, memoryLimitFile), rssLimit*MB)
	// set swap memory limit to zero
	writeFile(filepath.Join(cgroupRoot, memoryPath, swapLimitFile), 0)

	go startCmd("./cgapp", memoryPath)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	s := <-c
	logrus.Infoln("Got signal:", s)
}

func check(err error) {
	if err != nil {
		logrus.Panic(err)
	}
}

func writeFile(path string, value int) {
	check(ioutil.WriteFile(path, []byte(fmt.Sprintf("%d", value)), 0755))
}

type ExitStatus struct {
	Signal os.Signal
	Code   int
}

func startCmd(command, memoryPath string) {
	restart := make(chan ExitStatus, 1)

	runner := func() {
		cmd := exec.Cmd{Path: command}
		cmd.Stdout = os.Stdout

		// start app
		check(cmd.Start())

		logrus.Infoln("add pid", cmd.Process.Pid, "to file cgroup.procs")

		// set cgroup procs id
		writeFile(filepath.Join(cgroupRoot, memoryPath, procsFile), cmd.Process.Pid)

		if err := cmd.Wait(); err != nil {
			logrus.Infoln("cmd return with error:", err)
		}

		status := cmd.ProcessState.Sys().(syscall.WaitStatus)

		options := ExitStatus{Code: status.ExitStatus()}
		if status.Signaled() {
			options.Signal = status.Signal()
		}

		if err := cmd.Process.Kill(); err != nil {
			logrus.Infoln("cmd.Process.Kill error:", err)
		}

		restart <- options
	}

	go runner()

	for {
		status := <-restart

		switch status.Signal {
		case os.Kill:
			logrus.Infoln("app is killed by system")
		default:
			logrus.Infoln("app exit with code:", status.Code)
			return
		}

		logrus.Infoln("restart app..")
		go runner()
	}
}

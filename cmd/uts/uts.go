package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/sirupsen/logrus"
)

// Go语言实现容器namespace和cgroups
// https://juejin.im/entry/59abdb83f265da249412463a

// Build Container With Namespace and Cgroups in Go:
// https://kasheemlew.github.io/2017/09/02/build-container-with-go/

/*
[root@BJCA-device typhon-server]# ls -l /proc/$$/ns
总用量 0
lrwxrwxrwx. 1 root root 0 8月   7 13:52 ipc -> ipc:[4026531839]
lrwxrwxrwx. 1 root root 0 8月   7 13:52 mnt -> mnt:[4026531840]
lrwxrwxrwx. 1 root root 0 8月   7 13:52 net -> net:[4026531962]
lrwxrwxrwx. 1 root root 0 8月   7 13:52 pid -> pid:[4026531836]
lrwxrwxrwx. 1 root root 0 8月   7 13:52 user -> user:[4026531837]
lrwxrwxrwx. 1 root root 0 8月   7 13:52 uts -> uts:[4026531838]
[root@BJCA-device typhon-server]# ./uts run sh
INFO[0000] Running [sh]
sh-4.2# ls -l /proc/$$/ns
总用量 0
lrwxrwxrwx. 1 root root 0 8月   7 13:52 ipc -> ipc:[4026531839]
lrwxrwxrwx. 1 root root 0 8月   7 13:52 mnt -> mnt:[4026531840]
lrwxrwxrwx. 1 root root 0 8月   7 13:52 net -> net:[4026531962]
lrwxrwxrwx. 1 root root 0 8月   7 13:52 pid -> pid:[4026531836]
lrwxrwxrwx. 1 root root 0 8月   7 13:52 user -> user:[4026531837]
lrwxrwxrwx. 1 root root 0 8月   7 13:52 uts -> uts:[4026532199]
sh-4.2# hostname newhost
sh-4.2# hostname
newhost
sh-4.2# exit
exit
[root@BJCA-device typhon-server]# hostname
BJCA-device
[root@BJCA-device typhon-server]#
*/
func main() {
	if len(os.Args) < 2 {
		logrus.Errorf("missing commands")
		return
	}
	switch os.Args[1] {
	case "run1":
		run1()
	case "run2":
		run2()
	case "child2":
		child2()
	case "run3":
		run3()
	case "child3":
		child3()
	default:
		logrus.Errorf("wrong command")
		return
	}
}

func cg() {
	cgPath := "/sys/fs/cgroup/"
	pidsPath := filepath.Join(cgPath, "pids")
	// 在/sys/fs/cgroup/pids下创建container目录
	os.Mkdir(filepath.Join(pidsPath, "container"), 0755)
	// 设置最大进程数目为20
	check(ioutil.WriteFile(filepath.Join(pidsPath, "container/pids.max"), []byte("20"), 0700))
	// 将notify_on_release值设为1，当cgroup不再包含任何任务的时候将执行release_agent的内容
	check(ioutil.WriteFile(filepath.Join(pidsPath, "container/notify_on_release"), []byte("1"), 0700))
	// 加入当前正在执行的进程
	check(ioutil.WriteFile(filepath.Join(pidsPath, "container/cgroup.procs"), []byte(strconv.Itoa(os.Getpid())), 0700))
}

func run3() {
	logrus.Info("Setting up...")
	// 这里的/proc/self/exe就是当前正在执行的命令
	cmd := exec.Command("/proc/self/exe", append([]string{"child3"}, os.Args[2:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = cloneNewUtsNewPidSysProcAttr
	check(cmd.Run())
}

func child3() {
	logrus.Infof("Running %v", os.Args[2:])
	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cg()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	check(Sethostname("newhost"))
	/*
		https://kasheemlew.github.io/2017/09/02/build-container-with-go/：

		下面获取一个unix文件系统，可以选择docker的busybox镜像，并将其导出。

		docker pull busybox
		docker run -d busybox top -b
		此时获得刚刚的容器的containerID，然后执行

		docekr export -o busybox.tar <刚才容器的ID>
		即可在当前目录下得到一个busybox的压缩包，用

		mkdir busybox
		tar -xf busybox.tar -C busybox/
		解压即可得到我们需要的文件系统

		查看一下busybox目录

		$ ls busybox
		bin  dev  etc  home  proc  root  sys  tmp  usr  var

	*/
	check(syscall.Chroot("busybox"))
	check(os.Chdir("/"))
	// func Mount(source string, target string, fstype string, flags uintptr, data string) (err error)
	// 前三个参数分别是文件系统的名字，挂载到的路径，文件系统的类型
	check(Mount("proc", "proc", "proc", 0, ""))
	check(Mount("tmpdir", "tmp", "tmpfs", 0, ""))
	check(cmd.Run())
	check(syscall.Unmount("proc", 0))
	check(syscall.Unmount("tmp", 0))
}

func run2() {
	logrus.Info("Setting up...")
	// 这里的/proc/self/exe就是当前正在执行的命令
	cmd := exec.Command("/proc/self/exe", append([]string{"child2"}, os.Args[2:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = cloneNewUtsSysProcAttr
	check(cmd.Run())
}

func child2() {
	logrus.Infof("Running %v", os.Args[2:])
	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	check(Sethostname("newhost"))
	check(cmd.Run())
}

func run1() {
	logrus.Infof("Running %v", os.Args[2:])
	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = cloneNewUtsSysProcAttr
	check(cmd.Run())
}

func check(err error) {
	if err != nil {
		logrus.Errorln(err)
	}
}

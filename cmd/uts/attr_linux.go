package main

import "syscall"

/*
Linux内核实现了namespace，进而实现了轻量级虚拟化服务，在同一个namespace下的进程可以感知彼此的变化，
但是不能看到其他的进程，从而达到了环境隔离的目的。namespace有6项隔离，分别是:

1. UTS(Unix Time-sharing System, 主机和域名),
2. IPC(InterProcess Comms, 信号量、消息队列和共享内存),
3. PID(Process IDs, 进程编号),
4. Network(网络设备，网络栈，端口等),
5. Mount(挂载点[文件系统]),
6. User(用户和用户组)。

C语言中可以通过clone()指定flags参数，在创建进程的同时创建namespace。
Linux内核版本3.8之后的用户可以通过ls -l /proc/$$/ns查看当前进程指向的namespace编号。($$表示当前运行的进程ID号)

*/

var cloneNewUtsSysProcAttr = &syscall.SysProcAttr{
	Cloneflags: syscall.CLONE_NEWUTS, // 创建一个UTS隔离的新进程
}

var cloneNewUtsNewPidSysProcAttr = &syscall.SysProcAttr{
	Cloneflags: syscall.CLONE_NEWUTS | // 创建一个UTS隔离的新进程
		syscall.CLONE_NEWPID | // 进行PID的隔离
		syscall.CLONE_NEWNS, // mount point的隔离
}

func Sethostname(hostname string) error {
	return syscall.Sethostname([]byte(hostname))
}

func Mount(source string, target string, fstype string, flags uintptr, data string) error {
	return syscall.Mount(source, target, fstype, flags, data)
}

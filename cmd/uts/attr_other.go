// +build !linux

package main

import "syscall"

var cloneNewUtsSysProcAttr = &syscall.SysProcAttr{}
var cloneNewUtsNewPidSysProcAttr = &syscall.SysProcAttr{}

func Sethostname(hostname string) error {
	return nil
}

func Mount(source string, target string, fstype string, flags uintptr, data string) error {
	return nil
}

//go:build linux

package main

import "syscall"

// EVIOCGRAB from linux/input.h (_IOW('E', 0x90, int)) -- ARM/Linux ABI.
const eviocGrab = 0x40044590

func evdevGrab(fd uintptr) error {
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, fd, uintptr(eviocGrab), 1)
	if errno != 0 {
		return errno
	}
	return nil
}

func evdevUngrab(fd uintptr) error {
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, fd, uintptr(eviocGrab), 0)
	if errno != 0 {
		return errno
	}
	return nil
}

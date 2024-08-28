//go:build !windows
// +build !windows

package helper

import "golang.org/x/sys/unix"

func FcntlSet(fd int, flag int) error {
	var old int
	var err error
	for {
		old, err = unix.FcntlInt(uintptr(fd), unix.F_SETFL, flag)
		if old == -1 && err == unix.EINTR {
			continue
		}
		if err != nil {
			return err
		}
		break
	}

	return err
}

func FcntlAdd(fd int, flag int) error {
	var old int
	var err error
	for {
		old, err = unix.FcntlInt(uintptr(fd), unix.F_GETFL, 0)
		if old == -1 && err == unix.EINTR {
			continue
		}
		if err != nil {
			return err
		}
		if old&flag == flag {
			return nil
		}
		break
	}

	old |= flag

	return FcntlSet(fd, old)
}

func FcntlCancel(fd int, flag int) error {
	var old int
	var err error
	for {
		old, err = unix.FcntlInt(uintptr(fd), unix.F_GETFL, 0)
		if old == -1 && err == unix.EINTR {
			continue
		}
		if err != nil {
			return err
		}
		if old&flag == 0 {
			return nil
		}
		break
	}
	old &= ^flag

	return FcntlSet(fd, old)
}

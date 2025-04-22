package namespace

import (
	"fmt"
	"syscall"
)

func mountProc() error {
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	if err := syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), ""); err != nil {
		return fmt.Errorf("Failed to mount /proc: %s", err.Error())
	} else {
		return nil
	}
}

func mount() error {
	return mountProc()
}

package namespace

import (
	"fmt"
	"os"
	"strconv"
	"syscall"

	// "github.com/opencontainers/runtime-spec/specs-go"
	"golang.org/x/sys/unix"
)

func mountProc() error {
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	if err := syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), ""); err != nil {
		return fmt.Errorf("Failed to mount /proc: %s", err.Error())
	} else {
		return nil
	}
}

func prepareRootfs(rootfsDir string) error {
	flag := syscall.MS_SLAVE | syscall.MS_REC

	if err := syscall.Mount("", "/", "", uintptr(flag), ""); err != nil {
		return err
	}

	if err := syscall.Mount(rootfsDir, rootfsDir, "bind", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		return err
	}

	if err := syscall.Mount("", rootfsDir, "", syscall.MS_PRIVATE|syscall.MS_REC, ""); err != nil {
		return err
	}

	return nil
}

func pivotRootfs(rootfsDir string) error {

	oldroot, err := syscall.Open("/", syscall.O_DIRECTORY|syscall.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer syscall.Close(oldroot)

	newroot, err := syscall.Open(rootfsDir, unix.O_DIRECTORY|unix.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer syscall.Close(newroot)

	// Change to the new root so that the pivot_root actually acts on it.
	if err := syscall.Fchdir(newroot); err != nil {
		return &os.PathError{Op: "fchdir", Path: "fd " + strconv.Itoa(newroot), Err: err}
	}

	if err := unix.PivotRoot(".", "."); err != nil {
		return &os.PathError{Op: "pivot_root", Path: ".", Err: err}
	}

	// Currently our "." is oldroot (according to the current kernel code).
	// However, purely for safety, we will fchdir(oldroot) since there isn't
	// really any guarantee from the kernel what /proc/self/cwd will be after a
	// pivot_root(2).

	if err := unix.Fchdir(oldroot); err != nil {
		return &os.PathError{Op: "fchdir", Path: "fd " + strconv.Itoa(oldroot), Err: err}
	}

	// Make oldroot rslave to make sure our unmounts don't propagate to the
	// host (and thus bork the machine). We don't use rprivate because this is
	// known to cause issues due to races where we still have a reference to a
	// mount while a process in the host namespace are trying to operate on
	// something they think has no mounts (devicemapper in particular).
	if err := syscall.Mount("", ".", "", unix.MS_SLAVE|unix.MS_REC, ""); err != nil {
		return err
	}
	// Perform the unmount. MNT_DETACH allows us to unmount /proc/self/cwd.
	if err := syscall.Unmount(".", unix.MNT_DETACH); err != nil {
		return err
	}

	// Switch back to our shiny new root.
	if err := unix.Chdir("/"); err != nil {
		return &os.PathError{Op: "chdir", Path: "/", Err: err}
	}
	return nil
}

func mountRootfs(rootfsDir string) error {
	if err := prepareRootfs(rootfsDir); err != nil {
		return err
	}
	return pivotRootfs(rootfsDir)
}

func mount() error {

	// if err := mountRootfs("/home/tuchunxu/WorkSpace/Project/miniDocker/fs"); err != nil {
	// 	return err
	// }

	return mountProc()
	// return nil
}

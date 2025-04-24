package namespace

import (
	"fmt"
	"syscall"
)

func mountProc(src, dst string, options int) error {
	if err := syscall.Mount(src, dst, "proc", uintptr(options), ""); err != nil {
		return fmt.Errorf("Failed to mount %s to %s as proc: %v", src, dst, err)
	} else {
		return nil
	}
}

func mountTmpfs(src, dst string, options int) error {
	if err := syscall.Mount(src, dst, "tmpfs", uintptr(options), "mode=755"); err != nil {
		return fmt.Errorf("Failed to mount %s to %s as tmpfs: %v", src, dst, err)
	} else {
		return nil
	}
}

func pivotRootfs(rootfs string) error {
	// 把roofs bind mount到自己，这样它就是一个挂载点了
	if err := syscall.Mount(rootfs, rootfs, "bind", syscall.MS_BIND|syscall.MS_REC|syscall.MS_PRIVATE, ""); err != nil {
		return fmt.Errorf("Failed to mount rootfs to itself: %v", err)
	}
	// 用文件描述符号打开old root和new root，因为一会pivot root后就没法用名字找到了
	oldRoot, err := syscall.Open("/", syscall.O_DIRECTORY|syscall.O_RDONLY, 0)
	if err != nil {
		return fmt.Errorf("Failed to open \"/\": %v", err)
	}
	defer syscall.Close(oldRoot) //开的别忘了关

	newRoot, err := syscall.Open(rootfs, syscall.O_DIRECTORY|syscall.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer syscall.Close(newRoot) //开的别忘了关

	// cd到新的root里
	if err := syscall.Fchdir(newRoot); err != nil {
		return fmt.Errorf("Failed to fchdir to rootfs: %v", err)
	}

	// 这个神奇的操作会把原来的old root挂载到/proc/self/cwd，应该
	if err := syscall.PivotRoot(".", "."); err != nil {
		return fmt.Errorf("Failed to pivot root: %v", err)
	}

	// 到old root里把它unmount
	if err := syscall.Fchdir(oldRoot); err != nil {
		return fmt.Errorf("Failed to fchdir to old root: %v", err)
	}
	// 安全地递归unmount掉old root里的挂载点，不然你会看到host的monunt全留着
	if err := syscall.Mount("", ".", "", syscall.MS_SLAVE|syscall.MS_REC, ""); err != nil {
		return err
	}
	if err := syscall.Unmount(".", syscall.MNT_DETACH); err != nil {
		return err
	}

	if err := syscall.Chdir("/"); err != nil {
		return fmt.Errorf("Failed to chdir to new root: %v", err)
	}
	return nil
}

func mount() error {
	rootfs, err := syscall.Getwd()
	if err != nil {
		return fmt.Errorf("Failed to get work space: %v", err)
	}

	if err := mountProc("proc", rootfs+"/proc", syscall.MS_NOEXEC|syscall.MS_NOSUID|syscall.MS_NODEV); err != nil {
		return err
	}

	if err := mountTmpfs("tmpfs", rootfs+"/dev", syscall.MS_NOSUID|syscall.MS_STRICTATIME); err != nil {
		return err
	}

	if err := pivotRootfs(rootfs); err != nil {
		return fmt.Errorf("Failed to pivot rootfs: %v", err)
	}

	return nil
}

package namespace

import (
	"os"
	"os/exec"
	"syscall"

	"github.com/sirupsen/logrus"
)

type InitConfig struct {
	Interactive bool
}

func YieldInitProcess(config *InitConfig) (*exec.Cmd, *os.File) {
	cmd := exec.Command("/proc/self/exe", "init")
	readInitPipe, writeInitPipe, err := os.Pipe()
	if err != nil {
		return nil, nil
	}
	cmd.ExtraFiles = []*os.File{readInitPipe}

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUSER | syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
		UidMappings: []syscall.SysProcIDMap{
			{ContainerID: 0, HostID: os.Getuid(), Size: 1}, // 容器内root映射到当前用户
		},
		GidMappings: []syscall.SysProcIDMap{
			{ContainerID: 0, HostID: os.Getgid(), Size: 1},
		},
	}

	if config.Interactive {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
	}
	cmd.Stderr = os.Stderr

	return cmd, writeInitPipe
}

func mount() error {
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	if err := syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), ""); err != nil {
		logrus.Error("failed to mount /proc: ", err.Error())
		return err
	} else {
		logrus.Info("mount /proc sucessfully")
		return nil
	}
}

// func uts() error {
// 	logrus.Info("debug ", cgroups.ContainerId)
// 	return syscall.Sethostname([]byte(cgroups.ContainerId))
// }

// 这里，init进程需要初始化namespace里的环境
func RunInitProcess(command string, argv []string) error {

	if err := mount(); err != nil {
		return err
	}

	// if err := uts(); err != nil {
	// 	logrus.Error("failed to set UTS name: ", err.Error())
	// }

	if err := syscall.Exec(command, argv, os.Environ()); err != nil {
		logrus.Error("failed to exec command: ", err.Error())
		return err
	}

	return nil
}

package namespace

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"syscall"

	"github.com/opencontainers/runtime-spec/specs-go"
	"github.com/sirupsen/logrus"
)

const (
	STDIO_FD_COUNT = 3
)

type InitConfig struct {
	ContainerId string
	Interactive bool
}

func YieldInitProcess(config *InitConfig, spec *specs.Spec) (*exec.Cmd, *os.File) {
	cmd := exec.Command("/proc/self/exe", "init")

	readInitPipe, writeInitPipe, err := os.Pipe()
	if err != nil {
		return nil, nil
	}
	cmd.ExtraFiles = append(cmd.ExtraFiles, readInitPipe)
	cmd.Env = append(cmd.Env, "_INIT_PIPE="+strconv.Itoa(STDIO_FD_COUNT+len(cmd.ExtraFiles)-1))

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

// 这里，init进程需要初始化namespace里的环境
func RunInitProcess(specPtr *specs.Spec) error {

	if err := mount(specPtr); err != nil {
		return err
	} else {
		logrus.Info("mount sucessfully")
	}

	if err := syscall.Exec(specPtr.Process.Args[0], specPtr.Process.Args, os.Environ()); err != nil {
		return fmt.Errorf("Failed to exec command: %s", err.Error())
	}

	return nil
}

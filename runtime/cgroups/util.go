package cgroups

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"os"
	"os/user"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

var CgroupsRoot string
var ContainerId string

func YieldContainerId() string {
	// 组合输入数据(时间戳 + 随机盐)
	data := fmt.Sprintf("%d-%d", time.Now().UnixNano(), rand.Intn(1000))
	// 计算 SHA256 哈希
	hash := sha256.Sum256([]byte(data))
	// 转换为十六进制并截取前 12 位
	return hex.EncodeToString(hash[:])[:12]
}

// 创建cgroup，如果成功的话会设置CgroupsRoot
func CreateCgroupsRoot(containerId string) (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		logrus.Error("failed to get current user info: ", err)
		return "", err
	}
	rootDir := fmt.Sprintf("/sys/fs/cgroup/user.slice/user-%s.slice/user@%s.service/user.slice/minidocker-%s", currentUser.Uid, currentUser.Uid, containerId)
	if err := os.MkdirAll(rootDir, 0755); err != nil {
		CgroupsRoot = ""
		logrus.Error("failed to create cgroups: ", err)
		return "", err
	}
	CgroupsRoot = rootDir
	ContainerId = containerId
	return rootDir, nil
}

func RemoveCgroupsRoot(cgroupsRoot string) error {
	if err := os.Remove(cgroupsRoot); err != nil {
		logrus.Error("failed to remove cgroups: ", err)
		return err
	}
	return nil
}

func AddProcess(pid int) error {
	path := fmt.Sprintf("%s/cgroup.procs", CgroupsRoot)
	if err := os.WriteFile(path, []byte(strconv.Itoa(pid)), 0644); err != nil {
		logrus.Errorf("failed to add process %d to cgroup %s: %s\n", pid, CgroupsRoot, err.Error())
		return err
	}
	return nil
}

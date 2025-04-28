package cgroups

import (
	"bufio"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

var CgroupsDir string

//! 后续可能需要改为两个方案:
//! 1. 用systemd
//! 2. 通过参数指定，目前就做了这个，而且很临时。
//! 但是我希望就直接manager向systemd申请一个cgroup，给runtime用。
//! runtime后面就弄一下cgroup v1和v2
//! runtime尽量不要和别的组件耦合，只提供基础的功能接口

// func getCgroupDir() string {
// 	currentUser, _ := user.Current()
// 	return fmt.Sprintf("/sys/fs/cgroup/user.slice/user-%s.slice/user@%s.service/user.slice", currentUser.Uid, currentUser.Uid)
// }

// func YieldContainerId() string {
// 	// 组合输入数据(时间戳 + 随机盐)
// 	data := fmt.Sprintf("%d-%d", time.Now().UnixNano(), rand.Intn(1000))
// 	// 计算 SHA256 哈希
// 	hash := sha256.Sum256([]byte(data))
// 	// 转换为十六进制并截取前 12 位
// 	return hex.EncodeToString(hash[:])[:12]
// }

// 创建cgroup，如果成功的话会设置CgroupsRoot
// func CreateCgroupsRoot(containerId string) (string, error) {
// 	rootDir := fmt.Sprintf("%s/minidocker-%s", getCgroupDir(), containerId)
// 	if err := os.MkdirAll(rootDir, 0755); err != nil {
// 		CgroupsRoot = ""
// 		return "", fmt.Errorf("Failed to create cgroups: %s", err.Error())
// 	}
// 	CgroupsRoot = rootDir
// 	ContainerId = containerId
// 	return rootDir, nil
// }

// func RemoveCgroupsRoot(cgroupsRoot string) error {
// 	if err := os.Remove(cgroupsRoot); err != nil {
// 		logrus.Error("failed to remove cgroups: ", err)
// 		return err
// 	}
// 	return nil
// }

func AddProcess(pid int) error {
	path := path.Join(CgroupsDir, "cgroup.procs")
	if err := os.WriteFile(path, []byte(strconv.Itoa(pid)), 0644); err != nil {
		logrus.Errorf("failed to add process %d to cgroup %s: %s\n", pid, CgroupsDir, err.Error())
		return err
	}
	return nil
}

func GetCgroupsMounts() ([]string, error) {
	res := []string{}
	file, err := os.Open("/proc/self/mounts")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, " ")
		if len(fields) >= 3 && strings.HasPrefix(fields[2], "cgroup") {
			res = append(res, fields[1])
		}
	}
	return res, nil
}

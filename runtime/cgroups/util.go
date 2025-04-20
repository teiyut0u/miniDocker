package cgroups

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"os"
	"os/user"
	// "reflect"
	"regexp"
	"strconv"
	// "strings"
	"time"

	"github.com/sirupsen/logrus"
)

// type CgroupsField interface {
// 	Value() (string, error)
// 	SetValue(value string) error
// 	Remove() error
// }

var CgroupsRoot string

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

func SplitWords(x string) string {
	re := regexp.MustCompile(`([a-z])([A-Z])`)
	return re.ReplaceAllString(x, "${1}.${2}")
}

// func GetFieldName(field string) string {
// 	fieldVal := reflect.TypeOf(field)
// 	return fieldType.Name
// }

// func NewTName[T any](name string) *T {
// 	var res T
//
// 	resVal := reflect.ValueOf(&res).Elem()
// 	resType := resVal.Type()
// 	for i := range resVal.NumField() {
// 		fieldVal := resVal.Field(i)
// 		fieldType := resType.Field(i)
// 		value := fmt.Sprintf("%s.%s", name, strings.ToLower(SplitWords(fieldType.Name)))
// 		fieldVal.SetString(value)
// 	}
//
// 	return &res
// }

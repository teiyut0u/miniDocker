package cgroups

// import (
// 	"fmt"
// 	"os"
// 	"os/user"
// 	"strconv"
// 	"testing"
// )
//
// func Test_yieldContainerId(t *testing.T) {
//
// 	container := make(map[string]int)
// 	var containerNum int
// 	if n, err := strconv.Atoi(os.Getenv("CONTAINER_NUM")); err != nil {
// 		containerNum = 5
// 	} else {
// 		containerNum = n
// 	}
//
// 	for range containerNum {
// 		newId := YieldContainerId()
// 		if container[newId]++; container[newId] > 1 {
// 			t.Errorf("%d duplications detected\n", container[newId])
// 		}
// 	}
//
// }
//
// func Test_createCgroupRoot_and_removeCgroupsRoot(t *testing.T) {
// 	container := make(map[string]string)
// 	var containerNum int
// 	if n, err := strconv.Atoi(os.Getenv("CONTAINER_NUM")); err != nil {
// 		containerNum = 5
// 	} else {
// 		containerNum = n
// 	}
//
// 	for range containerNum {
// 		id := YieldContainerId()
// 		CreateCgroupsRoot(id)
//
// 		currentUser, err := user.Current()
// 		if err != nil {
// 			t.Error("tester failed to get current user info: ", err)
// 		}
// 		rootDir := fmt.Sprintf("/sys/fs/cgroup/user.slice/user-%s.slice/user@%s.service/user.slice/minidocker-%s", currentUser.Uid, currentUser.Uid, id)
// 		if _, err := os.Stat(rootDir); err != nil {
// 			t.Errorf("creation of directory %s is failed\n", rootDir)
// 		}
//
// 		container[id] = rootDir
// 	}
//
// 	for _, v := range container {
// 		RemoveCgroupsRoot(v)
// 		if _, err := os.Stat(v); !os.IsNotExist(err) {
// 			t.Errorf("removing of directory %s is failed\n", err)
// 		}
// 	}
//
// }

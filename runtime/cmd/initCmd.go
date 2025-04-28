package cmd

import (
	"encoding/gob"
	"miniDocker/runtime/namespace"
	"os"
	"os/exec"
	"strconv"

	"github.com/opencontainers/runtime-spec/specs-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// var initCmdConfig namespace.InitConfig

func init() {
	// initCmd.PersistentFlags().BoolVarP(&initCmdConfig.Interactive, "interactive", "i", false, "Keep STDIN and STDOUT open")
	rootCmd.AddCommand(initCmd)
}

func getSpecFromPipe(pipe *os.File) (*specs.Spec, error) {
	decoder := gob.NewDecoder(pipe)
	var res specs.Spec
	if err := decoder.Decode(&res); err != nil {
		return nil, err
	}
	return &res, nil
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Init namespace in a container",
	Run: func(cmd *cobra.Command, args []string) {
		// 获得用来传递参数的pipe，解析参数
		fd_n, err := strconv.Atoi(os.Getenv("_INIT_PIPE"))
		if err != nil {
			logrus.Error("failed to get the Init Pipe: ", err.Error())
		}
		initPipe := os.NewFile(uintptr(fd_n), "Init Pipe")
		specPtr, err := getSpecFromPipe(initPipe)
		if err != nil {
			logrus.Error("init process failed to get argv through the Init Pipe: ", err.Error())
			return
		}
		// 查找要执行的命令
		commandPath, err := exec.LookPath(specPtr.Process.Args[0])
		if err != nil {
			logrus.Error("failed to find the start command: ", err.Error())
			return
		} else {
			specPtr.Process.Args[0] = commandPath
		}
		// 运行init process
		if err := namespace.RunInitProcess(specPtr); err != nil {
			logrus.Error("failed to init process: ", err.Error())
		}
	},
}

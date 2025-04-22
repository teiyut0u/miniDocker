package cmd

import (
	"io"
	"miniDocker/runtime/namespace"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// var initCmdConfig namespace.InitConfig

func init() {
	// initCmd.PersistentFlags().BoolVarP(&initCmdConfig.Interactive, "interactive", "i", false, "Keep STDIN and STDOUT open")
	rootCmd.AddCommand(initCmd)
}

func getArgvFromPipe(pipe *os.File) ([]string, error) {
	args, err := io.ReadAll(pipe)
	if err != nil {
		return nil, err
	}
	return strings.Split(string(args), ":"), nil
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
		argv, err := getArgvFromPipe(initPipe)
		if err != nil {
			logrus.Error("init process failed to get argv through the Init Pipe: ", err.Error())
			return
		}
		// 查找要执行的命令
		commandPath, err := exec.LookPath(argv[0])
		if err != nil {
			logrus.Error("failed to find the start command: ", err.Error())
			return
		} else {
			argv[0] = commandPath
		}
		// 运行init process
		if err := namespace.RunInitProcess(argv[0], argv); err != nil {
			logrus.Error("failed to init process: ", err.Error())
		}
	},
}

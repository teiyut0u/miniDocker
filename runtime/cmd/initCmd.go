package cmd

import (
	"io"
	"miniDocker/runtime/namespace"
	"os"
	"os/exec"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// var initCmdConfig namespace.InitConfig

func init() {
	// initCmd.PersistentFlags().BoolVarP(&initCmdConfig.Interactive, "interactive", "i", false, "Keep STDIN and STDOUT open")
	rootCmd.AddCommand(initCmd)
}

func getArgsFromPipe(pipe *os.File) ([]string, error) {
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
		initPipe := os.NewFile(uintptr(3), "argsPipe")
		args, err := getArgsFromPipe(initPipe)
		if err != nil {
			logrus.Error("init process failed to get args through the Init Pipe: ", err.Error())
			return
		}
		// 查找要执行的命令
		commandPath, err := exec.LookPath(args[0])
		if err != nil {
			logrus.Error("failed to find the start command: ", err.Error())
			return
		} else {
			args[0] = commandPath
		}
		// 运行init process
		if err := namespace.RunInitProcess(args[0], args); err != nil {
			logrus.Error("failed to init process: ", err.Error())
		}
	},
}

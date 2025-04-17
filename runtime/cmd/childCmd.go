package cmd

import (
	"miniDocker/runtime/namespace"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var childCmdConfig namespace.InitConfig

func init() {
	childCmd.PersistentFlags().BoolVarP(&childCmdConfig.Interactive, "interactive", "i", false, "Keep STDIN and STDOUT open")
	rootCmd.AddCommand(childCmd)
}

var childCmd = &cobra.Command{
	Use:   "child [flags] [args...]",
	Short: "Fork child process as init process in a container",
	Run: func(cmd *cobra.Command, args []string) {
		command := namespace.YieldInitProcess(&childCmdConfig, args)
		if err := command.Run(); err != nil {
			logrus.Error("failed to launch child process: ", err.Error())
		}
		command.Wait()
	},
}

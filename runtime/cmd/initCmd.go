package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"miniDocker/runtime/namespace"
)

// var initCmdConfig namespace.InitConfig

func init() {
	// initCmd.PersistentFlags().BoolVarP(&initCmdConfig.Interactive, "interactive", "i", false, "Keep STDIN and STDOUT open")
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init [args...]",
	Short: "Init namespace in a container",
	Run: func(cmd *cobra.Command, args []string) {
		if err := namespace.InitProcess(args[0], args); err != nil {
			logrus.Error("failed to init process: ", err.Error())
		}
	},
}

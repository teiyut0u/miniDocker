package cmd

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var version bool

func init() {
	rootCmd.PersistentFlags().BoolVarP(&version, "version", "v", false, "Print version information and quit")
}

var rootCmd = &cobra.Command{
	Use:   "minidocker",
	Short: "minidocker-runtime is a container runtime demo",
	Run: func(cmd *cobra.Command, args []string) {
		if version {
			fmt.Println("minidocker-runtime version 0.0.0")
			return
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Error("failed to exec command: ", err.Error())
		os.Exit(1)
	}
}

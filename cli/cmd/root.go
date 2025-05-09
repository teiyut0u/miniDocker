package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var version bool

func init() {
	rootCmd.PersistentFlags().BoolVarP(&version, "version", "v", false, "Print version information and quit")
}

var rootCmd = &cobra.Command{
	Use:   "minidocker-cli",
	Short: "minidocker-cli is a container CLI demo",
	Run: func(cmd *cobra.Command, args []string) {
		if version {
			fmt.Println("minidocker-cli version 0.0.0")
			return
		}
		fmt.Println("run hugo...")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

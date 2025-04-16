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
	Use:   "miniDocker",
	Short: "miniDocker is a container demo",
	Run: func(cmd *cobra.Command, args []string) {
		if version {
			fmt.Println("miniDocker version 0.0.0")
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

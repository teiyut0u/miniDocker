package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var Interactive bool

func init() {
	runCmd.PersistentFlags().BoolVarP(&Interactive, "interactive", "i", false, "Keep STDIN and STDOUT open")
	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run program in a container",
	Run: func(cmd *cobra.Command, args []string) {
		if Interactive {
			fmt.Println("interactive mode")
		} else {
			fmt.Println("not interactive")
		}
		// notdone
	},
}

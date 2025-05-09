package cmd

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const Name = "minidocker-runtime"
const Version = "0.0.0"

// var ContainerId string

func init() {
	rootCmd.PersistentFlags().StringP("config", "", "", "path to config.json")
	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))

	rootCmd.PersistentFlags().BoolP("version", "v", false, "Print version information and quit")
	viper.BindPFlag("version", rootCmd.PersistentFlags().Lookup("version"))
}

var rootCmd = &cobra.Command{
	Use:   Name,
	Short: fmt.Sprintf("%s is a container runtime demo", Name),
	Run: func(cmd *cobra.Command, args []string) {
		if viper.GetBool("version") {
			fmt.Printf("%s version %s", Name, Version)
			return
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Errorf("failed to exec %s: %v", Name, err.Error())
		os.Exit(1)
	}
}

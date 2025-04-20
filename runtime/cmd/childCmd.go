package cmd

import (
	"miniDocker/runtime/cgroups"
	"miniDocker/runtime/cgroups/controllers"
	"miniDocker/runtime/namespace"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var childCmdConfig namespace.InitConfig
var memory controllers.Memory

func init() {
	childCmd.PersistentFlags().BoolVarP(&childCmdConfig.Interactive, "interactive", "i", false, "Keep STDIN and STDOUT open")
	childCmd.PersistentFlags().StringVarP(&memory.Max, "memory", "m", "", "Memory limit in bytes")
	// --cpuset-cpus="0,1" cpuset.cpus
	// --cpus cpu.max
	rootCmd.AddCommand(childCmd)
}

func setControllers() {
	// 这里好像需要处理错误
	cgroups.SetField("memory", &memory)
}

var childCmd = &cobra.Command{
	Use:   "child [flags] [args...]",
	Short: "Fork child process as init process in a container",
	Run: func(cmd *cobra.Command, args []string) {
		command := namespace.YieldInitProcess(&childCmdConfig, args)
		containerId := cgroups.YieldContainerId()
		cgroupsRoot, err := cgroups.CreateCgroupsRoot(containerId)
		if err != nil {
			return
		}
		defer cgroups.RemoveCgroupsRoot(cgroupsRoot)
		setControllers()

		if err := command.Start(); err != nil {
			logrus.Error("failed to launch child process: ", err.Error())
		}
		cgroups.AddProcess(command.Process.Pid)
		command.Wait()
	},
}

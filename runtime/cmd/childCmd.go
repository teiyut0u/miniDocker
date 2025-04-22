package cmd

import (
	"miniDocker/runtime/cgroups"
	"miniDocker/runtime/cgroups/controllers"
	"miniDocker/runtime/namespace"
	"strings"

	"github.com/opencontainers/runtime-spec/specs-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var childCmdConfig namespace.InitConfig
var memory controllers.Memory
var containerSpec specs.Spec

func init() {
	childCmd.PersistentFlags().BoolVarP(&childCmdConfig.Interactive, "interactive", "i", false, "Keep STDIN and STDOUT open")
	childCmd.PersistentFlags().StringVarP(&memory.Max, "memory", "m", "", "Memory limit in bytes")
	// childCmd.PersistentFlags().StringVarP(&(containerSpec.Root.Path), "rootfs", "", "", "Root fs directory")
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
		command, writeInitPipe := namespace.YieldInitProcess(&childCmdConfig)

		containerId := cgroups.YieldContainerId()
		cgroupsRoot, err := cgroups.CreateCgroupsRoot(containerId)
		if err != nil {
			logrus.Error("failed to create cgroups: ", err.Error())
			return
		}
		defer cgroups.RemoveCgroupsRoot(cgroupsRoot)
		setControllers()

		if err := command.Start(); err != nil {
			logrus.Error("failed to launch child process: ", err.Error())
		}
		cgroups.AddProcess(command.Process.Pid)
		writeInitPipe.WriteString(strings.Join(args, ":"))
		if err := writeInitPipe.Close(); err != nil {
			logrus.Error("failed to close Init Pipe: ", err.Error())
		}

		command.Wait()
	},
}

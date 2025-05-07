package cmd

import (
	"encoding/gob"
	"fmt"
	"miniDocker/runtime/cgroups"
	"miniDocker/runtime/cgroups/controllers"
	"miniDocker/runtime/namespace"
	"path"

	"github.com/opencontainers/runtime-spec/specs-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// var test string

func init() {
	childCmd.PersistentFlags().BoolP("interactive", "i", false, "Keep STDIN and STDOUT open")
	viper.BindPFlag("interactive", childCmd.PersistentFlags().Lookup("interactive"))

	childCmd.PersistentFlags().Int64P("memory", "m", 0, "Memory limit in bytes")
	viper.BindPFlag("memory", childCmd.PersistentFlags().Lookup("memory"))

	childCmd.PersistentFlags().String("id", "", "Container ID (required)")
	viper.BindPFlag("id", childCmd.PersistentFlags().Lookup("id"))
	// --cpuset-cpus="0,1" cpuset.cpus
	// --cpus cpu.max
	rootCmd.AddCommand(childCmd)
}

// 目前就设置了memory limit，以后的再说吧
func setControllers(resourcePtr *specs.LinuxResources) error {
	var memory controllers.Memory

	if resourcePtr != nil {
		if resourcePtr.Memory != nil {
			memory.SetFromConfig(resourcePtr.Memory)
		}
	}

	memory.SetFromCLI()

	return cgroups.SetField("memory", &memory)
}

func readConfig() error {
	if configPath := viper.GetString("config"); configPath != "" {
		viper.SetConfigFile(configPath)
	} else {
		// 设置配置文件的名称和路径
		viper.SetConfigName("config") // 配置文件名称（不包括扩展名）
		viper.SetConfigType("json")   // REQUIRED if the config file does not have the extension in the name
		viper.AddConfigPath(".")      // 配置文件的查找目录
	}

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("Fatal error config file: %w", err)
	}

	logrus.Infof("Using config file: %s", viper.ConfigFileUsed())
	return nil
}

func getInitSpec() (*specs.Spec, error) {
	if err := readConfig(); err != nil {
		return nil, err
	}

	var res specs.Spec
	res.Process = &specs.Process{}
	res.Root = &specs.Root{}
	res.Mounts = []specs.Mount{}
	res.Hooks = &specs.Hooks{}
	res.Annotations = make(map[string]string)
	res.Linux = &specs.Linux{}

	res.Version = viper.GetString("ociVersion")
	res.Hostname = viper.GetString("hostname")
	res.Domainname = viper.GetString("domainname") //?
	if err := viper.UnmarshalKey("process", res.Process); err != nil {
		return nil, err
	}
	if err := viper.UnmarshalKey("root", res.Root); err != nil {
		return nil, err
	}
	if err := viper.UnmarshalKey("mounts", &(res.Mounts)); err != nil {
		return nil, err
	}
	if err := viper.UnmarshalKey("hooks", res.Hooks); err != nil { //?
		return nil, err
	}
	if err := viper.UnmarshalKey("annotations", &(res.Annotations)); err != nil {
		return nil, err
	}
	if err := viper.UnmarshalKey("linux", res.Linux); err != nil {
		return nil, err
	}

	return &res, nil
}

func getInitConfig() *namespace.InitConfig {
	res := namespace.InitConfig{}
	res.Interactive = viper.GetBool("interactive")
	res.ContainerId = viper.GetString("id")
	logrus.Infof("id is %s", res.ContainerId)

	return &res
}

var childCmd = &cobra.Command{
	Use:   "child [flags] [args...]",
	Short: "Fork child process as init process in a container",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if viper.GetString("id") == "" {
			return fmt.Errorf("id is reuqired")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		// 获取spec
		spec, err := getInitSpec()
		if err != nil {
			logrus.Errorf("failed to read config.json: %v", err)
			return
		}
		// 获取其他配置
		childCmdConfig := getInitConfig()
		// 获取init cmd
		command, writeInitPipe := namespace.YieldInitProcess(childCmdConfig, spec)
		// 获取cgroup的挂载点，逻辑上可能有多个挂载点，先假设第一个是我们要的
		cgroupsMounts, err := cgroups.GetCgroupsMounts()
		if len(cgroupsMounts) == 0 || err != nil {
			logrus.Error("failed to get cgroups mounts: ", err.Error())
			return
		}
		// 获得该容器所使用的cgroup组
		cgroups.CgroupsDir = path.Join(cgroupsMounts[0], spec.Linux.CgroupsPath)
		// 设置cgroup
		var resourcePtr *specs.LinuxResources = nil
		if spec.Linux != nil {
			resourcePtr = spec.Linux.Resources
		}
		setControllers(resourcePtr)
		// 运行init进程
		if err := command.Start(); err != nil {
			logrus.Error("failed to launch child process: ", err.Error())
		}
		// 容器的init进程加入cgroups
		cgroups.AddProcess(command.Process.Pid)
		// spec通过管道写给init
		encoder := gob.NewEncoder(writeInitPipe)
		if err := encoder.Encode(*spec); err != nil {
			logrus.Errorf("failed to encode spec: %v", err)
		}
		if err := writeInitPipe.Close(); err != nil {
			logrus.Error("failed to close Init Pipe: ", err.Error())
		}
		// 这里要等子进程的，不然成孤儿了。detach用manager实现吧
		command.Wait()
	},
}

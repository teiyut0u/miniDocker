package config

import (
	"fmt"

	"github.com/opencontainers/runtime-spec/specs-go"
	"github.com/spf13/viper"
)

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

	return nil
}

func configSpec(spec *specs.Spec) {
	spec.Version = viper.GetString("ociVersion")
}

func configProcess() *specs.Process {
	if !viper.IsSet("process") {
		return nil
	}

	res := specs.Process{}

	res.Terminal = viper.GetBool("terminal")
	res.ConsoleSize = configConsoleSize()
	res.Args = viper.GetStringSlice("args")
	res.CommandLine = viper.GetString("commandLine")
	res.Env = viper.GetStringSlice("env")
	res.Cwd = viper.GetString("cwd")
	res.Capabilities = configCapabilities()
	res.Rlimits = configRlimits()
	res.NoNewPrivileges = viper.GetBool("noNewPrivileges")

	return &res
}

func configConsoleSize() *specs.Box {
	if !viper.IsSet("consoleSize") {
		return nil
	}
	res := specs.Box{}
	if err := viper.UnmarshalKey("consoleSize", &res); err != nil {
		return nil
	}
	return &res
}

func configCapabilities() *specs.LinuxCapabilities {
	if !viper.IsSet("capabilities") {
		return nil
	}

	res := specs.LinuxCapabilities{}

	res.Bounding = viper.GetStringSlice("bounding")
	res.Effective = viper.GetStringSlice("effective")
	res.Inheritable = viper.GetStringSlice("inheritable")
	res.Permitted = viper.GetStringSlice("permitted")
	res.Ambient = viper.GetStringSlice("ambient")

	return &res
}

func configRlimits() []specs.POSIXRlimit {
	if !viper.IsSet("rlimits") {
		return nil
	}

	res := []specs.POSIXRlimit{}

	if err := viper.UnmarshalKey("rlimits", &res); err != nil {
		return []specs.POSIXRlimit{}
	}
	return res
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

	// res.Version = viper.GetString("ociVersion")
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

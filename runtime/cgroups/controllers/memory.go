package controllers

import (
	"strconv"

	"github.com/opencontainers/runtime-spec/specs-go"
	"github.com/spf13/viper"
)

type Memory struct {
	Max string
}

func (mem *Memory) SetFromConfig(memoryPtr *specs.LinuxMemory) {
	if memoryPtr == nil {
		return
	}
	if memoryPtr.Limit != nil {
		mem.Max = strconv.FormatInt(*(memoryPtr.Limit), 10)
	}
}

func (mem *Memory) SetFromCLI() {
	if memLimit := viper.GetInt64("memory"); memLimit != 0 {
		mem.Max = strconv.FormatInt(memLimit, 10)
	}
}

package controllers

import "github.com/opencontainers/runtime-spec/specs-go"

type CPU struct {
	Stat       string
	Weight     string
	WeightNice string
	Max        string
}

func (cpu *CPU) SetFromConfig(cpuPtr *specs.LinuxCPU) {}

func (cpu *CPU) SetFromCLI() {}

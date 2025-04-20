package controllers

type CPU struct {
	Stat       string
	Weight     string
	WeightNice string
	Max        string
}

// func NewCPU() *CPU {
// 	return cgroups.NewT[CPU]("cpu")
// }

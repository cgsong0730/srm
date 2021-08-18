package main

import (
	"fmt"
	"os"
	"srm/lib/logger"
	"srm/workload/w1"
	"time"

	"github.com/containerd/cgroups"
	specs "github.com/opencontainers/runtime-spec/specs-go"
)

func genWorkload(goal int) {
	for true {
		w1.GeneratePrimeNumber(goal)
		//time.Sleep(1 * time.Second)
		fmt.Println("hello")
	}
}

func init() {
	err := logger.Init()
	if err != nil {
		logger.Fatal("Fail to initialize logger.")
		end()
	}
	logger.Info("#####[srm start]#####")
}

func main() {
	pid := os.Getpid()
	shares := uint64(100)
	var cpus string = "0-1"
	control, err := cgroups.New(cgroups.V1, cgroups.StaticPath("/cgs"), &specs.LinuxResources{
		CPU: &specs.LinuxCPU{
			Shares: &shares,
			Cpus:   cpus,
		},
	})
	if err != nil {
		logger.Fatal("Fail to create cgroup.")
	}

	if err := control.Add(cgroups.Process{Pid: pid}); err != nil {
		logger.Fatal("Fail to add cgroup.")
	}

	defer control.Delete()

	println("helloworld")
	//	genWorkload(1000000)
	go genWorkload(1000000)
	go genWorkload(1000000)

	time.Sleep(120 * time.Second)
	end()
}

func end() {
	logger.Info("#####[srm end]#####")
}

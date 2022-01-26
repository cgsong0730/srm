package controller

import (
	"srm/lib/logger"
	"strconv"

	"github.com/containerd/cgroups"
	"github.com/opencontainers/runtime-spec/specs-go"
)

/*
func init() {
	pid := os.Getpid()
	shares := uint64(100)
	var cpus string = "0-1"

	control, err := cgroups.New(cgroups.V1, cgroups.StaticPath("/cgs"),
		&specs.LinuxResources{
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
}
*/

func CreateResourcePolicy(pid int, cpus string) error {
	shares := uint64(100)

	control, err := cgroups.New(cgroups.V1, cgroups.StaticPath(strconv.Itoa(pid)),
		&specs.LinuxResources{
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

	return nil
}

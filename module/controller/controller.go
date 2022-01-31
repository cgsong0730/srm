package controller

import (
	"srm/lib/logger"
	"srm/lib/ptree"
	"strconv"

	"github.com/containerd/cgroups"
	"github.com/opencontainers/runtime-spec/specs-go"
)

func CreateResourcePolicy(node *ptree.Node, cpus string) (cgroups.Cgroup, error) {
	shares := uint64(100)

	control, err := cgroups.New(cgroups.V1, cgroups.StaticPath(strconv.Itoa(node.Pid)),
		&specs.LinuxResources{
			CPU: &specs.LinuxCPU{
				Shares: &shares,
				Cpus:   cpus,
			},
		})
	if err != nil {
		logger.Fatal("Fail to create cgroup.")
		return nil, err
	}

	return control, nil
}

func UpdateResourcePolicy(control cgroups.Cgroup, cpus string) error {
	shares := uint64(100)

	err := control.Update(
		&specs.LinuxResources{
			CPU: &specs.LinuxCPU{
				Shares: &shares,
				Cpus:   cpus,
			},
		})
	if err != nil {
		logger.Fatal("Fail to create cgroup.")
		return err
	}
	return nil
}

/*
func CreateResourcePolicy(node *ptree.Node, cpus string) error {
	shares := uint64(100)

	control, err := cgroups.New(cgroups.V1, cgroups.StaticPath(strconv.Itoa(node.Pid)),
		&specs.LinuxResources{
			CPU: &specs.LinuxCPU{
				Shares: &shares,
				Cpus:   cpus,
			},
		})
	if err != nil {
		logger.Fatal("Fail to create cgroup.")
	}

	if err := control.Add(cgroups.Process{Pid: node.Pid}); err != nil {
		logger.Fatal("Fail to add cgroup.")
	}

	for _, child := range node.Children {
		CreateResourcePolicy(child, cpus)
	}

	defer control.Delete()

	return nil
}
*/

func AddResourcePolicy(node *ptree.Node, control cgroups.Cgroup) error {

	if err := control.Add(cgroups.Process{Pid: node.Pid}); err != nil {
		logger.Fatal("Fail to add cgroup.")
		return err
	}

	for _, child := range node.Children {
		AddResourcePolicy(child, control)
	}

	return nil
}

/*
func DeleteResourcePolicy() error {
}
*/

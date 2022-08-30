package mape

import (
	"fmt"
	"os"
	"os/exec"
	config "srm/lib/config_parser"
	"srm/lib/ptree"
	"srm/module/analyzer"
	"srm/module/controller"
	"srm/module/monitor"
	"strconv"
	"time"

	"github.com/containerd/cgroups"
)

func Run() error {

	var containers []int
	var root ptree.Node

	pid := os.Getpid()
	println("[CSS] PID -> ", pid)

	root = ptree.Node{
		Pid: pid,
		Cnt: 0,
	}

	var ioContainerList []*ptree.Node
	var cpuContainerList []*ptree.Node
	var mapeCnt int = 0
	var useCleaning bool = false
	var useManagement bool = false
	var loopCnt int = 0

	var containerCgroup map[int]cgroups.Cgroup
	containerCgroup = make(map[int]cgroups.Cgroup)
	var oldContainerList []int

	numOfCpu := analyzer.GetCpuInfo()
	allOfCpu := "0-" + strconv.Itoa(numOfCpu-1)
	fmt.Println("cpu:", allOfCpu)

	for true {
		// M
		containers = monitor.FindContainer()
		for _, pid := range containers {
			ptree.CreateRootChild(&root, pid)
			monitor.GetChildTask(&root, pid)
		}

		cntOfDeletion := 0
		for index, pid := range oldContainerList {
			isContainer := false
			for _, cpid := range containers {
				if pid == cpid {
					isContainer = true
				}
			}
			if isContainer == false {
				root.Children = remove(root.Children, index-cntOfDeletion)
				cntOfDeletion++
			}
		}

		for _, node := range root.Children {
			isPid := false
			for _, pid := range oldContainerList {
				if node.Pid == pid {
					isPid = true
				}
			}

			if isPid == false {
				containerCgroup[node.Pid], _ = controller.CreateResourcePolicy(node, allOfCpu)
				controller.AddResourcePolicy(node, containerCgroup[node.Pid])
			}
		}

		oldContainerList = nil
		for _, node := range root.Children {
			oldContainerList = append(oldContainerList, node.Pid)
		}

		go monitor.GetSystemcall(&root, "futex")
		time.Sleep(time.Duration(config.Setting.Mape) * time.Second)

		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		if err := cmd.Run(); err != nil {
			fmt.Println(err)
		}
		ptree.PrintTree(&root, 0, false)

		/*
					if loopCnt%100 == 0 {
						logger.Info("print tree")
						ptree.LogTree(&root, 0, false)
			      ptree.CleanRootChild(&root)
					}
		*/
		loopCnt++

		// A
		for _, node := range root.Children {
			sum := ptree.SumContainerTree(node)
			if sum >= config.Setting.Threshold {
				ioContainerList = append(ioContainerList, node)
				//				useCleaning = true
				useManagement = true
			} else {
				cpuContainerList = append(cpuContainerList, node)
			}
		}

		// PE
		if len(cpuContainerList) != 0 {
			for _, node := range ioContainerList {
				controller.UpdateResourcePolicy(containerCgroup[node.Pid], config.Setting.Minimum)
			}
		} else {
			for i, node := range ioContainerList {
				numOfContainer := len(ioContainerList)
				if numOfCpu >= numOfContainer {
					part := numOfCpu / numOfContainer
					start := i * part
					cpus := "" + strconv.Itoa(start) + "-" + strconv.Itoa(start+part-1)
					controller.UpdateResourcePolicy(containerCgroup[node.Pid], cpus)
				}
			}
		}

		for i, node := range cpuContainerList {
			numOfContainer := len(cpuContainerList)
			if numOfCpu >= numOfContainer {
				part := numOfCpu / numOfContainer
				start := i * part
				cpus := "" + strconv.Itoa(start) + "-" + strconv.Itoa(start+part-1)
				controller.UpdateResourcePolicy(containerCgroup[node.Pid], cpus)
			}
		}

		if useCleaning == true {
			if mapeCnt == config.Setting.Clean-1 {
				ptree.CleanRootChild(&root)
				mapeCnt = 0
				useCleaning = false
			} else {
				mapeCnt += 1
			}
		}

		if useManagement == true {
			useManagement = false
		}

		ioContainerList = nil
		cpuContainerList = nil

	} // of MAPE-K

	return nil
}

func remove(s []*ptree.Node, i int) []*ptree.Node {
	return append(s[:i], s[i+1:]...)
}

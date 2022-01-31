package main

import (
	"fmt"
	"os"
	"os/exec"
	"srm/config"
	"srm/lib/logger"
	"srm/lib/ptree"
	"srm/module/analyzer"
	"srm/module/controller"
	"srm/module/monitor"
	"strconv"
	"time"

	"github.com/containerd/cgroups"
)

var containers []int
var root ptree.Node

func init() {
	err := logger.Init()
	if err != nil {
		logger.Fatal("Fail to initialize logger.")
		end()
	}
	logger.Info("#####[srm start]#####")

	pid := os.Getpid()
	println("[CSS] PID -> ", pid)

	root = ptree.Node{
		Pid: pid,
		Cnt: 0,
	}
}

func main() {

	// MAPE-K Loop
	var ioContainerList []*ptree.Node
	var cpuContainerList []*ptree.Node
	var mapeCnt int = 0
	var useCleaning bool = false
	var useManagement bool = false
	//var isCgroup bool = false
	//var ioControl cgroups.Cgroup

	var containerCgroup map[int]cgroups.Cgroup
	containerCgroup = make(map[int]cgroups.Cgroup)
	var oldContainerList []int

	for true {

		// M
		containers = monitor.FindContainer()
		for _, pid := range containers {
			ptree.CreateRootChild(&root, pid)
			monitor.GetChildTask(&root, pid)
		}

		for _, node := range root.Children {
			isPid := false
			for _, pid := range oldContainerList {
				if node.Pid == pid {
					isPid = true
				}
			}

			if isPid == false {
				containerCgroup[node.Pid], _ = controller.CreateResourcePolicy(node, "0-3")
				controller.AddResourcePolicy(node, containerCgroup[node.Pid])
				oldContainerList = append(oldContainerList, node.Pid)
			}
		}

		go monitor.GetSystemcall(&root, "mmap")
		time.Sleep(time.Duration(config.Interval) * time.Second)

		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		if err := cmd.Run(); err != nil {
			fmt.Println(err)
		}
		ptree.PrintTree(&root, 0, false)

		// A
		for _, node := range root.Children {
			sum := ptree.SumContainerTree(node)
			fmt.Println("Singularity -> pid:", node.Pid, ", sum: ", sum)
			if sum >= config.IoThresholdValue {
				ioContainerList = append(ioContainerList, node)
				useCleaning = true
				useManagement = true
			} else {
				cpuContainerList = append(cpuContainerList, node)
			}
		}

		// PE
		for _, node := range ioContainerList {
			fmt.Println("io-intensive: ", node.Pid)
			controller.UpdateResourcePolicy(containerCgroup[node.Pid], config.MCpus)
		}

		for i, node := range cpuContainerList {

			fmt.Println("cpu-intensive: ", node.Pid, ", len: ", len(cpuContainerList))

			numOfCpu := analyzer.GetCpuInfo()
			numOfContainer := len(cpuContainerList)

			if numOfCpu >= numOfContainer {
				part := numOfCpu / numOfContainer
				start := i * part
				cpus := "" + strconv.Itoa(i*start) + "-" + strconv.Itoa(start+part-1)
				fmt.Println("cpus: ", cpus)
				controller.UpdateResourcePolicy(containerCgroup[node.Pid], cpus)
			}
		}

		if useManagement == true {
			useManagement = false
		}

		if useCleaning == true {
			if mapeCnt == config.CleaningInterval-1 {
				ptree.CleanRootChild(&root)
				mapeCnt = 0
				useCleaning = false
			} else {
				mapeCnt += 1
			}
		}
		ioContainerList = nil
		cpuContainerList = nil
		root.Children = nil
	} // of MAPE-K

	end()
}

func end() {
	logger.Info("#####[srm end]#####")
}

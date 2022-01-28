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
	_ "srm/module/monitor"
	"strconv"
	"time"
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

	root = ptree.Node{
		Pid: 0,
		Cnt: 0,
	}

	pid := os.Getpid()
	println("[CSS] PID -> ", pid)
}

func main() {

	// MAPE-K Loop
	var containerNodeList []*ptree.Node
	var ioContainerList []*ptree.Node
	var cpuContainerList []*ptree.Node
	var mapeCnt int = 0
	var useCleaning bool = false
	var useManagement bool = false
	for true {
		// go monitor.GetSystemcall("clone")
		// go monitor.GetSystemcall("fork")

		// M
		containerNodeList = nil
		containers = monitor.FindContainer()
		for _, pid := range containers {

			monitor.GetChildTask(&root, pid)
			ptree.CreateRootChild(&root, pid)
			for _, node := range root.Children {
				if pid == node.Pid {
					containerNodeList = append(containerNodeList, node)
				}
			}
		}
		root.Children = containerNodeList

		go monitor.GetSystemcall(&root, "mmap")
		time.Sleep(time.Duration(config.Interval) * time.Second)

		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		if err := cmd.Run(); err != nil {
			fmt.Println(err)
		}
		ptree.PrintTree(&root, 0)

		// A
		for _, node := range root.Children {
			sum := ptree.SumTree(node)
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
			controller.CreateResourcePolicy(node, config.MCpus)
		}

		for i, node := range cpuContainerList {

			fmt.Println("cpu-intensive: ", node.Pid, ", len: ", len(cpuContainerList))

			numOfCpu := analyzer.GetCpuInfo()
			numOfContainer := len(cpuContainerList)

			if numOfCpu >= numOfContainer {
				part := numOfCpu / numOfContainer
				start := i * part
				//fmt.Println("", i*start, "-", start+part-1)
				cpus := "" + strconv.Itoa(i*start) + "-" + strconv.Itoa(start+part-1)
				fmt.Println("cpus: ", cpus)

				controller.CreateResourcePolicy(node, cpus)
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
	} // of MAPE-K

	end()
}

func end() {
	logger.Info("#####[srm end]#####")
}

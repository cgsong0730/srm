package main

import (
	"fmt"
	"os"
	"os/exec"
	config "srm/lib/config_parser"
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

	err = config.Init()
	if err != nil {
		logger.Fatal("Fail to initialize config_parser")
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
	//	var mapeCnt int = 0
	//  var useCleaning bool = false
	var useManagement bool = false

	var containerCgroup map[int]cgroups.Cgroup
	containerCgroup = make(map[int]cgroups.Cgroup)
	var oldContainerList []int

	numOfCpu := analyzer.GetCpuInfo()

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
				ptree.LogTree(&root, 0, false)
				root.Children = remove(root.Children, index-cntOfDeletion)
				//fmt.Println("deleted container:", pid, ", index:", index)
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
				containerCgroup[node.Pid], _ = controller.CreateResourcePolicy(node, "0-15")
				controller.AddResourcePolicy(node, containerCgroup[node.Pid])
			}
		}

		oldContainerList = nil
		for _, node := range root.Children {
			oldContainerList = append(oldContainerList, node.Pid)
		}

		go monitor.GetSystemcall(&root, "futex")
		time.Sleep(time.Duration(config.Setting.Mape) * time.Second)
		//fmt.Println("mape:", config.Setting.Mape)
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		if err := cmd.Run(); err != nil {
			fmt.Println(err)
		}
		ptree.PrintTree(&root, 0, false)

		// A
		for _, node := range root.Children {
			sum := ptree.SumContainerTree(node)
			//fmt.Println("Singularity -> pid:", node.Pid, ", sum: ", sum)
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
				//fmt.Println("io-intensive: ", node.Pid)
				controller.UpdateResourcePolicy(containerCgroup[node.Pid], config.Setting.Minimum)
			}
		} else {
			for i, node := range ioContainerList {
				//fmt.Println("io-intensive: ", node.Pid)

				//numOfCpu := analyzer.GetCpuInfo()
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

			//fmt.Println("cpu-intensive: ", node.Pid, ", len: ", len(cpuContainerList))

			//numOfCpu := analyzer.GetCpuInfo()
			numOfContainer := len(cpuContainerList)

			if numOfCpu >= numOfContainer {
				part := numOfCpu / numOfContainer
				start := i * part
				cpus := "" + strconv.Itoa(start) + "-" + strconv.Itoa(start+part-1)
				//fmt.Println("cpus: ", cpus)
				controller.UpdateResourcePolicy(containerCgroup[node.Pid], cpus)
			}
		}

		if useManagement == true {
			useManagement = false
		}

		/*
			if useCleaning == true {
					if mapeCnt == config.Setting.Clean-1 {
						//fmt.Println("clean start")
						ptree.CleanRootChild(&root)
						mapeCnt = 0
						useCleaning = false
					} else {
						mapeCnt += 1
					}
					//fmt.Println("clean:", config.Setting.Clean, ", mapeCnt:", mapeCnt)
			}
		*/

		ioContainerList = nil
		cpuContainerList = nil

	} // of MAPE-K

	end()
}

/*
func remove(s []*ptree.Node, i int) []*ptree.Node {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
*/
/*
func remove(slice []*ptree.Node, i int) []*ptree.Node {
	ret := make([]*ptree.Node, 0)
	if len(slice)-1 == i { // when s is last index
		return append(ret, slice[:i]...)
	} else {
		ret = append(ret, slice[:i]...)
		return append(ret, slice[i+1:]...)
	}
}
*/

func remove(s []*ptree.Node, i int) []*ptree.Node {
	return append(s[:i], s[i+1:]...)
}

func end() {
	logger.Info("#####[srm end]#####")
}

package main

import (
	"fmt"
	"os"
	"os/exec"
	"srm/lib/logger"
	"srm/lib/process_tree"
	"srm/workload/w1"
	"strconv"
	"strings"
	"time"

	"github.com/containerd/cgroups"
	specs "github.com/opencontainers/runtime-spec/specs-go"
)

/*
type Node struct {
	pid      int
	cnt      int
	children []*Node
}
*/

var rootPid int

// breadth-first search
/*
func findById(root *Node, pid int) *Node {
	queue := make([]*Node, 0)
	queue = append(queue, root)
	for len(queue) > 0 {
		nextUp := queue[0]
		queue = queue[1:]
		if nextUp.pid == pid {
			return nextUp
		}
		if len(nextUp.children) > 0 {
			for _, child := range nextUp.children {
				queue = append(queue, child)
			}
		}
	}
	return nil
}
*/

func getChildNode(pid int) {
	cmd := exec.Command("pgrep", "-P", strconv.Itoa(pid))
	output, err := cmd.Output()
	if err != nil {
		logger.Fatal("Fail to find child nodes.")
	}
	println("child: ", string(output))
}

func getSystemcall(root *process_tree.Node, systemcall string) {

	var pid int
	cmd := exec.Command("bpftrace", systemcall+".bt")
	output, err := cmd.Output()
	if err != nil {
		logger.Fatal("Fail to execute bpftrace command.")
	}

	str := strings.Split(string(output), "\n")
	pids := strings.Split(str[1], ",")
	pids = pids[:len(pids)-1]

	for _, str := range pids {
		pid, _ = strconv.Atoi(str)
		if process_tree.FindById(root, pid) != nil {
			fmt.Println(systemcall, ":", str)
			getChildNode(pid)
		}
	}
}

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

	p1 := process_tree.Node{
		pid: 1,
		cnt: 5,
	}

	p2 := process_tree.Node{
		pid: 2,
		cnt: 10,
	}

	rootPid, _ = strconv.Atoi(os.Args[1])
	rootNode := process_tree.Node{
		pid:      rootPid,
		cnt:      1,
		children: []*process_tree.Node{&p2},
	}
	rootNode.children = append(rootNode.children, &p1)

	//println("cnt:", px.cnt)
	//println("container shell pid:", os.Args[1])
	println("root process's pid:", rootNode.pid)

	getChildNode(rootPid)

	println("[CSS] Start", len(rootNode.children))

	for true {
		go getSystemcall(&rootNode, "clone")
		go getSystemcall(&rootNode, "mmap")
		go getSystemcall(&rootNode, "fork")
		time.Sleep(10 * time.Second)
	}

	end()
}

func end() {
	logger.Info("#####[srm end]#####")
}

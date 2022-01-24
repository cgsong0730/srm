package main

import (
	"fmt"
	"os"
	"os/exec"
	"srm/config"
	"srm/lib/logger"
	"srm/lib/ptree"
	"srm/module/monitor"
	_ "srm/module/monitor"
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
		Pid: 1,
		Cnt: 0,
	}

	/*
		     p1 := ptree.Node{
		         Pid: 2,
		         Cnt: 0,
		     }

			   p2 := ptree.Node{
			       Pid: 3,
			       Cnt: 0,
			   }

			   ptree.AddChild(&root, 1, &p1)
			   ptree.AddChild(&root, 1, &p2)

			   ptree.CreateChild(&root, 2, 4)
			   ptree.CreateChild(&root, 2, 5)

			   ptree.PlusCount(&root, 4)
			   ptree.DeleteChild(&root, 2)
	*/
	pid := os.Getpid()
	//  rootNode = ptree.CreateRootNode(pid)
	println("[CSS] PID -> ", pid)
	containers = monitor.FindContainer()
	for _, pid := range containers {
		monitor.GetChildTask(pid)
		ptree.CreateChild(&root, 1, pid)
	}

	ptree.PrintTree(&root, 0)
}

func main() {
	/*
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
	*/

	// MAPE-K Loop
	for true {
		//      go monitor.GetSystemcall("clone")
		go monitor.GetSystemcall(&root, "mmap")
		//      go monitor.GetSystemcall("fork")

		time.Sleep(time.Duration(config.Interval) * time.Second)

		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		if err := cmd.Run(); err != nil {
			fmt.Println(err)
		}
		ptree.PrintTree(&root, 0)
	}

	end()
}

func end() {
	logger.Info("#####[srm end]#####")
}

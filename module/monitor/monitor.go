package monitor

import (
	"fmt"
	"os/exec"
	"srm/lib/logger"
	"strconv"
	"strings"
)

type Node struct {
	Pid      int
	Cnt      int
	Children []*Node
}

func GetSystemcall(root *Node, systemcall string) {

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
		if findById(root, pid) != nil {
			fmt.Println(systemcall, ":", str)
			getChildNode(pid)
		}
	}
}

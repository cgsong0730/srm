package monitor

import (
    "fmt"
    "os/exec"
    "srm/lib/logger"
    "strconv"
    "strings"
)

type Node struct {
    pid      int
    cnt      int
    children []*Node
}

var rootPid int
var rootNode Node

/*
func Init() {

    p1 := Node{
        pid: 1,
        cnt: 5,
    }

    p2 := Node{
        pid: 2,
        cnt: 10,
    }

    rootPid, _ = strconv.Atoi(os.Args[1])
    rootNode = Node{
        pid:      rootPid,
        cnt:      1,
        children: []*Node{&p2},
    }
    rootNode.children = append(rootNode.children, &p1)

    //println("cnt:", px.cnt)
    //println("container shell pid:", os.Args[1])
    println("root process's pid:", rootNode.pid)

    getChildNode(rootPid)
}

// breadth-first search
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

func FindContainer() []int {

    var containers []int

    cmd := exec.Command("./fc.sh")
    output, err := cmd.Output()
    if err != nil {
        logger.Fatal("Fail to find namespace.")
    }
    fmt.Println(string(output))

    pids := strings.Split(string(output), "\n")
    pids = pids[:len(pids)-1]

    for _, str := range pids {
        pid, _ := strconv.Atoi(str)
        fmt.Println("Singularity -> ", pid)
        containers = append(containers, pid)
    }

    return containers
}

func GetChildTask(pid int) {
    cmd := exec.Command("pgrep", "-P", strconv.Itoa(pid))
    output, err := cmd.Output()
    if err != nil {
        logger.Fatal("Fail to find child nodes.")
    }
    pids := strings.Split(string(output), "\n")
    pids = pids[:len(pids)-1]

    for _, str := range pids {
        println("child: ", str)
    }
}

func GetSystemcall(systemcall string) {

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
        /*
            if findById(&rootNode, pid) != nil {
                fmt.Println(systemcall, ":", str)
                getChildNode(pid)
            }
        */
        fmt.Println(pid, "->", systemcall)
    }
}


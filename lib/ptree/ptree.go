package ptree

import "fmt"

type Node struct {
    Pid      int
    Cnt      int
    Children []*Node
}

func CreateChild(root *Node, ppid int, cpid int) {
    parent := FindByPid(root, ppid)
    child := Node{
        Pid: cpid,
        Cnt: 0,
    }
    parent.Children = append(parent.Children, &child)
}

func AddChild(root *Node, ppid int, child *Node) {
    parent := FindByPid(root, ppid)
    parent.Children = append(parent.Children, child)
}

func DeleteChild(root *Node, pid int) {
    target := FindByPid(root, pid)
    target.Children = nil
}

func FindByPid(root *Node, pid int) *Node {
    if root.Pid == pid {
        return root
    } else if len(root.Children) > 0 {
        for _, child := range root.Children {
            tmp := FindByPid(child, pid)
            if tmp != nil {
                return tmp
            }
        }
    }
    return nil
}

func PrintTree(root *Node, depth int) {
    if depth != 0 {
        for i := 1; i <= depth; i++ {
            fmt.Printf("   ")
        }
        fmt.Printf("└─")
    } else {
        fmt.Printf("  ")
    }
    fmt.Printf("(%d)[%d]\n", root.Pid, root.Cnt)
    if len(root.Children) > 0 {
        //fmt.Println("# of children:", len(root.children))
        for _, child := range root.Children {
            PrintTree(child, depth+1)
        }
    }
}

func SumTree(root *Node) int {
    sum := root.Cnt
    if len(root.Children) > 0 {
        for _, child := range root.Children {
            sum += SumTree(child)
        }
    }
    return sum
}

func PlusCount(root *Node, pid int) {
    tmp := FindByPid(root, pid)
    tmp.Cnt += 1
}

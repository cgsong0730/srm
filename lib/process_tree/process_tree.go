package process_tree

type Node struct {
	pid      int
	cnt      int
	children []*Node
}

// breath-first search
func FindById(root *Node, pid int) *Node {
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

func FreateRootNode(pid int) *Node {
	rootNode := Node{
		pid: pid,
		cnt: 0,
	}
	return &rootNode
}

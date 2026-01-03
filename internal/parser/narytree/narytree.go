package narytree

type Node struct {
	Type     string
	Data     string
	Children []Node
}

func CreateNewNode(data string) Node {
	node := Node{
		Data:     data,
		Children: []Node{},
	}

	return node
}

func (n *Node) AddChild(node Node) {
	n.Children = append(n.Children, node)
}

func (n *Node) PrintTree() {
	n.printTreeHelper(0)
}

func (n *Node) printTreeHelper(indent int) {
	for i := 0; i < indent; i++ {
		print("  ")
	}
	if n.Type != "" {
		println("[" + n.Type + "] " + n.Data)
	} else {
		println(n.Data)
	}
	for i := range n.Children {
		n.Children[i].printTreeHelper(indent + 1)
	}
}
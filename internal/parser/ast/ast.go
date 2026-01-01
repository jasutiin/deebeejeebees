package ast

type ASTNode struct {
	Data     string
	Children []ASTNode
}

func CreateNewNode(data string) ASTNode {
	node := ASTNode{
		Data:     data,
		Children: []ASTNode{},
	}

	return node
}

func (n *ASTNode) AddChild(node ASTNode) {
	n.Children = append(n.Children, node)
}

func (n *ASTNode) PrintTree() {
	n.printTreeHelper(0)
}

func (n *ASTNode) printTreeHelper(indent int) {
	for i := 0; i < indent; i++ {
		print("  ")
	}
	println(n.Data)
	for i := range n.Children {
		n.Children[i].printTreeHelper(indent + 1)
	}
}
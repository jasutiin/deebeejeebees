package ast

type ASTNode struct {
	data     string
	children []ASTNode
}

func CreateNewNode(data string) ASTNode {
	node := ASTNode{
		data:     data,
		children: []ASTNode{},
	}

	return node
}

func (n *ASTNode) AddChild(node ASTNode) {
	n.children = append(n.children, node)
}
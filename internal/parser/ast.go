package parser

import (
	"github.com/jasutiin/deebeejeebees/internal/parser/narytree"
)

const (
	IdentifierNode      = "IdentifierNode"
	ColumnListNode      = "ColumnListNode"
	ValuesListNode      = "ValuesListNode"
	DataTypeNode        = "DataTypeNode"
	ConstraintNode      = "ConstraintNode"
	ColumnDefNode       = "ColumnDefNode"
	Operator            = "Operator"
	Left                = "Left"
	Right               = "Right"
	BinaryOperationNode = "BinaryOperationNode"
	ValueNode           = "ValueNode"
	TableNameNode       = "TableNameNode"
	SelectNode          = "SelectNode"
	CreateTableNode     = "CreateTableNode"
	InsertNode          = "InsertNode"
)

var transformationRules = map[string]string{
	"<column_list>":           ColumnListNode,
	"<table_name>":            TableNameNode,
	"<optional_where>":        BinaryOperationNode,
	"<value_list>":            ValuesListNode,
	"<column_defs_list>":      ColumnListNode,
	"<column_def>":            ColumnDefNode,
	"<data_type>":             DataTypeNode,

	"<column_name>":           "DEL",
	"<column_list_tail>":      "DEL",
	"<condition>":             "DEL",
	"<value>":                 "DEL",
	"<value_list_tail>":       "DEL",
	"<column_defs_list_tail>": "DEL",

	"SELECT": "DEL",
	"FROM":   "DEL",
	"WHERE":  "DEL",
	"INSERT": "DEL",
	"INTO":   "DEL",
	"VALUES": "DEL",
	"CREATE": "DEL",
	"TABLE":  "DEL",
	",":      "DEL",
	";":      "DEL",
	"(":      "DEL",
	")":      "DEL",
}

func ConvertToAST(root narytree.Node) narytree.Node {
	astRoot := root
	astRoot.Type = determineQueryType(&astRoot) // set root node based on type
	astRoot.Data = ""

	i := 0
	for i < len(astRoot.Children) {
		if rule := transformationRules[astRoot.Children[i].Data]; rule == "DEL" {
			astRoot.Children = append(astRoot.Children[:i], astRoot.Children[i+1:]...) // remove child
			continue
		}

		// if we don't have '&', we get a copy of the node and do the transformation on the copy.
		// this means that the original node inside the slice is left unchanged.
		node := &astRoot.Children[i]
		transformNode(node)
		i++
	}

	return astRoot
}

func determineQueryType(node *narytree.Node) string {
	for _, child := range node.Children {
		switch child.Data {
			case "SELECT":
				return SelectNode
			case "CREATE":
				return CreateTableNode
			case "INSERT":
				return InsertNode
		}
	}
	return ""
}

func transformNode(node *narytree.Node) {
	ruleName := transformationRules[node.Data]

	if ruleName == "DEL" {
		return
	}

	// handle table name node
	if ruleName == TableNameNode {
		node.Type = TableNameNode
		if len(node.Children) > 0 {
			node.Data = node.Children[0].Data
			node.Children = nil
		}
		return
	}

	// handle binary operation node (WHERE clause)
	if ruleName == BinaryOperationNode {
		node.Type = BinaryOperationNode
		node.Data = "WhereClause"
		identifiers := collectIdentifiers(node)

		// label the children as Left, Operator, Right
		var newChildren []narytree.Node
		for i, child := range identifiers {
			switch i {
				case 0:
					child.Type = Left
				case 1:
					child.Type = Operator
				case 2:
					child.Type = Right
			}

			newChildren = append(newChildren, child)
		}

		node.Children = newChildren
		return
	}

	identifiers := collectIdentifiers(node)
	node.Children = identifiers

	if ruleName != "" {
		node.Type = ruleName
		node.Data = ""
	}
}

func collectIdentifiers(currentNode *narytree.Node) []narytree.Node {
	var result []narytree.Node

	for _, child := range currentNode.Children {
		ruleName := transformationRules[child.Data]

		switch ruleName {
			case "":
				child.Type = IdentifierNode
				result = append(result, child)
				
			case ColumnDefNode:
				child.Type = ColumnDefNode
				child.Data = ""
				var newChildren []narytree.Node

				for i, grandchild := range child.Children {
					grandchildRule := transformationRules[grandchild.Data]

					if i == 0 { // this is the name of the column
						grandchild.Type = IdentifierNode
						newChildren = append(newChildren, grandchild)
					} else if grandchildRule == DataTypeNode {
						grandchild.Type = DataTypeNode
						processDataType(&grandchild)
						// move ValueNode children up to be siblings of DataTypeNode
						valueNodes := grandchild.Children
						grandchild.Children = nil
						newChildren = append(newChildren, grandchild)
						newChildren = append(newChildren, valueNodes...)
					} else if grandchildRule == ConstraintNode {
						grandchild.Type = ConstraintNode
						newChildren = append(newChildren, grandchild)
					}
				}

				child.Children = newChildren
				result = append(result, child)

			case DataTypeNode:
				child.Type = DataTypeNode
				processDataType(&child)
				result = append(result, child)

			case ConstraintNode:
				child.Type = ConstraintNode
				result = append(result, child)

			default:
				result = append(result, collectIdentifiers(&child)...)
		}
	}

	return result
}

func processDataType(node *narytree.Node) {
	if len(node.Children) == 0 {
		return
	}

	node.Data = node.Children[0].Data // this is the data type name like varchar or int

	var newChildren []narytree.Node
	for i := 1; i < len(node.Children); i++ {
		child := node.Children[i]

		if child.Data == "(" || child.Data == ")" { // skip parentheses
			continue
		}

		child.Type = ValueNode
		newChildren = append(newChildren, child)
	}
	node.Children = newChildren
}
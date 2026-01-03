package parser

import (
	"github.com/jasutiin/deebeejeebees/internal/parser/narytree"
)

var transformationRules = map[string]string{
	"<column_list>": "Projection",
	"<table_name>": "Source",
	"<optional_where>": "WhereClause",
	"<value_list>": "InsertValues",
	"<column_defs_list>": "ColumnDefinitions",

	"<column_name>": "DEL", 
	"<column_list_tail>": "DEL",
	"<condition>": "DEL",
	"<value>": "DEL",
	"<value_list_tail>": "DEL",
	"<column_def>": "DEL",
	"<data_type>": "DEL",
	"<column_defs_list_tail>": "DEL",
	
	"SELECT": "DEL",
	"FROM": "DEL",
	"WHERE": "DEL",
	",": "DEL",
	";": "DEL",
}

func ConvertToAST(root narytree.Node) narytree.Node {
	queryType := root.Children[0].Data
	astRoot := root

	switch queryType {
		case "SELECT":
			for i := 0; i < len(astRoot.Children); i++ {
				// if we don't have '&', we get a copy of the node and do the transformation on the copy.
				// this means that the original node inside the slice is left unchanged.
				node := &astRoot.Children[i]
				transformNode(node)
			}
	}

	return astRoot
}

func transformNode(node *narytree.Node) {
	identifiers := collectIdentifiers(node)
	node.Children = identifiers

	ruleName := transformationRules[node.Data]
	
	if ruleName == "DEL" {
		node = nil
	}
	
	if ruleName != "" {
		node.Data = ruleName
	}
}

func collectIdentifiers(currentNode *narytree.Node) []narytree.Node {
	var result []narytree.Node
	
	for _, child := range currentNode.Children {
		ruleName := transformationRules[child.Data]
		if ruleName == "" {
			result = append(result, child)
		} else {
			result = append(result, collectIdentifiers(&child)...)
		}
	}

	return result
}
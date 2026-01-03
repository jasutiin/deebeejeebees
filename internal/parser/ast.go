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
	"INSERT": "DEL",
	"INTO": "DEL",
	"VALUES": "DEL",
	"CREATE": "DEL",
	"TABLE": "DEL",		
	",": "DEL",
	";": "DEL",
	"(": "DEL",
	")": "DEL",
}

func ConvertToAST(root narytree.Node) narytree.Node {
	astRoot := root

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

func transformNode(node *narytree.Node) {
	identifiers := collectIdentifiers(node)
	node.Children = identifiers

	ruleName := transformationRules[node.Data]
	
	if ruleName == "DEL" {
		return
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
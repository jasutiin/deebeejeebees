package parser

import (
	"errors"
	"fmt"

	"github.com/jasutiin/deebeejeebees/internal/parser/ast"
)

type Parser struct {
	tokens []string
	pos int
	rootNode *ast.ASTNode
}

func ParseTokens(tokens []string) ast.ASTNode {
	rootNode := ast.ASTNode{ Data: "<query>", Children: []ast.ASTNode{} }
	parser := Parser{ tokens: tokens, pos: 0, rootNode: &rootNode }
	queryType := tokens[parser.pos]

	switch queryType {
		case "SELECT":
			err := parser.parseSelect() // parsing select statements

			if err != nil {
				fmt.Println(err.Error())
			}

		case "INSERT":
			insertNode, err := parser.parseInsert() // parsing insert statements

			if err != nil {
				fmt.Println(err.Error())
			}

			rootNode.AddChild(insertNode)
		case "CREATE":
			createNode, err := parser.parseCreate() // parsing create statements

			if err != nil {
				fmt.Println(err.Error())
			}

			rootNode.AddChild(createNode)
	}

	return rootNode
}

// parseSelect is responsible for parsing the SELECT query type
func (p *Parser) parseSelect() error {
	selectNode := ast.ASTNode{ Data: "SELECT", Children: []ast.ASTNode{} }
	colListNode, err := p.parseColumnList() // should return a whole branch
	
	if err != nil {
		return err
	}

	p.rootNode.AddChild(selectNode)
	p.rootNode.AddChild(colListNode)
	// tableListNode := p.parseTableList()
	return nil
}

func (p *Parser) parseColumnList() (ast.ASTNode, error) {
	nextToken := p.peek()

	if nextToken == "FROM" {
		return ast.ASTNode{}, errors.New("missing at least one column name after SELECT!")
	}

	columnListNode := ast.ASTNode { Data: "<column_list>", Children: []ast.ASTNode{} }
	p.incrementPosition()

	columnName := p.parseColumnName()
	columnListNode.AddChild(columnName)
	p.parseColumnListTail(&columnListNode)

	return columnListNode, nil
}

func (p *Parser) parseColumnName() ast.ASTNode {
	columnNameNonTerminal := ast.ASTNode{ Data: "<column_name>", Children: []ast.ASTNode{} }
	columnName := ast.ASTNode{ Data: p.tokens[p.pos], Children: []ast.ASTNode{} }
	columnNameNonTerminal.Children = append(columnNameNonTerminal.Children, columnName)
	return columnNameNonTerminal
}

func (p *Parser) parseColumnListTail(parentNode *ast.ASTNode) ast.ASTNode {
	columnListTailNode := ast.ASTNode{ Data: "<column_list_tail>", Children: []ast.ASTNode{} }
	nextToken := p.peek()
	
	if nextToken == "," {
		p.incrementPosition()
		commaNode := ast.ASTNode{ Data: p.tokens[p.pos], Children: []ast.ASTNode{} }
		columnListTailNode.AddChild(commaNode)
		p.incrementPosition()
		columnName := p.parseColumnName()
		columnListTailNode.AddChild(columnName)
		p.parseColumnListTail(&columnListTailNode)
	}
	
	parentNode.AddChild(columnListTailNode)
	return columnListTailNode
}

// parseInsert is responsible for parsing the SELECT query type
func (p *Parser) parseInsert() (ast.ASTNode, error) {
	node := ast.ASTNode{}
	return node, nil
}

// parseCreate is responsible for parsing the SELECT query type
func (p *Parser) parseCreate() (ast.ASTNode, error) {
	node := ast.ASTNode{}
	return node, nil
}

func (p *Parser) incrementPosition() {
	p.pos += 1
}

func (p *Parser) peek() string {
	return p.tokens[p.pos + 1]
}
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
	fromNode := p.parseFromNode()
	tableNameNode := p.parseTableName()
	optionalWhereNode := p.parseOptionalWhere()
	semicolonNode := p.parseSemicolon()
	
	if err != nil {
		return err
	}

	p.rootNode.AddChild(selectNode)
	p.rootNode.AddChild(colListNode)
	p.rootNode.AddChild(fromNode)
	p.rootNode.AddChild(tableNameNode)
	p.rootNode.AddChild(optionalWhereNode)
	p.rootNode.AddChild(semicolonNode)
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

func (p *Parser) parseFromNode() ast.ASTNode {
	p.incrementPosition()
	fromNode := ast.ASTNode{ Data: p.tokens[p.pos], Children: []ast.ASTNode{} }
	return fromNode
}

func (p *Parser) parseTableName() ast.ASTNode {
	tableNameNoneTerminal := ast.ASTNode{ Data: "<table_name>", Children: []ast.ASTNode{} }
	p.incrementPosition()
	tableNameNode := ast.ASTNode{ Data: p.tokens[p.pos], Children: []ast.ASTNode{} }
	tableNameNoneTerminal.AddChild(tableNameNode)
	return tableNameNoneTerminal
}

func (p *Parser) parseOptionalWhere() ast.ASTNode {
	optionalWhereNode := ast.ASTNode{ Data: "<optional_where>", Children: []ast.ASTNode{} }
	
	nextToken := p.peek()
	if nextToken != "WHERE" {
		return optionalWhereNode
	}
	
	p.incrementPosition()
	whereNode := ast.ASTNode{ Data: p.tokens[p.pos], Children: []ast.ASTNode{} }
	optionalWhereNode.AddChild(whereNode)
	
	conditionNode := p.parseCondition()
	optionalWhereNode.AddChild(conditionNode)
	
	return optionalWhereNode
}

func (p *Parser) parseCondition() ast.ASTNode {
	conditionNode := ast.ASTNode{ Data: "<condition>", Children: []ast.ASTNode{} }
	
	p.incrementPosition()
	leftOperand := ast.ASTNode{ Data: p.tokens[p.pos], Children: []ast.ASTNode{} }
	conditionNode.AddChild(leftOperand)
	
	p.incrementPosition()
	operator := ast.ASTNode{ Data: p.tokens[p.pos], Children: []ast.ASTNode{} }
	conditionNode.AddChild(operator)
	
	p.incrementPosition()
	rightOperand := ast.ASTNode{ Data: p.tokens[p.pos], Children: []ast.ASTNode{} }
	conditionNode.AddChild(rightOperand)
	
	return conditionNode
}

func (p *Parser) parseSemicolon() ast.ASTNode {
	p.incrementPosition()
	semicolonNode := ast.ASTNode{ Data: p.tokens[p.pos], Children: []ast.ASTNode{} }
	return semicolonNode
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
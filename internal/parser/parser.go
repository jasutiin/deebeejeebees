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
			err := parser.parseInsert() // parsing insert statements

			if err != nil {
				fmt.Println(err.Error())
			}
		case "CREATE":
			err := parser.parseCreate() // parsing create statements

			if err != nil {
				fmt.Println(err.Error())
			}
		default:
			fmt.Println("could not determine the query type! must be one of: SELECT, INSERT, CREATE")
	}

	return rootNode
}

// parseSelect is responsible for parsing the SELECT query type
// SELECT (col1, col2, ...) FROM table_name WHERE column_name (operator) (value);
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

// parseInsert is responsible for parsing the INSERT query type
// INSERT INTO table_name (col1, col2, ...) VALUES (val1, val2, ...);
func (p *Parser) parseInsert() error {
	insertNode := ast.ASTNode{ Data: "INSERT", Children: []ast.ASTNode{} }
	
	intoNode, err := p.parseIntoNode()
	if err != nil {
		return err
	}
	
	tableNameNode := p.parseTableName()
	
	openParenNode1, err := p.parseOpenParen()
	if err != nil {
		return err
	}
	
	insertColListNode, err := p.parseInsertColumnList()
	if err != nil {
		return err
	}
	
	closeParenNode1, err := p.parseCloseParen()
	if err != nil {
		return err
	}
	
	valuesNode, err := p.parseValuesNode()
	if err != nil {
		return err
	}
	
	openParenNode2, err := p.parseOpenParen()
	if err != nil {
		return err
	}
	
	valueListNode, err := p.parseValueList()
	if err != nil {
		return err
	}
	
	closeParenNode2, err := p.parseCloseParen()
	if err != nil {
		return err
	}
	
	semicolonNode := p.parseSemicolon()
	
	p.rootNode.AddChild(insertNode)
	p.rootNode.AddChild(intoNode)
	p.rootNode.AddChild(tableNameNode)
	p.rootNode.AddChild(openParenNode1)
	p.rootNode.AddChild(insertColListNode)
	p.rootNode.AddChild(closeParenNode1)
	p.rootNode.AddChild(valuesNode)
	p.rootNode.AddChild(openParenNode2)
	p.rootNode.AddChild(valueListNode)
	p.rootNode.AddChild(closeParenNode2)
	p.rootNode.AddChild(semicolonNode)
	return nil
}

func (p *Parser) parseIntoNode() (ast.ASTNode, error) {
	nextToken := p.peek()
	if nextToken != "INTO" {
		return ast.ASTNode{}, fmt.Errorf("expected 'INTO' but got '%s'", nextToken)
	}
	p.incrementPosition()
	intoNode := ast.ASTNode{ Data: p.tokens[p.pos], Children: []ast.ASTNode{} }
	return intoNode, nil
}

func (p *Parser) parseOpenParen() (ast.ASTNode, error) {
	nextToken := p.peek()
	if nextToken != "(" {
		return ast.ASTNode{}, fmt.Errorf("expected '(' but got '%s'", nextToken)
	}
	p.incrementPosition()
	openParenNode := ast.ASTNode{ Data: p.tokens[p.pos], Children: []ast.ASTNode{} }
	return openParenNode, nil
}

func (p *Parser) parseCloseParen() (ast.ASTNode, error) {
	nextToken := p.peek()
	if nextToken != ")" {
		return ast.ASTNode{}, fmt.Errorf("expected ')' but got '%s'", nextToken)
	}
	p.incrementPosition()
	closeParenNode := ast.ASTNode{ Data: p.tokens[p.pos], Children: []ast.ASTNode{} }
	return closeParenNode, nil
}

func (p *Parser) parseInsertColumnList() (ast.ASTNode, error) {
	nextToken := p.peek()
	
	if nextToken == ")" {
		return ast.ASTNode{}, errors.New("missing at least one column name in INSERT!")
	}
	
	columnListNode := ast.ASTNode{ Data: "<column_list>", Children: []ast.ASTNode{} }
	p.incrementPosition()
	
	columnName := p.parseColumnName()
	columnListNode.AddChild(columnName)
	p.parseInsertColumnListTail(&columnListNode)
	
	return columnListNode, nil
}

func (p *Parser) parseInsertColumnListTail(parentNode *ast.ASTNode) {
	columnListTailNode := ast.ASTNode{ Data: "<column_list_tail>", Children: []ast.ASTNode{} }
	nextToken := p.peek()
	
	if nextToken == "," {
		p.incrementPosition()
		commaNode := ast.ASTNode{ Data: p.tokens[p.pos], Children: []ast.ASTNode{} }
		columnListTailNode.AddChild(commaNode)
		p.incrementPosition()
		columnName := p.parseColumnName()
		columnListTailNode.AddChild(columnName)
		p.parseInsertColumnListTail(&columnListTailNode)
	}
	
	parentNode.AddChild(columnListTailNode)
}

func (p *Parser) parseValuesNode() (ast.ASTNode, error) {
	nextToken := p.peek()
	if nextToken != "VALUES" {
		return ast.ASTNode{}, fmt.Errorf("expected 'VALUES' but got '%s'", nextToken)
	}
	p.incrementPosition()
	valuesNode := ast.ASTNode{ Data: p.tokens[p.pos], Children: []ast.ASTNode{} }
	return valuesNode, nil
}

func (p *Parser) parseValueList() (ast.ASTNode, error) {
	nextToken := p.peek()
	
	if nextToken == ")" {
		return ast.ASTNode{}, errors.New("missing at least one value in VALUES!")
	}
	
	valueListNode := ast.ASTNode{ Data: "<value_list>", Children: []ast.ASTNode{} }
	p.incrementPosition()
	
	value := p.parseValue()
	valueListNode.AddChild(value)
	p.parseValueListTail(&valueListNode)
	
	return valueListNode, nil
}

func (p *Parser) parseValue() ast.ASTNode {
	valueNonTerminal := ast.ASTNode{ Data: "<value>", Children: []ast.ASTNode{} }
	value := ast.ASTNode{ Data: p.tokens[p.pos], Children: []ast.ASTNode{} }
	valueNonTerminal.AddChild(value)
	return valueNonTerminal
}

func (p *Parser) parseValueListTail(parentNode *ast.ASTNode) {
	valueListTailNode := ast.ASTNode{ Data: "<value_list_tail>", Children: []ast.ASTNode{} }
	nextToken := p.peek()
	
	if nextToken == "," {
		p.incrementPosition()
		commaNode := ast.ASTNode{ Data: p.tokens[p.pos], Children: []ast.ASTNode{} }
		valueListTailNode.AddChild(commaNode)
		p.incrementPosition()
		value := p.parseValue()
		valueListTailNode.AddChild(value)
		p.parseValueListTail(&valueListTailNode)
	}
	
	parentNode.AddChild(valueListTailNode)
}

// parseCreate is responsible for parsing the CREATE query type
// CREATE TABLE table_name (col1 datatype1, col2 datatype2, ...);
func (p *Parser) parseCreate() error {
	createNode := ast.ASTNode{ Data: "CREATE", Children: []ast.ASTNode{} }
	
	tableKeywordNode, err := p.parseTableKeyword()
	if err != nil {
		return err
	}

	tableNameNode := p.parseTableName()

	openParenNode, err := p.parseOpenParen()
	if err != nil {
		return err
	}

	columnDefsNode, err := p.parseColumnDefsList()
	if err != nil {
		return err
	}

	closeParenNode, err := p.parseCloseParen()
	if err != nil {
		return err
	}

	semicolonNode := p.parseSemicolon()
	
	p.rootNode.AddChild(createNode)
	p.rootNode.AddChild(tableKeywordNode)
	p.rootNode.AddChild(tableNameNode)
	p.rootNode.AddChild(openParenNode)
	p.rootNode.AddChild(columnDefsNode)
	p.rootNode.AddChild(closeParenNode)
	p.rootNode.AddChild(semicolonNode)
	return nil
}

func (p *Parser) parseTableKeyword() (ast.ASTNode, error) {
	nextToken := p.peek()
	if nextToken != "TABLE" {
		return ast.ASTNode{}, fmt.Errorf("expected 'TABLE' but got '%s'", nextToken)
	}
	p.incrementPosition()
	tableKeywordNode := ast.ASTNode{ Data: p.tokens[p.pos], Children: []ast.ASTNode{} }
	return tableKeywordNode, nil
}

func (p *Parser) parseColumnDefsList() (ast.ASTNode, error) {
	nextToken := p.peek()
	
	if nextToken == ")" {
		return ast.ASTNode{}, errors.New("missing at least one column definition in CREATE TABLE!")
	}
	
	columnDefsNode := ast.ASTNode{ Data: "<column_defs_list>", Children: []ast.ASTNode{} }
	
	columnDef := p.parseColumnDef()
	columnDefsNode.AddChild(columnDef)
	p.parseColumnDefsListTail(&columnDefsNode)
	
	return columnDefsNode, nil
}

func (p *Parser) parseColumnDef() ast.ASTNode {
	columnDefNode := ast.ASTNode{ Data: "<column_def>", Children: []ast.ASTNode{} }
	
	p.incrementPosition()
	columnNameNode := ast.ASTNode{ Data: p.tokens[p.pos], Children: []ast.ASTNode{} }
	columnDefNode.AddChild(columnNameNode)
	
	p.incrementPosition()
	dataTypeNode := p.parseDataType()
	columnDefNode.AddChild(dataTypeNode)
	
	return columnDefNode
}

func (p *Parser) parseDataType() ast.ASTNode {
	dataTypeNonTerminal := ast.ASTNode{ Data: "<data_type>", Children: []ast.ASTNode{} }
	dataType := ast.ASTNode{ Data: p.tokens[p.pos], Children: []ast.ASTNode{} }
	dataTypeNonTerminal.AddChild(dataType)
	
	nextToken := p.peek()
	if nextToken == "(" {
		p.incrementPosition()
		openParenNode := ast.ASTNode{ Data: p.tokens[p.pos], Children: []ast.ASTNode{} }
		dataTypeNonTerminal.AddChild(openParenNode)
		
		p.incrementPosition()
		sizeNode := ast.ASTNode{ Data: p.tokens[p.pos], Children: []ast.ASTNode{} }
		dataTypeNonTerminal.AddChild(sizeNode)
		
		p.incrementPosition()
		closeParenNode := ast.ASTNode{ Data: p.tokens[p.pos], Children: []ast.ASTNode{} }
		dataTypeNonTerminal.AddChild(closeParenNode)
	}
	
	return dataTypeNonTerminal
}

func (p *Parser) parseColumnDefsListTail(parentNode *ast.ASTNode) {
	columnDefsListTailNode := ast.ASTNode{ Data: "<column_defs_list_tail>", Children: []ast.ASTNode{} }
	nextToken := p.peek()
	
	if nextToken == "," {
		p.incrementPosition()
		commaNode := ast.ASTNode{ Data: p.tokens[p.pos], Children: []ast.ASTNode{} }
		columnDefsListTailNode.AddChild(commaNode)
		
		columnDef := p.parseColumnDef()
		columnDefsListTailNode.AddChild(columnDef)
		p.parseColumnDefsListTail(&columnDefsListTailNode)
	}
	
	parentNode.AddChild(columnDefsListTailNode)
}

func (p *Parser) incrementPosition() {
	p.pos += 1
}

func (p *Parser) peek() string {
	return p.tokens[p.pos + 1]
}
package parser

import (
	"errors"
	"fmt"

	"github.com/jasutiin/deebeejeebees/internal/parser/narytree"
)

type Parser struct {
	tokens []string
	pos int
	rootNode *narytree.Node
}

func ParseTokens(tokens []string) narytree.Node {
	rootNode := narytree.Node{ Data: "<query>", Children: []narytree.Node{} }
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
	selectNode := narytree.Node{ Data: "SELECT", Children: []narytree.Node{} }
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

func (p *Parser) parseColumnList() (narytree.Node, error) {
	nextToken := p.peek()

	if nextToken == "FROM" {
		return narytree.Node{}, errors.New("missing at lenarytree one column name after SELECT!")
	}

	columnListNode := narytree.Node { Data: "<column_list>", Children: []narytree.Node{} }
	p.incrementPosition()

	columnName := p.parseColumnName()
	columnListNode.AddChild(columnName)
	p.parseColumnListTail(&columnListNode)

	return columnListNode, nil
}

func (p *Parser) parseColumnName() narytree.Node {
	columnNameNonTerminal := narytree.Node{ Data: "<column_name>", Children: []narytree.Node{} }
	columnName := narytree.Node{ Data: p.tokens[p.pos], Children: []narytree.Node{} }
	columnNameNonTerminal.Children = append(columnNameNonTerminal.Children, columnName)
	return columnNameNonTerminal
}

func (p *Parser) parseColumnListTail(parentNode *narytree.Node) narytree.Node {
	columnListTailNode := narytree.Node{ Data: "<column_list_tail>", Children: []narytree.Node{} }
	nextToken := p.peek()
	
	if nextToken == "," {
		p.incrementPosition()
		commaNode := narytree.Node{ Data: p.tokens[p.pos], Children: []narytree.Node{} }
		columnListTailNode.AddChild(commaNode)
		p.incrementPosition()
		columnName := p.parseColumnName()
		columnListTailNode.AddChild(columnName)
		p.parseColumnListTail(&columnListTailNode)
	}
	
	parentNode.AddChild(columnListTailNode)
	return columnListTailNode
}

func (p *Parser) parseFromNode() narytree.Node {
	p.incrementPosition()
	fromNode := narytree.Node{ Data: p.tokens[p.pos], Children: []narytree.Node{} }
	return fromNode
}

func (p *Parser) parseTableName() narytree.Node {
	tableNameNoneTerminal := narytree.Node{ Data: "<table_name>", Children: []narytree.Node{} }
	p.incrementPosition()
	tableNameNode := narytree.Node{ Data: p.tokens[p.pos], Children: []narytree.Node{} }
	tableNameNoneTerminal.AddChild(tableNameNode)
	return tableNameNoneTerminal
}

func (p *Parser) parseOptionalWhere() narytree.Node {
	optionalWhereNode := narytree.Node{ Data: "<optional_where>", Children: []narytree.Node{} }
	
	nextToken := p.peek()
	if nextToken != "WHERE" {
		return optionalWhereNode
	}
	
	p.incrementPosition()
	whereNode := narytree.Node{ Data: p.tokens[p.pos], Children: []narytree.Node{} }
	optionalWhereNode.AddChild(whereNode)
	
	conditionNode := p.parseCondition()
	optionalWhereNode.AddChild(conditionNode)
	
	return optionalWhereNode
}

func (p *Parser) parseCondition() narytree.Node {
	conditionNode := narytree.Node{ Data: "<condition>", Children: []narytree.Node{} }
	
	p.incrementPosition()
	leftOperand := narytree.Node{ Data: p.tokens[p.pos], Children: []narytree.Node{} }
	conditionNode.AddChild(leftOperand)
	
	p.incrementPosition()
	operator := narytree.Node{ Data: p.tokens[p.pos], Children: []narytree.Node{} }
	conditionNode.AddChild(operator)
	
	p.incrementPosition()
	rightOperand := narytree.Node{ Data: p.tokens[p.pos], Children: []narytree.Node{} }
	conditionNode.AddChild(rightOperand)
	
	return conditionNode
}

func (p *Parser) parseSemicolon() narytree.Node {
	p.incrementPosition()
	semicolonNode := narytree.Node{ Data: p.tokens[p.pos], Children: []narytree.Node{} }
	return semicolonNode
}

// parseInsert is responsible for parsing the INSERT query type
// INSERT INTO table_name (col1, col2, ...) VALUES (val1, val2, ...);
func (p *Parser) parseInsert() error {
	insertNode := narytree.Node{ Data: "INSERT", Children: []narytree.Node{} }
	
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

func (p *Parser) parseIntoNode() (narytree.Node, error) {
	nextToken := p.peek()
	if nextToken != "INTO" {
		return narytree.Node{}, fmt.Errorf("expected 'INTO' but got '%s'", nextToken)
	}
	p.incrementPosition()
	intoNode := narytree.Node{ Data: p.tokens[p.pos], Children: []narytree.Node{} }
	return intoNode, nil
}

func (p *Parser) parseOpenParen() (narytree.Node, error) {
	nextToken := p.peek()
	if nextToken != "(" {
		return narytree.Node{}, fmt.Errorf("expected '(' but got '%s'", nextToken)
	}
	p.incrementPosition()
	openParenNode := narytree.Node{ Data: p.tokens[p.pos], Children: []narytree.Node{} }
	return openParenNode, nil
}

func (p *Parser) parseCloseParen() (narytree.Node, error) {
	nextToken := p.peek()
	if nextToken != ")" {
		return narytree.Node{}, fmt.Errorf("expected ')' but got '%s'", nextToken)
	}
	p.incrementPosition()
	closeParenNode := narytree.Node{ Data: p.tokens[p.pos], Children: []narytree.Node{} }
	return closeParenNode, nil
}

func (p *Parser) parseInsertColumnList() (narytree.Node, error) {
	nextToken := p.peek()
	
	if nextToken == ")" {
		return narytree.Node{}, errors.New("missing at lenarytree one column name in INSERT!")
	}
	
	columnListNode := narytree.Node{ Data: "<column_list>", Children: []narytree.Node{} }
	p.incrementPosition()
	
	columnName := p.parseColumnName()
	columnListNode.AddChild(columnName)
	p.parseInsertColumnListTail(&columnListNode)
	
	return columnListNode, nil
}

func (p *Parser) parseInsertColumnListTail(parentNode *narytree.Node) {
	columnListTailNode := narytree.Node{ Data: "<column_list_tail>", Children: []narytree.Node{} }
	nextToken := p.peek()
	
	if nextToken == "," {
		p.incrementPosition()
		commaNode := narytree.Node{ Data: p.tokens[p.pos], Children: []narytree.Node{} }
		columnListTailNode.AddChild(commaNode)
		p.incrementPosition()
		columnName := p.parseColumnName()
		columnListTailNode.AddChild(columnName)
		p.parseInsertColumnListTail(&columnListTailNode)
	}
	
	parentNode.AddChild(columnListTailNode)
}

func (p *Parser) parseValuesNode() (narytree.Node, error) {
	nextToken := p.peek()
	if nextToken != "VALUES" {
		return narytree.Node{}, fmt.Errorf("expected 'VALUES' but got '%s'", nextToken)
	}
	p.incrementPosition()
	valuesNode := narytree.Node{ Data: p.tokens[p.pos], Children: []narytree.Node{} }
	return valuesNode, nil
}

func (p *Parser) parseValueList() (narytree.Node, error) {
	nextToken := p.peek()
	
	if nextToken == ")" {
		return narytree.Node{}, errors.New("missing at lenarytree one value in VALUES!")
	}
	
	valueListNode := narytree.Node{ Data: "<value_list>", Children: []narytree.Node{} }
	p.incrementPosition()
	
	value := p.parseValue()
	valueListNode.AddChild(value)
	p.parseValueListTail(&valueListNode)
	
	return valueListNode, nil
}

func (p *Parser) parseValue() narytree.Node {
	valueNonTerminal := narytree.Node{ Data: "<value>", Children: []narytree.Node{} }
	value := narytree.Node{ Data: p.tokens[p.pos], Children: []narytree.Node{} }
	valueNonTerminal.AddChild(value)
	return valueNonTerminal
}

func (p *Parser) parseValueListTail(parentNode *narytree.Node) {
	valueListTailNode := narytree.Node{ Data: "<value_list_tail>", Children: []narytree.Node{} }
	nextToken := p.peek()
	
	if nextToken == "," {
		p.incrementPosition()
		commaNode := narytree.Node{ Data: p.tokens[p.pos], Children: []narytree.Node{} }
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
	createNode := narytree.Node{ Data: "CREATE", Children: []narytree.Node{} }
	
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

func (p *Parser) parseTableKeyword() (narytree.Node, error) {
	nextToken := p.peek()
	if nextToken != "TABLE" {
		return narytree.Node{}, fmt.Errorf("expected 'TABLE' but got '%s'", nextToken)
	}
	p.incrementPosition()
	tableKeywordNode := narytree.Node{ Data: p.tokens[p.pos], Children: []narytree.Node{} }
	return tableKeywordNode, nil
}

func (p *Parser) parseColumnDefsList() (narytree.Node, error) {
	nextToken := p.peek()
	
	if nextToken == ")" {
		return narytree.Node{}, errors.New("missing at lenarytree one column definition in CREATE TABLE!")
	}
	
	columnDefsNode := narytree.Node{ Data: "<column_defs_list>", Children: []narytree.Node{} }
	
	columnDef := p.parseColumnDef()
	columnDefsNode.AddChild(columnDef)
	p.parseColumnDefsListTail(&columnDefsNode)
	
	return columnDefsNode, nil
}

func (p *Parser) parseColumnDef() narytree.Node {
	columnDefNode := narytree.Node{ Data: "<column_def>", Children: []narytree.Node{} }
	
	p.incrementPosition()
	columnNameNode := narytree.Node{ Data: p.tokens[p.pos], Children: []narytree.Node{} }
	columnDefNode.AddChild(columnNameNode)
	
	p.incrementPosition()
	dataTypeNode := p.parseDataType()
	columnDefNode.AddChild(dataTypeNode)
	
	return columnDefNode
}

func (p *Parser) parseDataType() narytree.Node {
	dataTypeNonTerminal := narytree.Node{ Data: "<data_type>", Children: []narytree.Node{} }
	dataType := narytree.Node{ Data: p.tokens[p.pos], Children: []narytree.Node{} }
	dataTypeNonTerminal.AddChild(dataType)
	
	nextToken := p.peek()
	if nextToken == "(" {
		p.incrementPosition()
		openParenNode := narytree.Node{ Data: p.tokens[p.pos], Children: []narytree.Node{} }
		dataTypeNonTerminal.AddChild(openParenNode)
		
		p.incrementPosition()
		sizeNode := narytree.Node{ Data: p.tokens[p.pos], Children: []narytree.Node{} }
		dataTypeNonTerminal.AddChild(sizeNode)
		
		p.incrementPosition()
		closeParenNode := narytree.Node{ Data: p.tokens[p.pos], Children: []narytree.Node{} }
		dataTypeNonTerminal.AddChild(closeParenNode)
	}
	
	return dataTypeNonTerminal
}

func (p *Parser) parseColumnDefsListTail(parentNode *narytree.Node) {
	columnDefsListTailNode := narytree.Node{ Data: "<column_defs_list_tail>", Children: []narytree.Node{} }
	nextToken := p.peek()
	
	if nextToken == "," {
		p.incrementPosition()
		commaNode := narytree.Node{ Data: p.tokens[p.pos], Children: []narytree.Node{} }
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
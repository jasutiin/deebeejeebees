package main

import (
	"fmt"
	"os"

	"github.com/jasutiin/deebeejeebees/internal/lexer"
	"github.com/jasutiin/deebeejeebees/internal/parser"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("you need to provide the query")
	} else {
		tokens := lexer.AnalyzeString(args[0])

		fmt.Println("=== TOKENS ===")
		for _, val := range tokens {
			fmt.Println(val)
		}

		fmt.Println("=== PARSE TREE ===")
		cstTree := parser.ParseTokensToCST(tokens)
		cstTree.PrintTree()

		fmt.Println("=== AST ===")
		astTree := parser.ConvertToAST(cstTree)
		astTree.PrintTree()
	}
}
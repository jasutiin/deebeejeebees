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
		for _, val := range tokens {
			fmt.Println(val)
		}
		astTree := parser.ParseTokens(tokens)
	}
}
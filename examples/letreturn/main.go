package main

import (
	"fmt"

	"github.com/animeshmukherjeework911-jpg/goInterpreter/ast"
	"github.com/animeshmukherjeework911-jpg/goInterpreter/lexer"
	"github.com/animeshmukherjeework911-jpg/goInterpreter/parser"
)

func parseAndPrint(input string) {
	fmt.Printf("input: %s\n", input)

	p := parser.New(lexer.New(input))
	program := p.ParseProgram()

	for i, stmt := range program.Statements {
		switch s := stmt.(type) {
		case *ast.LetStatement:
			fmt.Printf("  %d: let, Name=%s, Value=%v\n", i, s.Name.Value, s.Value)
		case *ast.ReturnStatement:
			fmt.Printf("  %d: return, ReturnValue=%v\n", i, s.ReturnValue)
		}
	}

	for _, msg := range p.Errors() {
		fmt.Printf("  error: %s\n", msg)
	}
	fmt.Println()
}

func main() {
	parseAndPrint("let x = 5; let total = x + 10; return total;")
	parseAndPrint("let = 5; let y = 1; return y;")
}

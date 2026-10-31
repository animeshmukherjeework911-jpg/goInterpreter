package parser

import (
	"testing"

	"github.com/animeshmukherjeework911-jpg/goInterpreter/ast"
	"github.com/animeshmukherjeework911-jpg/goInterpreter/lexer"
)

func TestLetStatement(t *testing.T) {
	input := `

	let x = 5;
	let y = 10;
	let foobar = 838383;

	`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements, got=%d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"}, {"y"}, {"foobar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func TestLetStatementWithoutSemicolon(t *testing.T) {
	l := lexer.New("let x = 5")
	p := New(l)

	program := p.ParseProgram()

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement, got=%d", len(program.Statements))
	}

	testLetStatement(t, program.Statements[0], "x")
}

func TestInvalidLetStatements(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedError string
	}{
		{"missing identifier", "let = 5;", "expected next token to be IDENT, got = instead"},
		{"missing assign", "let x 5;", "expected next token to be =, got INT instead"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New(lexer.New(tt.input))
			program := p.ParseProgram()

			if len(program.Statements) != 0 {
				t.Fatalf("program.Statements does not contain 0 statements, got=%d", len(program.Statements))
			}

			errors := p.Errors()
			if len(errors) != 1 {
				t.Fatalf("expected 1 parser error, got=%d: %q", len(errors), errors)
			}

			if errors[0] != tt.expectedError {
				t.Errorf("wrong error. expected=%q, got=%q", tt.expectedError, errors[0])
			}
		})
	}
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got=%q", s.TokenLiteral())
		return false
	}

	letStmt, ok := s.(*ast.LetStatement)

	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%T", s)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("letStmt.Name.TokenLiteral() not '%s'. got=%s", name, letStmt.Name.TokenLiteral())
		return false
	}

	return true
}

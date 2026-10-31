# 02: AST and parser (Oct 31)

Designed the AST, wrote a parser for `let` and `return` that skips their values for now, fixed two bugs that bad input exposed, and made the parser collect errors. This covers Sessions 1 and 2 of the plan.

## The AST

```go
type Node interface {
	TokenLiteral() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}
```

`statementNode()` and `expressionNode()` are empty and never called. Go interfaces match on method sets, so without them the two interfaces would be identical, and the compiler couldn't stop an expression being used where a statement belongs.

```go
type Program struct {
	Statements []Statement
}

type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
```

`Identifier` is an expression, even as the name in a `let`, which keeps the node types simple.

## The parser holds two tokens

```go
type Parser struct {
	l         *lexer.Lexer
	curToken  token.Token
	peekToken token.Token
	errors    []string
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}
```

`New` calls `nextToken` twice to fill both slots. `peekToken` is the same idea as `peekChar` in the lexer, one level up.

## `expectPeek` drives parsing

```go
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}
```

It advances only if the check passes, so `parseLetStatement` reads like the grammar:

```go
func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// TODO: We are skipping the expression until we encounter a semicolon

	for !p.curTokenIs(token.SEMICOLON) && !p.curTokenIs(token.EOF) {
		p.nextToken()
	}

	return stmt
}
```

Tracing `let x = 5;`:

```
                        cur    peek
after New               let    x
expectPeek(IDENT)       x      =      ← Name is set here
expectPeek(ASSIGN)      =      5
skip value              5      ;
skip stops              ;      EOF
ParseProgram advances   EOF    EOF    ← loop ends
```

`Value` stays nil. That's deliberate: expressions get Pratt parsing in Session 3.

## `return` is simpler

```go
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}
	p.nextToken()

	// TODO: We are skipping the expression until we encounter a semicolon
	for !p.curTokenIs(token.SEMICOLON) && !p.curTokenIs(token.EOF) {
		p.nextToken()
	}

	return stmt
}
```

There's no name or `=` to check, so it doesn't need `expectPeek` and can't fail yet.

## Bug 1: the hang

`let x = 5` with no semicolon hung forever. The skip loop only stopped at `;`, and the lexer keeps returning `EOF`.

```go
// before
for !p.curTokenIs(token.SEMICOLON) {
	p.nextToken()
}

// after
for !p.curTokenIs(token.SEMICOLON) && !p.curTokenIs(token.EOF) {
	p.nextToken()
}
```

**Rule:** any loop that consumes tokens must also stop at `EOF`.

## Bug 2: the typed nil

`let = 5;` panicked even though `ParseProgram` checks `stmt != nil`. An interface only equals nil when both its type and value are empty:

```go
var p *ast.LetStatement = nil // a nil pointer
var s ast.Statement = p       // s holds (type: *LetStatement, value: nil)
s == nil                      // false!
```

So the check passed, the broken statement joined the program, and the first method call on it panicked. The fix checks while it's still a plain pointer:

```go
func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		// A nil *ast.LetStatement returned as ast.Statement would not compare equal to nil
		if stmt := p.parseLetStatement(); stmt != nil {
			return stmt
		}
		return nil

	case token.RETURN:
		if stmt := p.parseReturnStatement(); stmt != nil {
			return stmt
		}
		return nil

	default:
		return nil
	}
}
```

**Rule:** be suspicious when a function returns a concrete pointer through an interface return type.

## Collecting errors

```go
func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}
```

This is called from `expectPeek`, the one place that knows both what was wanted and what came, so every caller gets error reporting for free.

```
let x 5;   →  expected next token to be =, got INT instead
let = 5;   →  expected next token to be IDENT, got = instead
```

`let = 5;` gives only one error. After the failure, `ParseProgram` advances through `=`, `5` and `;`, and each one hits the silent `default` case above. A stray `5;` gets no error for the same reason. `let x = ;` and `return ;` are also accepted, because the value is skipped.

## Write these yourself

Claude wrote the two bug fixes and their tests, so type those in by hand: the `EOF` check in the skip loop and the guard in `parseStatement`.

## Commits

| Commit | What it added |
|---|---|
| `ebfdfe7` | VS Code debug configs |
| `7be4b04` | AST and `let` parsing, including both bug fixes |
| `789be92` | Error collection through `expectPeek` |
| `0bcef45` | `return` statements |

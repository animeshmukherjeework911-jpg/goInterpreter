# 01: Lexer and REPL (Oct 26)

Built a lexer that turns Monkey source into tokens, put a REPL on top so I could see them, then cleaned up before starting the parser.

## A token is a type plus its text

```go
type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}
```

`5` and `10` are both `INT`, so the literal has to come along to tell them apart.

## The lexer is a cursor

```go
type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition += 1
}
```

`ch == 0` means the input has ended, and `NextToken` turns that into `EOF`. Calling it again gives `EOF` again, forever. That caused a real bug in session 02.

## Looking ahead tells `=` from `==`

```go
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}
```

```go
case '=':
	if l.peekChar() == '=' {
		ch := l.ch
		l.readChar()
		literal := string(ch) + string(l.ch)
		tok = token.Token{Type: token.EQ, Literal: literal}
	} else {
		tok = newToken(token.ASSIGN, l.ch)
	}
```

`peekChar` looks without moving. `!` and `!=` use the same pattern.

## Words first, keywords second

```go
var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
```

`let` and `foo` are read the same way, and the map decides which is which. A new keyword only needs a map entry.

## The easy-to-miss early return

```go
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok // already one past the word
		} else if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok // already one past the number
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}
```

`readIdentifier` and `readNumber` stop one past what they read. Reaching the final `readChar()` would swallow a character, so these branches return first.

A bad character doesn't stop the lexer. It becomes an `ILLEGAL` token and lexing carries on.

## What it does with odd input

| Input | Tokens | Why |
|---|---|---|
| `x1` | `x`, `1` | Identifiers can't contain digits |
| `<=` | `<`, `=` | Only `==` and `!=` are two-character operators |
| `5x` | `5`, `x` | A number stops at the first non-digit |
| `@` | `ILLEGAL` | No case matches it |

## The REPL

```go
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Fprint(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)

		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Fprintf(out, "%+v\n", tok)
		}
	}
}
```

It takes `io.Reader` and `io.Writer` rather than stdin and stdout, so it isn't tied to a terminal. The prompt uses `Fprint` because `PROMPT` has no format verbs.

```
>> let x 5;
{Type:LET Literal:let}
{Type:IDENT Literal:x}
{Type:INT Literal:5}
{Type:; Literal:;}
```

The lexer doesn't object to the missing `=`. Whether the tokens make sense is the parser's job.

## Cleanup before the parser

Commit `9a62a82` changed the module path from `monkey` to `github.com/animeshmukherjeework911-jpg/goInterpreter`, fixed typos, switched to `Fprint`, and added lexer tests for `@` and empty input.

## Commits

| Commit | What it added |
|---|---|
| `6e461dc` | Go module and a lexer for single-character tokens |
| `2d927cf` | Identifiers, integers, `let` and `fn` |
| `2f4c810` | Operators, `==` and `!=`, the remaining keywords |
| `3a237b4` | REPL |
| `07ca34a` | `.gitignore` for build artifacts |
| `9a62a82` | Module path, typos, `Fprint`, edge-case tests |
| `3799799` | Study plan for the next chapters |

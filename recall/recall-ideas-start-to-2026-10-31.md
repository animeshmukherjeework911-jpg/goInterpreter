# The ideas behind my interpreter

*Read this to get the mental model back. Before each section, try its "Recall first" prompt from memory, even if you get it wrong. The struggle is what makes it stick.*

---

## 1. The big picture

**Recall first:** Without looking, sketch the stages from source text to result. What does each stage take in and produce?

To a computer, `let x = 5;` is just ten characters. It has no idea that a variable exists. An interpreter builds up meaning in stages, and each stage answers one question:

```
"let x = 5;"
     │  Lexer: what are the words?
     ▼
LET  IDENT(x)  =  INT(5)  ;
     │  Parser: how do the words relate?
     ▼
LetStatement { Name: x, Value: 5 }
     │  Evaluator: what does it mean? (not built yet)
     ▼
x now holds 5
```

**The idea to keep:** each stage only solves its own problem. The lexer has no idea what a valid program is, so it will happily tokenize `let let let`. Spotting that nonsense is the parser's job. This split is why each stage can be small.

---

## 2. Tokens: giving characters a meaning

**Recall first:** Why have tokens at all? Why does a token need both a `Type` and a `Literal`?

Working on raw characters would be miserable: the parser would have to handle spaces, and work out letter by letter that `l-e-t` means `let`. Tokens remove all that, so later stages can think in words.

```go
type Token struct {
    Type    TokenType // what kind of word: INT
    Literal string    // the actual text: "5"
}
```

The parser mostly needs the type: "is this an integer?" But `5` and `10` have the same type and differ only in their text, so the literal has to come along too.

`TokenType` is a string, so printed tokens are readable when debugging. Integers would compare faster, but learning from readable output matters more here.

---

## 3. The lexer: walking through text

**Recall first:** What do `position`, `readPosition`, and `ch` each hold? How does the lexer tell `=` from `==`? Why do the word and number branches `return` early?

The lexer is a cursor moving over the text, with `ch` the character under it. Each `NextToken()` call skips whitespace, looks at `ch`, and decides what it's looking at. It produces **one token per call**, only when asked.

Three small ideas do most of the work:

- **Looking ahead:** at `=`, the lexer can't decide yet, because it could be `=` or `==`. `peekChar()` looks at the next character without moving. This is the same trick the parser uses later.
- **Words first, keywords second:** it reads any word, then asks `LookupIdent` whether that word is a keyword. So `let` and `foo` are read the same way, and adding a keyword only means adding a map entry.
- **The easy-to-miss early return:** `NextToken` normally finishes by moving past the current character. But `readIdentifier` and `readNumber` stop *one past* the end of what they read, so they're already in position. Moving again would swallow a character, which is why they `return` early.

**Two behaviours that matter later:**

- A bad character doesn't stop the lexer. It becomes an `ILLEGAL` token, and a later stage deals with it.
- At the end of the input, the lexer **returns `EOF` forever**. This caused a real bug (section 7).

---

## 4. The AST: why a tree?

**Recall first:** Why isn't a flat list of tokens enough? What's the difference between a statement and an expression? What are the empty `statementNode()` methods for?

Tokens are flat, but programs are nested. In `1 + 2 * 3`, the `2 * 3` happens first and its result feeds into the `+`. A tree shows that structure directly, without any brackets:

```
    +
   / \
  1   *
     / \
    2   3
```

A tree can be evaluated by visiting children first, so the structure itself gives the order. This is the payoff that the parser sets up.

**Two kinds of node:**

- An **expression** produces a value: `5`, `x`, `add(1, 2)`
- A **statement** does something and produces no value: `let x = 5;`

**The neat Go trick:** `statementNode()` and `expressionNode()` are empty and never called. They exist because Go interfaces match on method sets. Without them, the two interfaces would be identical, and the compiler couldn't stop me putting an expression where a statement belongs. They turn a design rule into a compile error.

So far there are four nodes: `Program` (the root, a list of statements), `LetStatement`, `ReturnStatement`, and `Identifier`.

---

## 5. The parser: turning a flat list into a tree

**Recall first:** Why does the parser hold two tokens? Trace `let x = 5;` by hand, writing `curToken` and `peekToken` at each step. What does `expectPeek` do besides check?

The parser is **recursive descent**: one function per grammar rule, and each function looks like the rule it parses. A `let` is "the word `let`, a name, `=`, a value", and `parseLetStatement` reads in that same order:

```
ParseProgram → parseStatement → parseLetStatement
```

**Why two tokens?** `curToken` is where I am, and `peekToken` tells me what to do next. Sometimes the current token isn't enough: after `5`, is the expression over (`5;`), or just starting (`5 + 5`)? It's the same looking-ahead idea as `peekChar()`, one level up.

**`expectPeek` is the key move:** it checks *and* advances, but only if the check passes. That's why `parseLetStatement` reads almost like the grammar:

```go
if !p.expectPeek(token.IDENT)  { return nil } // need a name
stmt.Name = ...
if !p.expectPeek(token.ASSIGN) { return nil } // need an =
```

The trace to check yourself against:

```
                        cur    peek
after New               let    x
expectPeek(IDENT)       x      =      ← Name is set here
expectPeek(ASSIGN)      =      5
skip value              5      ;
skip stops              ;      EOF
ParseProgram advances   EOF    EOF    ← loop ends
```

---

## 6. The deliberate shortcut

**Recall first:** What's in `Value` after parsing `let x = 5;`? Why did I leave it that way?

`Value` is **nil**, and so is `ReturnValue` on a `return`. The parser skips to the `;` without reading what's there, because expressions are the hard part and get their own technique (Pratt parsing, Session 3).

Skipping lets me get the whole skeleton working first and fill in the hardest piece later. The cost is that the parser currently accepts `let x = ;` and `return ;`, and `let a = 1 let b = 2;` loses `b`. Those aren't bugs to chase, because real expression parsing replaces these loops.

---

## 7. The two bugs: the most valuable part

**Recall first:** Why did `let x = 5` hang? Why did `let = 5;` crash even though there's a nil check? Explain the second one out loud, as if teaching someone.

### The hang

The skip loop only stopped at `;`. With no `;`, it kept asking for tokens, and the lexer kept answering `EOF` forever.

**The rule:** any loop that consumes tokens must also stop at `EOF`.

### The typed nil trap

This one is about Go itself. An interface value holds two things: a type and a value. It only equals `nil` when **both** are empty.

```go
var p *ast.LetStatement = nil  // a nil pointer
var s ast.Statement = p        // s holds (type: *LetStatement, value: nil)
s == nil                       // false!
```

`parseLetStatement` returned a nil pointer, and `parseStatement` handed it back as an `ast.Statement`. That put a type in the interface, so `ParseProgram`'s `stmt != nil` check passed, the broken statement joined the program, and the first method call on it panicked.

The fix checks while it's still a plain pointer:

```go
if stmt := p.parseLetStatement(); stmt != nil {
    return stmt
}
return nil // a true nil interface
```

**The rule:** be suspicious when a function returns a concrete pointer through an interface return type. `parseStatement` uses the same guard for `parseReturnStatement`. That function can't return nil today, so the guard is only a precaution for when expression parsing can fail.

---

## 8. Return statements and parser errors

**Recall first:** Trace `return 5;` by hand. Why doesn't it use `expectPeek`? Where do parser errors get recorded, and why there? How many errors does `let = 5;` produce, and why not more?

A `return` is simpler than a `let`: the word, then a value. There's no name or `=` to check, so nothing can fail yet.

```
                        cur     peek
after New               return  5
nextToken               5       ;
skip stops              ;       EOF
ParseProgram advances   EOF     EOF    ← loop ends
```

**Errors live in `expectPeek`,** because it's the one place that knows both what was wanted and what actually came. It calls `peekError`, which appends to a list, and `Errors()` exposes that list. Every `expectPeek` caller gets error reporting without writing any code for it.

**The parser collects errors and keeps going.** It doesn't stop at the first problem, so one run could report several.

`let = 5;` still produces only one error, for a less impressive reason:

1. `expectPeek(IDENT)` fails, records the error, and leaves `curToken` on `let`.
2. `ParseProgram` advances to `=`, then `5`, then `;`.
3. `parseStatement` silently returns nil for each, because they aren't `let` or `return`.

That silent `default` case is why a stray `5;` currently gives no error at all, and it's also the gap Session 3 fills.

**Reading the messages:** `got = instead` looks odd, but the type of `ASSIGN` really is the string `"="`. Compare `got INT instead` for `let x 5;`.

---

## 9. Limits I've accepted for now

**Recall first:** What does each of these lex as: `x1`, `<=`, `5x`?

The lexer keeps things simple, so:

- **`x1`** becomes `x`, `1`, because identifiers can't contain digits.
- **`<=`** becomes `<`, `=`, because only `==` and `!=` are two-character operators.
- **`5x`** becomes `5`, `x`, because a number just stops at the first non-digit.
- **`é`** becomes two `ILLEGAL` tokens, because the lexer reads bytes, not Unicode characters.

Tokens also don't record their line or column, so error messages won't be able to say where a problem is.

---

## Teach it back

Close this file and explain these out loud. Reread only what you can't explain.

1. What happens to `let x = 5;` from characters to tree, and what's still missing?
2. Why do the lexer and parser both need to look ahead?
3. Why does the AST have empty marker methods?
4. What is a typed nil, and where could it bite me next?
5. Why is `expectPeek` the right home for error reporting, and which bad input still produces no error?

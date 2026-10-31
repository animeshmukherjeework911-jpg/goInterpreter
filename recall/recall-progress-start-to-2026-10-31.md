# How far I've got

*Read this when coming back to the project, to get back into the thread. Last updated 2026-10-31.*

## Before reading, try this

Without opening any code, write down:

1. What can my interpreter do right now, and what can't it do?
2. What did I build most recently?
3. What am I building next, and what will be the first thing to change?

Then read on and notice what you forgot. Those gaps are what to reread in [recall-ideas-start-to-2026-10-31.md](recall-ideas-start-to-2026-10-31.md).

---

## The story so far

### Step 1: Teaching it to read words (Oct 26)

I started with the smallest possible lexer, one that only knew single symbols like `+` and `(`. Then I grew it in layers: first identifiers, numbers, and `let`/`fn`, then operators and the remaining keywords.

**What clicked here:** a lexer is only a cursor plus a big `switch`. Looking one character ahead is all it takes to tell `=` from `==`.

### Step 2: Making it something I could touch (Oct 26)

I added the REPL, so I could type a line and see its tokens. Then I did a cleanup pass before moving on: a proper module path, typo fixes, and edge-case tests for illegal characters and empty input.

**What clicked here:** `Start` takes an `io.Reader` and `io.Writer` rather than stdin and stdout directly, so it isn't tied to a terminal.

### Step 3: Teaching it structure (Oct 30–31)

This was the jump from words to meaning. I designed the AST (`Node`, `Statement`, `Expression`), then wrote a parser that understands `let` statements. I wrote it with a deliberate shortcut: it reads the name and skips the value.

**What clicked here:** the tree is the real goal of the parser, and the interfaces with empty marker methods let the compiler enforce the design.

### Step 4: Finding out where it breaks (Oct 31)

The tests all passed, but they only tried valid input. Two bad inputs exposed real bugs: `let x = 5` hung forever, and `let = 5;` crashed despite a nil check. The second one taught me Go's typed nil trap.

**What clicked here:** passing tests only prove what they try. Bad input is where the bugs were hiding.

**Own these fixes:** Claude wrote these two fixes and their tests, so I didn't type them. Since my rule is to type everything myself, I should be able to rewrite them from memory. Section 7 of [recall-ideas-start-to-2026-10-31.md](recall-ideas-start-to-2026-10-31.md) explains them.

### Step 5: Teaching it to complain (Oct 31)

I added `return` statements, which skip their value just as `let` does, and made the parser collect errors through `peekError` and `Errors()`. The tests now cover `let = 5;`, `let x 5;`, and `return 5` with no semicolon.

**What clicked here:** `expectPeek` already knew what was wanted and what came, so one small change there gave every caller error reporting.

---

## Where that leaves me

**Can do:** tokenize any Monkey code, print the tokens in the REPL, parse the name in a `let` statement, recognise a `return` statement, and collect errors for a malformed `let`.

**Can't do yet:**

- Understand any value or expression: `Value` and `ReturnValue` are always nil.
- Parse anything that isn't a `let` or `return`: it's silently dropped, with no error.
- Show errors to the user: they're collected, but only the tests read them.
- Run anything: the REPL doesn't even use the parser yet, and there's no evaluator.

**Position in the plan:** Sessions 1 and 2 of 12 are done, with tests passing.

---

## Picking up: Session 3

The goal is the start of Pratt parsing: identifiers, integer literals, and the prefix operators `!` and `-`. Before looking at the book, try answering these. They connect what I've done to what's next.

1. `5;` is currently dropped. What new kind of statement lets an expression stand alone?
2. `-` always means the same thing here, but in `5 - 3` it sits between two values. What does the parser need to know to tell those apart?
3. Instead of one big `switch`, how could the parser look up which function handles each token type?
4. What should the parser report for a token that can't start an expression? Which silent case does that replace?

---

## The road ahead

Each phase answers a bigger question:

| Sessions | Question it answers | Payoff |
|---|---|---|
| **3–6** | Can it understand expressions? | `Value` gets filled in, and the REPL prints trees |
| **7–9** | Can it run code? | Closures, which the plan calls the "aha" moment |
| **10–12** | Can it do more? | Strings, arrays, hash maps |
| **Week 5** | Can I make it mine? | My own feature, such as `spawn` with goroutines |

**Heads-up:** Session 4 (precedence) is the hard one. Draw it on paper first, and expect it to take two sittings.

---

## Update this after each session

Add a step to the story while it's fresh, and write it from memory first:

```
### Step N: <what I set out to do> (date)
What I did, in two or three sentences.
**What clicked here:** the one idea I'd want future me to remember.
**What tripped me up:** the bug or confusion, and what resolved it.
```

Writing it from memory is the recall practice. If I can't write it, I haven't got it yet.

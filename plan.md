Week 1: parsing, part 1

Session 1: AST (Node, Statement, Expression interfaces) and parsing let statements
Session 2: return statements and parser error collection
Session 3: Pratt parser intro: identifiers, integer literals, prefix operators (!, -)

Week 2: parsing, part 2

Session 4: infix operators and precedence. This is the hardest section. Draw the precedence climbing on paper before you code it.
Session 5: booleans, grouped expressions, if/else
Session 6: function literals, call expressions, then switch the REPL to print the parsed AST

Week 3: evaluation (chapter 3)

Session 7: object system, integers, booleans, prefix and infix evaluation
Session 8: conditionals, return values, error handling
Session 9: environments, let bindings, functions and closures (the “aha” moment of the book)

Week 4: extending the interpreter (chapter 4)

Session 10: strings and built-in functions (len)
Session 11: arrays, indexing, first/rest/push
Session 12: hash maps, then a README with REPL examples

Week 5 (optional): your own extension. A spawn built-in using goroutines and channels, or a decimal type. This is what makes the repo yours rather than a copy of the book.

Rules for pacing

One section, one commit, tests green. Keep doing what you did tonight.
If a session runs over, split it. Don’t skip it. Precedence parsing often takes two sessions; that’s normal.
Read Go by Example alongside: its interfaces chapter suits the AST work in week 1, and its testing chapter suits week 2.
Type the code yourself and understand each function before moving on. Don’t paste from the book or from AI. The point of this project is building Go fluency in your hands.
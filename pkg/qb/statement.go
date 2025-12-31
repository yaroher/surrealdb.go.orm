package qb

import "strings"

// Statement is a renderable SurrealQL statement.
type Statement interface {
	build(*Builder)
}

// Build renders a statement into a Query.
func Build(stmt Statement) Query {
	b := NewBuilder()
	stmt.build(b)
	return Query{Text: strings.TrimSpace(b.String()), Args: b.Args()}
}

// Chain is a list of statements rendered sequentially.
type Chain struct {
	Statements []Statement
}

func QueryChain(stmts ...Statement) Chain {
	return Chain{Statements: stmts}
}

func (c Chain) build(b *Builder) {
	for i, stmt := range c.Statements {
		if i > 0 {
			b.Write("; ")
		}
		stmt.build(b)
	}
}

func (c Chain) Build() Query {
	return Build(c)
}

// RawStatement injects raw SurrealQL statements.
type RawStatement struct {
	Text string
	Args map[string]any
}

func RawStmt(text string, args map[string]any) RawStatement {
	return RawStatement{Text: text, Args: args}
}

func (r RawStatement) build(b *Builder) {
	b.Write(r.Text)
	for name, val := range r.Args {
		b.Bind(name, val)
	}
}

// Subquery wraps a statement as an expression.
type Subquery struct {
	Stmt Statement
}

func (s Subquery) build(b *Builder) {
	b.Write("(")
	s.Stmt.build(b)
	b.Write(")")
}

// StmtExpr wraps a statement as an expression without parentheses.
type StmtExpr struct {
	Stmt Statement
}

func StatementExpr(stmt Statement) StmtExpr {
	return StmtExpr{Stmt: stmt}
}

func (s StmtExpr) build(b *Builder) {
	s.Stmt.build(b)
}

// Block wraps a body in braces.
type Block struct {
	Body Node
}

func BlockOf(body Node) Block {
	return Block{Body: body}
}

func (b Block) build(builder *Builder) {
	builder.Write("{")
	if b.Body != nil {
		builder.Write(" ")
		b.Body.build(builder)
		builder.Write(" ")
	}
	builder.Write("}")
}

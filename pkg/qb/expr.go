package qb

import "fmt"

// Node is a renderable AST node.
type Node interface {
	build(*Builder)
}

// Expr is a typed expression.
type Expr[T any] struct {
	node Node
}

func (e Expr[T]) build(b *Builder) {
	e.node.build(b)
}

// Condition is a boolean expression.
type Condition = Expr[bool]

// RawExpr injects raw SurrealQL.
type RawExpr struct {
	Text string
}

func Raw(text string) Expr[any] {
	return Expr[any]{node: RawExpr{Text: text}}
}

// RawCond injects a raw boolean condition.
func RawCond(text string) Condition {
	return Expr[bool]{node: RawExpr{Text: text}}
}

func (r RawExpr) build(b *Builder) {
	b.Write(r.Text)
}

// Value represents a bound value.
type Value struct {
	Val any
}

func V(val any) Expr[any] {
	return Expr[any]{node: Value{Val: val}}
}

func (v Value) build(b *Builder) {
	b.Write(b.Arg(v.Val))
}

// Param is a named parameter.
type Param struct {
	Name  string
	Value any
	Has   bool
}

func P(name string) Expr[any] {
	return Expr[any]{node: Param{Name: name}}
}

func PWith(name string, value any) Expr[any] {
	return Expr[any]{node: Param{Name: name, Value: value, Has: true}}
}

func (p Param) build(b *Builder) {
	if p.Has {
		b.Bind(p.Name, p.Value)
	}
	b.Write("$")
	b.Write(p.Name)
}

// Ident represents an identifier.
type Ident struct {
	Name string
}

func I(name string) Expr[any] {
	return Expr[any]{node: Ident{Name: name}}
}

func (i Ident) build(b *Builder) {
	b.Write(i.Name)
}

// Alias represents "expr AS alias".
type Alias struct {
	Expr  Node
	Alias string
}

func As(expr Node, alias string) Expr[any] {
	return Expr[any]{node: Alias{Expr: expr, Alias: alias}}
}

func (a Alias) build(b *Builder) {
	a.Expr.build(b)
	b.Write(" AS ")
	b.Write(a.Alias)
}

// Binary is a binary operator node.
type Binary struct {
	Left  Node
	Op    string
	Right Node
}

func (bop Binary) build(b *Builder) {
	b.Write("(")
	bop.Left.build(b)
	b.Write(" ")
	b.Write(bop.Op)
	b.Write(" ")
	bop.Right.build(b)
	b.Write(")")
}

// Unary is a unary operator node.
type Unary struct {
	Op   string
	Expr Node
}

func (u Unary) build(b *Builder) {
	b.Write(u.Op)
	u.Expr.build(b)
}

// FuncCall represents "fn(arg1, arg2, ...)".
type FuncCall struct {
	Name string
	Args []Node
}

func Fn(name string, args ...Node) Expr[any] {
	return Expr[any]{node: FuncCall{Name: name, Args: args}}
}

func (f FuncCall) build(b *Builder) {
	b.Write(f.Name)
	b.Write("(")
	for i, arg := range f.Args {
		if i > 0 {
			b.Write(", ")
		}
		arg.build(b)
	}
	b.Write(")")
}

// List renders "[a, b, c]".
type List struct {
	Items []Node
}

func L(items ...Node) Expr[any] {
	return Expr[any]{node: List{Items: items}}
}

func (l List) build(b *Builder) {
	b.Write("[")
	for i, it := range l.Items {
		if i > 0 {
			b.Write(", ")
		}
		it.build(b)
	}
	b.Write("]")
}

func ensureNode(v any) Node {
	switch t := v.(type) {
	case Node:
		return t
	case Expr[any]:
		return t
	default:
		return Value{Val: t}
	}
}

func debugNode(n Node) string {
	b := NewBuilder()
	n.build(b)
	return fmt.Sprintf("%s %v", b.String(), b.Args())
}

func trimParamName(name string) string {
	if len(name) == 0 {
		return name
	}
	if name[0] == '$' {
		return name[1:]
	}
	return name
}

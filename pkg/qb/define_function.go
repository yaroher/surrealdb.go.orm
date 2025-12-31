package qb

import "strings"

// DefineFunctionStatement builds DEFINE FUNCTION.
type DefineFunctionStatement struct {
	Name   string
	Params []string
	Body   Node
}

func DefineFunction(name string, params ...string) *DefineFunctionStatement {
	return &DefineFunctionStatement{Name: name, Params: params}
}

func (d *DefineFunctionStatement) BodyExpr(expr Node) *DefineFunctionStatement {
	d.Body = expr
	return d
}

func (d *DefineFunctionStatement) build(b *Builder) {
	b.Write("DEFINE FUNCTION ")
	b.Write(d.Name)
	b.Write("(")
	for i, p := range d.Params {
		if i > 0 {
			b.Write(", ")
		}
		b.Write("$")
		b.Write(trimParamName(p))
	}
	b.Write(") ")
	BlockOf(d.Body).build(b)
}

func (d *DefineFunctionStatement) Build() Query {
	return Build(d)
}

// Helper to build function signature from a slice.
func FuncSignature(name string, params []string) string {
	var sb strings.Builder
	sb.WriteString(name)
	sb.WriteString("(")
	for i, p := range params {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString("$")
		sb.WriteString(trimParamName(p))
	}
	sb.WriteString(")")
	return sb.String()
}

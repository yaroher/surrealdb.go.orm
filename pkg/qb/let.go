package qb

// LetStatement builds LET statements.
type LetStatement struct {
	Name  string
	Value Node
}

// Let declares a variable with a value.
func Let(name string, value any) *LetStatement {
	return &LetStatement{Name: name, Value: ensureValueNode(value)}
}

func (l *LetStatement) build(b *Builder) {
	b.Write("LET $")
	b.Write(trimParamName(l.Name))
	b.Write(" = ")
	l.Value.build(b)
}

func (l *LetStatement) Build() Query {
	return Build(l)
}

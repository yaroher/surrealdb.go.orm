package qb

// BreakStatement builds BREAK.
type BreakStatement struct{}

func Break() BreakStatement {
	return BreakStatement{}
}

func (b BreakStatement) build(builder *Builder) {
	builder.Write("BREAK")
}

// ContinueStatement builds CONTINUE.
type ContinueStatement struct{}

func Continue() ContinueStatement {
	return ContinueStatement{}
}

func (c ContinueStatement) build(builder *Builder) {
	builder.Write("CONTINUE")
}

// ThrowStatement builds THROW <expr>.
type ThrowStatement struct {
	Value Node
}

func Throw(value any) *ThrowStatement {
	return &ThrowStatement{Value: ensureValueNode(value)}
}

func (t *ThrowStatement) build(b *Builder) {
	b.Write("THROW ")
	t.Value.build(b)
}

func (t *ThrowStatement) Build() Query {
	return Build(t)
}

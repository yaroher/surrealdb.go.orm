package qb

// ReturnStatement builds RETURN statements.
type ReturnStatement struct {
	Value Node
}

func Return(value any) *ReturnStatement {
	return &ReturnStatement{Value: ensureValueNode(value)}
}

func (r *ReturnStatement) build(b *Builder) {
	b.Write("RETURN ")
	r.Value.build(b)
}

func (r *ReturnStatement) Build() Query {
	return Build(r)
}

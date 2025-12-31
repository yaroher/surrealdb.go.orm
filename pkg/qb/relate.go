package qb

// RelateBuilder builds RELATE statements.
type RelateBuilder struct {
	from      Node
	edge      Node
	to        Node
	set       []Assignment
	returning Node
}

func Relate(from Node, edge Node, to Node) *RelateBuilder {
	return &RelateBuilder{from: from, edge: edge, to: to}
}

func (r *RelateBuilder) Set(assignments ...Assignment) *RelateBuilder {
	r.set = assignments
	return r
}

func (r *RelateBuilder) Return(ret Node) *RelateBuilder {
	r.returning = ret
	return r
}

func (r *RelateBuilder) Build() Query {
	return Build(r)
}

func (r *RelateBuilder) build(b *Builder) {
	b.Write("RELATE ")
	r.from.build(b)
	b.Write(" -> ")
	r.edge.build(b)
	b.Write(" -> ")
	r.to.build(b)
	if len(r.set) > 0 {
		b.Write(" SET ")
		renderAssignments(b, r.set)
	}
	if r.returning != nil {
		b.Write(" RETURN ")
		r.returning.build(b)
	}
}

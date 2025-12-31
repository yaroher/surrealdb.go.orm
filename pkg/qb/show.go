package qb

// ShowChangesStatement builds SHOW CHANGES statements.
type ShowChangesStatement struct {
	Table Node
	since Node
	limit Node
}

func ShowChangesForTable(table Node) *ShowChangesStatement {
	return &ShowChangesStatement{Table: table}
}

func (s *ShowChangesStatement) SinceExpr(value Node) *ShowChangesStatement {
	s.since = value
	return s
}

func (s *ShowChangesStatement) Since(value any) *ShowChangesStatement {
	s.since = ensureValueNode(value)
	return s
}

func (s *ShowChangesStatement) Limit(value any) *ShowChangesStatement {
	s.limit = ensureValueNode(value)
	return s
}

func (s *ShowChangesStatement) build(b *Builder) {
	b.Write("SHOW CHANGES FOR TABLE ")
	s.Table.build(b)
	if s.since != nil {
		b.Write(" SINCE ")
		s.since.build(b)
	}
	if s.limit != nil {
		b.Write(" LIMIT ")
		s.limit.build(b)
	}
}

func (s *ShowChangesStatement) Build() Query {
	return Build(s)
}

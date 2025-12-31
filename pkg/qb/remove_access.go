package qb

// RemoveAccessStatement builds REMOVE ACCESS.
type RemoveAccessStatement struct {
	Name     Node
	ifExists bool
	Scope    AccessScope
}

func RemoveAccess(name string) *RemoveAccessStatement {
	return &RemoveAccessStatement{Name: Ident{Name: name}}
}

func RemoveAccessExpr(expr Node) *RemoveAccessStatement {
	return &RemoveAccessStatement{Name: expr}
}

func (r *RemoveAccessStatement) IfExists() *RemoveAccessStatement {
	r.ifExists = true
	return r
}

func (r *RemoveAccessStatement) OnNamespace() *RemoveAccessStatement {
	r.Scope = AccessNamespace
	return r
}

func (r *RemoveAccessStatement) OnDatabase() *RemoveAccessStatement {
	r.Scope = AccessDatabase
	return r
}

func (r *RemoveAccessStatement) build(b *Builder) {
	b.Write("REMOVE ACCESS ")
	if r.ifExists {
		b.Write("IF EXISTS ")
	}
	r.Name.build(b)
	if r.Scope != "" {
		b.Write(" ON ")
		b.Write(string(r.Scope))
	}
}

func (r *RemoveAccessStatement) Build() Query {
	return Build(r)
}

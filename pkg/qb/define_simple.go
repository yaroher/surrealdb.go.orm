package qb

// DefineNamespaceStatement builds DEFINE NAMESPACE.
type DefineNamespaceStatement struct {
	Name Node
}

func DefineNamespace(name string) *DefineNamespaceStatement {
	return &DefineNamespaceStatement{Name: Ident{Name: name}}
}

func DefineNamespaceExpr(expr Node) *DefineNamespaceStatement {
	return &DefineNamespaceStatement{Name: expr}
}

func (d *DefineNamespaceStatement) build(b *Builder) {
	b.Write("DEFINE NAMESPACE ")
	d.Name.build(b)
}

func (d *DefineNamespaceStatement) Build() Query {
	return Build(d)
}

// DefineDatabaseStatement builds DEFINE DATABASE.
type DefineDatabaseStatement struct {
	Name Node
}

func DefineDatabase(name string) *DefineDatabaseStatement {
	return &DefineDatabaseStatement{Name: Ident{Name: name}}
}

func DefineDatabaseExpr(expr Node) *DefineDatabaseStatement {
	return &DefineDatabaseStatement{Name: expr}
}

func (d *DefineDatabaseStatement) build(b *Builder) {
	b.Write("DEFINE DATABASE ")
	d.Name.build(b)
}

func (d *DefineDatabaseStatement) Build() Query {
	return Build(d)
}

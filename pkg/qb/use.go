package qb

// UseBuilder builds USE statements.
type UseBuilder struct {
	ns Node
	db Node
}

// Use starts a USE statement.
func Use() *UseBuilder {
	return &UseBuilder{}
}

// Namespace sets the namespace.
func (u *UseBuilder) Namespace(name string) *UseBuilder {
	u.ns = Ident{Name: name}
	return u
}

// Database sets the database.
func (u *UseBuilder) Database(name string) *UseBuilder {
	u.db = Ident{Name: name}
	return u
}

// NamespaceExpr sets the namespace expression.
func (u *UseBuilder) NamespaceExpr(expr Node) *UseBuilder {
	u.ns = expr
	return u
}

// DatabaseExpr sets the database expression.
func (u *UseBuilder) DatabaseExpr(expr Node) *UseBuilder {
	u.db = expr
	return u
}

func (u *UseBuilder) build(b *Builder) {
	b.Write("USE")
	if u.ns != nil {
		b.Write(" NS ")
		u.ns.build(b)
	}
	if u.db != nil {
		b.Write(" DB ")
		u.db.build(b)
	}
}

func (u *UseBuilder) Build() Query {
	return Build(u)
}

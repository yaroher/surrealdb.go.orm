package qb

// DefineEventStatement builds DEFINE EVENT.
type DefineEventStatement struct {
	Name     Node
	Table    Node
	whenCond Condition
	thenExpr Node
}

func DefineEvent(name string) *DefineEventStatement {
	return &DefineEventStatement{Name: Ident{Name: name}}
}

func DefineEventExpr(expr Node) *DefineEventStatement {
	return &DefineEventStatement{Name: expr}
}

func (d *DefineEventStatement) OnTable(table Node) *DefineEventStatement {
	d.Table = table
	return d
}

func (d *DefineEventStatement) OnTableName(name string) *DefineEventStatement {
	d.Table = Ident{Name: name}
	return d
}

func (d *DefineEventStatement) When(cond Condition) *DefineEventStatement {
	d.whenCond = cond
	return d
}

func (d *DefineEventStatement) ThenExpr(expr Node) *DefineEventStatement {
	d.thenExpr = expr
	return d
}

func (d *DefineEventStatement) ThenValue(value any) *DefineEventStatement {
	d.thenExpr = ensureValueNode(value)
	return d
}

func (d *DefineEventStatement) build(b *Builder) {
	b.Write("DEFINE EVENT ")
	d.Name.build(b)
	if d.Table != nil {
		b.Write(" ON TABLE ")
		d.Table.build(b)
	}
	if d.whenCond.node != nil {
		b.Write(" WHEN ")
		d.whenCond.build(b)
	}
	if d.thenExpr != nil {
		b.Write(" THEN ")
		d.thenExpr.build(b)
	}
}

func (d *DefineEventStatement) Build() Query {
	return Build(d)
}

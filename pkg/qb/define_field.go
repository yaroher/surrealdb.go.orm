package qb

// DefineFieldStatement builds DEFINE FIELD.
type DefineFieldStatement struct {
	Field        Node
	Table        Node
	fieldType    Node
	fieldValue   Node
	assertCond   Condition
	defaultValue Node
	Permissions  *Permissions
}

func DefineField(field Node, table Node) *DefineFieldStatement {
	return &DefineFieldStatement{Field: field, Table: table}
}

func DefineFieldName(field, table string) *DefineFieldStatement {
	return DefineField(Ident{Name: field}, Ident{Name: table})
}

func (d *DefineFieldStatement) TypeExpr(expr Node) *DefineFieldStatement {
	d.fieldType = expr
	return d
}

func (d *DefineFieldStatement) Type(value string) *DefineFieldStatement {
	d.fieldType = Ident{Name: value}
	return d
}

func (d *DefineFieldStatement) ValueExpr(expr Node) *DefineFieldStatement {
	d.fieldValue = expr
	return d
}

func (d *DefineFieldStatement) Value(value any) *DefineFieldStatement {
	d.fieldValue = ensureValueNode(value)
	return d
}

func (d *DefineFieldStatement) Assert(cond Condition) *DefineFieldStatement {
	d.assertCond = cond
	return d
}

func (d *DefineFieldStatement) DefaultExpr(expr Node) *DefineFieldStatement {
	d.defaultValue = expr
	return d
}

func (d *DefineFieldStatement) Default(value any) *DefineFieldStatement {
	d.defaultValue = ensureValueNode(value)
	return d
}

func (d *DefineFieldStatement) PermissionsNone() *DefineFieldStatement {
	if d.Permissions == nil {
		d.Permissions = &Permissions{}
	}
	d.Permissions.NoneOnly()
	return d
}

func (d *DefineFieldStatement) PermissionsFull() *DefineFieldStatement {
	if d.Permissions == nil {
		d.Permissions = &Permissions{}
	}
	d.Permissions.FullOnly()
	return d
}

func (d *DefineFieldStatement) PermissionsFor(lines ...Node) *DefineFieldStatement {
	if d.Permissions == nil {
		d.Permissions = &Permissions{}
	}
	d.Permissions.With(lines...)
	return d
}

func (d *DefineFieldStatement) build(b *Builder) {
	b.Write("DEFINE FIELD ")
	d.Field.build(b)
	b.Write(" ON TABLE ")
	d.Table.build(b)
	if d.fieldType != nil {
		b.Write(" TYPE ")
		d.fieldType.build(b)
	}
	if d.fieldValue != nil {
		b.Write(" VALUE ")
		d.fieldValue.build(b)
	}
	if d.assertCond.node != nil {
		b.Write(" ASSERT ")
		d.assertCond.build(b)
	}
	if d.defaultValue != nil {
		b.Write(" DEFAULT ")
		d.defaultValue.build(b)
	}
	renderPermissions(b, d.Permissions)
}

func (d *DefineFieldStatement) Build() Query {
	return Build(d)
}

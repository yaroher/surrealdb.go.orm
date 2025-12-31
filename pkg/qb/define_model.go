package qb

// DefineModelStatement builds DEFINE MODEL.
type DefineModelStatement struct {
	Name        Node
	Version     Node
	Comment     Node
	Permissions *Permissions
}

func DefineModel(name string) *DefineModelStatement {
	return &DefineModelStatement{Name: Ident{Name: name}}
}

func DefineModelExpr(expr Node) *DefineModelStatement {
	return &DefineModelStatement{Name: expr}
}

func (d *DefineModelStatement) VersionExpr(expr Node) *DefineModelStatement {
	d.Version = expr
	return d
}

func (d *DefineModelStatement) VersionValue(value any) *DefineModelStatement {
	d.Version = ensureValueNode(value)
	return d
}

func (d *DefineModelStatement) CommentExpr(expr Node) *DefineModelStatement {
	d.Comment = expr
	return d
}

func (d *DefineModelStatement) CommentValue(value any) *DefineModelStatement {
	d.Comment = ensureValueNode(value)
	return d
}

func (d *DefineModelStatement) PermissionsNone() *DefineModelStatement {
	if d.Permissions == nil {
		d.Permissions = &Permissions{}
	}
	d.Permissions.NoneOnly()
	return d
}

func (d *DefineModelStatement) PermissionsFull() *DefineModelStatement {
	if d.Permissions == nil {
		d.Permissions = &Permissions{}
	}
	d.Permissions.FullOnly()
	return d
}

func (d *DefineModelStatement) PermissionsFor(lines ...Node) *DefineModelStatement {
	if d.Permissions == nil {
		d.Permissions = &Permissions{}
	}
	d.Permissions.With(lines...)
	return d
}

func (d *DefineModelStatement) build(b *Builder) {
	b.Write("DEFINE MODEL ml::")
	d.Name.build(b)
	if d.Version != nil {
		b.Write("<")
		d.Version.build(b)
		b.Write(">")
	}
	if d.Comment != nil {
		b.Write("\nCOMMENT ")
		d.Comment.build(b)
	}
	renderPermissions(b, d.Permissions)
}

func (d *DefineModelStatement) Build() Query {
	return Build(d)
}

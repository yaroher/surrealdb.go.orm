package qb

// SchemaType defines schema mode.
type SchemaType string

const (
	SchemaFull  SchemaType = "SCHEMAFULL"
	SchemaLess  SchemaType = "SCHEMALESS"
	SchemaUnset SchemaType = ""
)

// DefineTableStatement builds DEFINE TABLE.
type DefineTableStatement struct {
	Table       Node
	Drop        bool
	Schema      SchemaType
	As          Statement
	Permissions *Permissions
}

func DefineTable(table Node) *DefineTableStatement {
	return &DefineTableStatement{Table: table}
}

func DefineTableName(name string) *DefineTableStatement {
	return DefineTable(Ident{Name: name})
}

func (d *DefineTableStatement) DropTable() *DefineTableStatement {
	d.Drop = true
	return d
}

func (d *DefineTableStatement) SchemaFull() *DefineTableStatement {
	d.Schema = SchemaFull
	return d
}

func (d *DefineTableStatement) SchemaLess() *DefineTableStatement {
	d.Schema = SchemaLess
	return d
}

func (d *DefineTableStatement) AsSelect(stmt Statement) *DefineTableStatement {
	d.As = stmt
	return d
}

func (d *DefineTableStatement) PermissionsNone() *DefineTableStatement {
	if d.Permissions == nil {
		d.Permissions = &Permissions{}
	}
	d.Permissions.NoneOnly()
	return d
}

func (d *DefineTableStatement) PermissionsFull() *DefineTableStatement {
	if d.Permissions == nil {
		d.Permissions = &Permissions{}
	}
	d.Permissions.FullOnly()
	return d
}

func (d *DefineTableStatement) PermissionsFor(lines ...Node) *DefineTableStatement {
	if d.Permissions == nil {
		d.Permissions = &Permissions{}
	}
	d.Permissions.With(lines...)
	return d
}

func (d *DefineTableStatement) build(b *Builder) {
	b.Write("DEFINE TABLE ")
	d.Table.build(b)
	if d.Drop {
		b.Write(" DROP")
	}
	if d.Schema != SchemaUnset {
		b.Write(" ")
		b.Write(string(d.Schema))
	}
	if d.As != nil {
		b.Write(" AS ")
		d.As.build(b)
	}
	renderPermissions(b, d.Permissions)
}

func (d *DefineTableStatement) Build() Query {
	return Build(d)
}

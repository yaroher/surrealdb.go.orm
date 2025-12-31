package qb

// UserScope defines scope for DEFINE USER.
type UserScope string

const (
	UserRoot      UserScope = "ROOT"
	UserNamespace UserScope = "NAMESPACE"
	UserDatabase  UserScope = "DATABASE"
)

// DefineUserStatement builds DEFINE USER.
type DefineUserStatement struct {
	Name     Node
	On       UserScope
	password Node
	Roles    []string
}

func DefineUser(name string) *DefineUserStatement {
	return &DefineUserStatement{Name: Ident{Name: name}}
}

func DefineUserExpr(expr Node) *DefineUserStatement {
	return &DefineUserStatement{Name: expr}
}

func (d *DefineUserStatement) OnRoot() *DefineUserStatement {
	d.On = UserRoot
	return d
}

func (d *DefineUserStatement) OnNamespace() *DefineUserStatement {
	d.On = UserNamespace
	return d
}

func (d *DefineUserStatement) OnDatabase() *DefineUserStatement {
	d.On = UserDatabase
	return d
}

func (d *DefineUserStatement) Password(value any) *DefineUserStatement {
	d.password = ensureValueNode(value)
	return d
}

func (d *DefineUserStatement) RolesList(roles ...string) *DefineUserStatement {
	d.Roles = append(d.Roles, roles...)
	return d
}

func (d *DefineUserStatement) build(b *Builder) {
	b.Write("DEFINE USER ")
	d.Name.build(b)
	if d.On != "" {
		b.Write(" ON ")
		b.Write(string(d.On))
	}
	if d.password != nil {
		b.Write(" PASSWORD ")
		d.password.build(b)
	}
	if len(d.Roles) > 0 {
		b.Write(" ROLES ")
		for i, r := range d.Roles {
			if i > 0 {
				b.Write(", ")
			}
			b.Write(r)
		}
	}
}

func (d *DefineUserStatement) Build() Query {
	return Build(d)
}

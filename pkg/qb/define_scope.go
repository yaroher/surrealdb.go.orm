package qb

// DefineScopeStatement builds DEFINE SCOPE.
type DefineScopeStatement struct {
	Name    Node
	Session Node
	Signup  Node
	Signin  Node
}

func DefineScope(name string) *DefineScopeStatement {
	return &DefineScopeStatement{Name: Ident{Name: name}}
}

func DefineScopeExpr(expr Node) *DefineScopeStatement {
	return &DefineScopeStatement{Name: expr}
}

func (d *DefineScopeStatement) SessionExpr(expr Node) *DefineScopeStatement {
	d.Session = expr
	return d
}

func (d *DefineScopeStatement) SessionValue(value any) *DefineScopeStatement {
	d.Session = ensureValueNode(value)
	return d
}

func (d *DefineScopeStatement) SignupExpr(expr Node) *DefineScopeStatement {
	d.Signup = expr
	return d
}

func (d *DefineScopeStatement) SignupValue(value any) *DefineScopeStatement {
	d.Signup = ensureValueNode(value)
	return d
}

func (d *DefineScopeStatement) SigninExpr(expr Node) *DefineScopeStatement {
	d.Signin = expr
	return d
}

func (d *DefineScopeStatement) SigninValue(value any) *DefineScopeStatement {
	d.Signin = ensureValueNode(value)
	return d
}

func (d *DefineScopeStatement) build(b *Builder) {
	b.Write("DEFINE SCOPE ")
	d.Name.build(b)
	if d.Session != nil {
		b.Write(" SESSION ")
		d.Session.build(b)
	}
	if d.Signup != nil {
		b.Write(" SIGNUP ")
		BlockOf(d.Signup).build(b)
	}
	if d.Signin != nil {
		b.Write(" SIGNIN ")
		BlockOf(d.Signin).build(b)
	}
}

func (d *DefineScopeStatement) Build() Query {
	return Build(d)
}

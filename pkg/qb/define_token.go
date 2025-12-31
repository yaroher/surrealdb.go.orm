package qb

// TokenScope defines token scope.
type TokenScope string

const (
	TokenNamespace TokenScope = "NAMESPACE"
	TokenDatabase  TokenScope = "DATABASE"
	TokenScopeName TokenScope = "SCOPE"
)

// DefineTokenStatement builds DEFINE TOKEN.
type DefineTokenStatement struct {
	Name       Node
	On         TokenScope
	Scope      Node
	tokenType  Node
	tokenValue Node
}

func DefineToken(name string) *DefineTokenStatement {
	return &DefineTokenStatement{Name: Ident{Name: name}}
}

func DefineTokenExpr(expr Node) *DefineTokenStatement {
	return &DefineTokenStatement{Name: expr}
}

func (d *DefineTokenStatement) OnNamespace() *DefineTokenStatement {
	d.On = TokenNamespace
	return d
}

func (d *DefineTokenStatement) OnDatabase() *DefineTokenStatement {
	d.On = TokenDatabase
	return d
}

func (d *DefineTokenStatement) OnScope(name string) *DefineTokenStatement {
	d.On = TokenScopeName
	d.Scope = Ident{Name: name}
	return d
}

func (d *DefineTokenStatement) OnScopeExpr(expr Node) *DefineTokenStatement {
	d.On = TokenScopeName
	d.Scope = expr
	return d
}

func (d *DefineTokenStatement) TypeExpr(expr Node) *DefineTokenStatement {
	d.tokenType = expr
	return d
}

func (d *DefineTokenStatement) Type(name string) *DefineTokenStatement {
	d.tokenType = Ident{Name: name}
	return d
}

func (d *DefineTokenStatement) ValueExpr(expr Node) *DefineTokenStatement {
	d.tokenValue = expr
	return d
}

func (d *DefineTokenStatement) Value(value any) *DefineTokenStatement {
	d.tokenValue = ensureValueNode(value)
	return d
}

func (d *DefineTokenStatement) build(b *Builder) {
	b.Write("DEFINE TOKEN ")
	d.Name.build(b)
	if d.On != "" {
		b.Write(" ON ")
		b.Write(string(d.On))
		if d.On == TokenScopeName && d.Scope != nil {
			b.Write(" ")
			d.Scope.build(b)
		}
	}
	if d.tokenType != nil {
		b.Write(" TYPE ")
		d.tokenType.build(b)
	}
	if d.tokenValue != nil {
		b.Write(" VALUE ")
		d.tokenValue.build(b)
	}
}

func (d *DefineTokenStatement) Build() Query {
	return Build(d)
}

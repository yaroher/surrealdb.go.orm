package qb

// AccessGrantStatement builds ACCESS GRANT.
type AccessGrantStatement struct {
	Access   Node
	On       AccessScope
	subject  Node
	token    Node
	duration Node
}

func AccessGrant(name string) *AccessGrantStatement {
	return &AccessGrantStatement{Access: Ident{Name: name}}
}

func AccessGrantExpr(expr Node) *AccessGrantStatement {
	return &AccessGrantStatement{Access: expr}
}

func (a *AccessGrantStatement) OnNamespace() *AccessGrantStatement {
	a.On = AccessNamespace
	return a
}

func (a *AccessGrantStatement) OnDatabase() *AccessGrantStatement {
	a.On = AccessDatabase
	return a
}

// Subject sets FOR <subject>.
func (a *AccessGrantStatement) Subject(expr Node) *AccessGrantStatement {
	a.subject = expr
	return a
}

// Token sets TOKEN <expr>.
func (a *AccessGrantStatement) Token(expr Node) *AccessGrantStatement {
	a.token = expr
	return a
}

// Duration sets DURATION <expr>.
func (a *AccessGrantStatement) Duration(expr Node) *AccessGrantStatement {
	a.duration = expr
	return a
}

func (a *AccessGrantStatement) build(b *Builder) {
	b.Write("ACCESS ")
	a.Access.build(b)
	if a.On != "" {
		b.Write(" ON ")
		b.Write(string(a.On))
	}
	b.Write(" GRANT")
	if a.subject != nil {
		b.Write(" FOR ")
		a.subject.build(b)
	}
	if a.token != nil {
		b.Write(" TOKEN ")
		a.token.build(b)
	}
	if a.duration != nil {
		b.Write(" DURATION ")
		a.duration.build(b)
	}
}

func (a *AccessGrantStatement) Build() Query {
	return Build(a)
}

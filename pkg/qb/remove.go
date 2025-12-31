package qb

// RemoveNamespaceStatement builds REMOVE NAMESPACE.
type RemoveNamespaceStatement struct {
	Name Node
}

func RemoveNamespace(name string) *RemoveNamespaceStatement {
	return &RemoveNamespaceStatement{Name: Ident{Name: name}}
}

func RemoveNamespaceExpr(expr Node) *RemoveNamespaceStatement {
	return &RemoveNamespaceStatement{Name: expr}
}

func (r *RemoveNamespaceStatement) build(b *Builder) {
	b.Write("REMOVE NAMESPACE ")
	r.Name.build(b)
}

func (r *RemoveNamespaceStatement) Build() Query {
	return Build(r)
}

// RemoveDatabaseStatement builds REMOVE DATABASE.
type RemoveDatabaseStatement struct {
	Name Node
}

func RemoveDatabase(name string) *RemoveDatabaseStatement {
	return &RemoveDatabaseStatement{Name: Ident{Name: name}}
}

func RemoveDatabaseExpr(expr Node) *RemoveDatabaseStatement {
	return &RemoveDatabaseStatement{Name: expr}
}

func (r *RemoveDatabaseStatement) build(b *Builder) {
	b.Write("REMOVE DATABASE ")
	r.Name.build(b)
}

func (r *RemoveDatabaseStatement) Build() Query {
	return Build(r)
}

// RemoveTableStatement builds REMOVE TABLE.
type RemoveTableStatement struct {
	Name Node
}

func RemoveTable(name string) *RemoveTableStatement {
	return &RemoveTableStatement{Name: Ident{Name: name}}
}

func RemoveTableExpr(expr Node) *RemoveTableStatement {
	return &RemoveTableStatement{Name: expr}
}

func (r *RemoveTableStatement) build(b *Builder) {
	b.Write("REMOVE TABLE ")
	r.Name.build(b)
}

func (r *RemoveTableStatement) Build() Query {
	return Build(r)
}

// RemoveFieldStatement builds REMOVE FIELD.
type RemoveFieldStatement struct {
	Name  Node
	Table Node
}

func RemoveField(name string) *RemoveFieldStatement {
	return &RemoveFieldStatement{Name: Ident{Name: name}}
}

func RemoveFieldExpr(expr Node) *RemoveFieldStatement {
	return &RemoveFieldStatement{Name: expr}
}

func (r *RemoveFieldStatement) OnTable(table Node) *RemoveFieldStatement {
	r.Table = table
	return r
}

func (r *RemoveFieldStatement) OnTableName(name string) *RemoveFieldStatement {
	r.Table = Ident{Name: name}
	return r
}

func (r *RemoveFieldStatement) build(b *Builder) {
	b.Write("REMOVE FIELD ")
	r.Name.build(b)
	if r.Table != nil {
		b.Write(" ON TABLE ")
		r.Table.build(b)
	}
}

func (r *RemoveFieldStatement) Build() Query {
	return Build(r)
}

// RemoveIndexStatement builds REMOVE INDEX.
type RemoveIndexStatement struct {
	Name  Node
	Table Node
}

func RemoveIndex(name string) *RemoveIndexStatement {
	return &RemoveIndexStatement{Name: Ident{Name: name}}
}

func RemoveIndexExpr(expr Node) *RemoveIndexStatement {
	return &RemoveIndexStatement{Name: expr}
}

func (r *RemoveIndexStatement) OnTable(table Node) *RemoveIndexStatement {
	r.Table = table
	return r
}

func (r *RemoveIndexStatement) OnTableName(name string) *RemoveIndexStatement {
	r.Table = Ident{Name: name}
	return r
}

func (r *RemoveIndexStatement) build(b *Builder) {
	b.Write("REMOVE INDEX ")
	r.Name.build(b)
	if r.Table != nil {
		b.Write(" ON TABLE ")
		r.Table.build(b)
	}
}

func (r *RemoveIndexStatement) Build() Query {
	return Build(r)
}

// RemoveEventStatement builds REMOVE EVENT.
type RemoveEventStatement struct {
	Name  Node
	Table Node
}

func RemoveEvent(name string) *RemoveEventStatement {
	return &RemoveEventStatement{Name: Ident{Name: name}}
}

func RemoveEventExpr(expr Node) *RemoveEventStatement {
	return &RemoveEventStatement{Name: expr}
}

func (r *RemoveEventStatement) OnTable(table Node) *RemoveEventStatement {
	r.Table = table
	return r
}

func (r *RemoveEventStatement) OnTableName(name string) *RemoveEventStatement {
	r.Table = Ident{Name: name}
	return r
}

func (r *RemoveEventStatement) build(b *Builder) {
	b.Write("REMOVE EVENT ")
	r.Name.build(b)
	if r.Table != nil {
		b.Write(" ON TABLE ")
		r.Table.build(b)
	}
}

func (r *RemoveEventStatement) Build() Query {
	return Build(r)
}

// RemoveFunctionStatement builds REMOVE FUNCTION.
type RemoveFunctionStatement struct {
	Name string
}

func RemoveFunction(name string) *RemoveFunctionStatement {
	return &RemoveFunctionStatement{Name: name}
}

func (r *RemoveFunctionStatement) build(b *Builder) {
	b.Write("REMOVE FUNCTION ")
	b.Write(r.Name)
}

func (r *RemoveFunctionStatement) Build() Query {
	return Build(r)
}

// RemoveParamStatement builds REMOVE PARAM.
type RemoveParamStatement struct {
	Name string
}

func RemoveParam(name string) *RemoveParamStatement {
	return &RemoveParamStatement{Name: name}
}

func (r *RemoveParamStatement) build(b *Builder) {
	b.Write("REMOVE PARAM $")
	b.Write(trimParamName(r.Name))
}

func (r *RemoveParamStatement) Build() Query {
	return Build(r)
}

// RemoveScopeStatement builds REMOVE SCOPE.
type RemoveScopeStatement struct {
	Name Node
}

func RemoveScope(name string) *RemoveScopeStatement {
	return &RemoveScopeStatement{Name: Ident{Name: name}}
}

func RemoveScopeExpr(expr Node) *RemoveScopeStatement {
	return &RemoveScopeStatement{Name: expr}
}

func (r *RemoveScopeStatement) build(b *Builder) {
	b.Write("REMOVE SCOPE ")
	r.Name.build(b)
}

func (r *RemoveScopeStatement) Build() Query {
	return Build(r)
}

// RemoveTokenStatement builds REMOVE TOKEN.
type RemoveTokenStatement struct {
	Name  Node
	On    TokenScope
	Scope Node
}

func RemoveToken(name string) *RemoveTokenStatement {
	return &RemoveTokenStatement{Name: Ident{Name: name}}
}

func RemoveTokenExpr(expr Node) *RemoveTokenStatement {
	return &RemoveTokenStatement{Name: expr}
}

func (r *RemoveTokenStatement) OnNamespace() *RemoveTokenStatement {
	r.On = TokenNamespace
	return r
}

func (r *RemoveTokenStatement) OnDatabase() *RemoveTokenStatement {
	r.On = TokenDatabase
	return r
}

func (r *RemoveTokenStatement) OnScope(name string) *RemoveTokenStatement {
	r.On = TokenScopeName
	r.Scope = Ident{Name: name}
	return r
}

func (r *RemoveTokenStatement) OnScopeExpr(expr Node) *RemoveTokenStatement {
	r.On = TokenScopeName
	r.Scope = expr
	return r
}

func (r *RemoveTokenStatement) build(b *Builder) {
	b.Write("REMOVE TOKEN ")
	r.Name.build(b)
	if r.On != "" {
		b.Write(" ON ")
		b.Write(string(r.On))
		if r.On == TokenScopeName && r.Scope != nil {
			b.Write(" ")
			r.Scope.build(b)
		}
	}
}

func (r *RemoveTokenStatement) Build() Query {
	return Build(r)
}

// RemoveUserStatement builds REMOVE USER.
type RemoveUserStatement struct {
	Name Node
	On   UserScope
}

func RemoveUser(name string) *RemoveUserStatement {
	return &RemoveUserStatement{Name: Ident{Name: name}}
}

func RemoveUserExpr(expr Node) *RemoveUserStatement {
	return &RemoveUserStatement{Name: expr}
}

func (r *RemoveUserStatement) OnRoot() *RemoveUserStatement {
	r.On = UserRoot
	return r
}

func (r *RemoveUserStatement) OnNamespace() *RemoveUserStatement {
	r.On = UserNamespace
	return r
}

func (r *RemoveUserStatement) OnDatabase() *RemoveUserStatement {
	r.On = UserDatabase
	return r
}

func (r *RemoveUserStatement) build(b *Builder) {
	b.Write("REMOVE USER ")
	r.Name.build(b)
	if r.On != "" {
		b.Write(" ON ")
		b.Write(string(r.On))
	}
}

func (r *RemoveUserStatement) Build() Query {
	return Build(r)
}

// RemoveAnalyzerStatement builds REMOVE ANALYZER.
type RemoveAnalyzerStatement struct {
	Name Node
}

func RemoveAnalyzer(name string) *RemoveAnalyzerStatement {
	return &RemoveAnalyzerStatement{Name: Ident{Name: name}}
}

func RemoveAnalyzerExpr(expr Node) *RemoveAnalyzerStatement {
	return &RemoveAnalyzerStatement{Name: expr}
}

func (r *RemoveAnalyzerStatement) build(b *Builder) {
	b.Write("REMOVE ANALYZER ")
	r.Name.build(b)
}

func (r *RemoveAnalyzerStatement) Build() Query {
	return Build(r)
}

// RemoveLoginStatement builds REMOVE LOGIN.
type RemoveLoginStatement struct {
	Name Node
	On   TokenScope
}

func RemoveLogin(name string) *RemoveLoginStatement {
	return &RemoveLoginStatement{Name: Ident{Name: name}}
}

func RemoveLoginExpr(expr Node) *RemoveLoginStatement {
	return &RemoveLoginStatement{Name: expr}
}

func (r *RemoveLoginStatement) OnNamespace() *RemoveLoginStatement {
	r.On = TokenNamespace
	return r
}

func (r *RemoveLoginStatement) OnDatabase() *RemoveLoginStatement {
	r.On = TokenDatabase
	return r
}

func (r *RemoveLoginStatement) build(b *Builder) {
	b.Write("REMOVE LOGIN ")
	r.Name.build(b)
	if r.On != "" {
		b.Write(" ON ")
		b.Write(string(r.On))
	}
}

func (r *RemoveLoginStatement) Build() Query {
	return Build(r)
}

// RemoveModelStatement builds REMOVE MODEL.
type RemoveModelStatement struct {
	Name    Node
	Version Node
}

func RemoveModel(name string) *RemoveModelStatement {
	return &RemoveModelStatement{Name: Ident{Name: name}}
}

func RemoveModelExpr(expr Node) *RemoveModelStatement {
	return &RemoveModelStatement{Name: expr}
}

func (r *RemoveModelStatement) VersionExpr(expr Node) *RemoveModelStatement {
	r.Version = expr
	return r
}

func (r *RemoveModelStatement) VersionValue(value any) *RemoveModelStatement {
	r.Version = ensureValueNode(value)
	return r
}

func (r *RemoveModelStatement) build(b *Builder) {
	b.Write("REMOVE MODEL ml::")
	r.Name.build(b)
	if r.Version != nil {
		b.Write("<")
		r.Version.build(b)
		b.Write(">")
	}
}

func (r *RemoveModelStatement) Build() Query {
	return Build(r)
}

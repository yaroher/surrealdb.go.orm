package qb

// DefineParamStatement builds DEFINE PARAM.
type DefineParamStatement struct {
	Name  string
	Value Node
}

func DefineParam(name string, value any) *DefineParamStatement {
	return &DefineParamStatement{Name: name, Value: ensureValueNode(value)}
}

func (d *DefineParamStatement) ValueExpr(expr Node) *DefineParamStatement {
	d.Value = expr
	return d
}

func (d *DefineParamStatement) build(b *Builder) {
	b.Write("DEFINE PARAM $")
	b.Write(trimParamName(d.Name))
	b.Write(" VALUE ")
	d.Value.build(b)
}

func (d *DefineParamStatement) Build() Query {
	return Build(d)
}

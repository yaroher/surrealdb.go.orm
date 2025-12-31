package qb

// DefineAnalyzerStatement builds DEFINE ANALYZER.
type DefineAnalyzerStatement struct {
	Name       Node
	Tokenizers []Node
	Filters    []Node
	Comment    Node
}

func DefineAnalyzer(name string) *DefineAnalyzerStatement {
	return &DefineAnalyzerStatement{Name: Ident{Name: name}}
}

func DefineAnalyzerExpr(expr Node) *DefineAnalyzerStatement {
	return &DefineAnalyzerStatement{Name: expr}
}

func (d *DefineAnalyzerStatement) TokenizersList(tokenizers ...Node) *DefineAnalyzerStatement {
	d.Tokenizers = append(d.Tokenizers, tokenizers...)
	return d
}

func (d *DefineAnalyzerStatement) FiltersList(filters ...Node) *DefineAnalyzerStatement {
	d.Filters = append(d.Filters, filters...)
	return d
}

func (d *DefineAnalyzerStatement) CommentExpr(expr Node) *DefineAnalyzerStatement {
	d.Comment = expr
	return d
}

func (d *DefineAnalyzerStatement) CommentValue(value any) *DefineAnalyzerStatement {
	d.Comment = ensureValueNode(value)
	return d
}

func (d *DefineAnalyzerStatement) build(b *Builder) {
	b.Write("DEFINE ANALYZER ")
	d.Name.build(b)
	if len(d.Tokenizers) > 0 {
		b.Write(" TOKENIZERS ")
		renderNodes(b, d.Tokenizers)
	}
	if len(d.Filters) > 0 {
		b.Write(" FILTERS ")
		renderNodes(b, d.Filters)
	}
	if d.Comment != nil {
		b.Write(" COMMENT ")
		d.Comment.build(b)
	}
}

func (d *DefineAnalyzerStatement) Build() Query {
	return Build(d)
}

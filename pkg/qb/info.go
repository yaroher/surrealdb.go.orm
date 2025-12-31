package qb

// InfoStatement builds INFO statements.
type InfoStatement struct {
	Target string
	Name   Node
}

func InfoForRoot() *InfoStatement {
	return &InfoStatement{Target: "ROOT"}
}

func InfoForNamespace(ns string) *InfoStatement {
	return &InfoStatement{Target: "NS", Name: Ident{Name: ns}}
}

func InfoForNamespaceExpr(expr Node) *InfoStatement {
	return &InfoStatement{Target: "NS", Name: expr}
}

func InfoForDatabase(db string) *InfoStatement {
	return &InfoStatement{Target: "DB", Name: Ident{Name: db}}
}

func InfoForDatabaseExpr(expr Node) *InfoStatement {
	return &InfoStatement{Target: "DB", Name: expr}
}

func InfoForTable(table Node) *InfoStatement {
	return &InfoStatement{Target: "TABLE", Name: table}
}

func (i *InfoStatement) build(b *Builder) {
	b.Write("INFO FOR ")
	b.Write(i.Target)
	if i.Name != nil {
		b.Write(" ")
		i.Name.build(b)
	}
}

func (i *InfoStatement) Build() Query {
	return Build(i)
}

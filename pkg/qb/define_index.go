package qb

// DefineIndexStatement builds DEFINE INDEX.
type DefineIndexStatement struct {
	Name    Node
	Table   Node
	fields  []Node
	columns []Node
	Unique  bool
	Search  *SearchAnalyzer
}

func DefineIndex(name string) *DefineIndexStatement {
	return &DefineIndexStatement{Name: Ident{Name: name}}
}

func DefineIndexExpr(expr Node) *DefineIndexStatement {
	return &DefineIndexStatement{Name: expr}
}

func (d *DefineIndexStatement) OnTable(table Node) *DefineIndexStatement {
	d.Table = table
	return d
}

func (d *DefineIndexStatement) OnTableName(name string) *DefineIndexStatement {
	d.Table = Ident{Name: name}
	return d
}

func (d *DefineIndexStatement) Fields(fields ...Node) *DefineIndexStatement {
	d.fields = append(d.fields, fields...)
	return d
}

func (d *DefineIndexStatement) Columns(cols ...Node) *DefineIndexStatement {
	d.columns = append(d.columns, cols...)
	return d
}

func (d *DefineIndexStatement) UniqueOnly() *DefineIndexStatement {
	d.Unique = true
	return d
}

func (d *DefineIndexStatement) SearchAnalyzer(search *SearchAnalyzer) *DefineIndexStatement {
	d.Search = search
	return d
}

func (d *DefineIndexStatement) build(b *Builder) {
	b.Write("DEFINE INDEX ")
	d.Name.build(b)
	if d.Table != nil {
		b.Write(" ON TABLE ")
		d.Table.build(b)
	}
	if len(d.fields) > 0 {
		b.Write(" FIELDS ")
		renderNodes(b, d.fields)
	}
	if len(d.columns) > 0 {
		b.Write(" COLUMNS ")
		renderNodes(b, d.columns)
	}
	if d.Unique {
		b.Write(" UNIQUE")
	} else if d.Search != nil {
		b.Write(" ")
		d.Search.build(b)
	}
}

func (d *DefineIndexStatement) Build() Query {
	return Build(d)
}

// SearchAnalyzer builds SEARCH ANALYZER clause.
type SearchAnalyzer struct {
	Analyzer        Node
	highlight       bool
	scoring         Node
	docIDsOrder     Node
	docLengthsOrder Node
	postingsOrder   Node
	termsOrder      Node
}

func SearchAnalyzerFor(name string) *SearchAnalyzer {
	return &SearchAnalyzer{Analyzer: Ident{Name: name}}
}

func SearchAnalyzerExpr(expr Node) *SearchAnalyzer {
	return &SearchAnalyzer{Analyzer: expr}
}

func (s *SearchAnalyzer) Highlight() *SearchAnalyzer {
	s.highlight = true
	return s
}

func (s *SearchAnalyzer) BM25(k1 any, b any) *SearchAnalyzer {
	s.scoring = scoringBM25{K1: ensureValueNode(k1), B: ensureValueNode(b)}
	return s
}

func (s *SearchAnalyzer) VS() *SearchAnalyzer {
	s.scoring = scoringVS{}
	return s
}

func (s *SearchAnalyzer) DocIDsOrder(value any) *SearchAnalyzer {
	s.docIDsOrder = ensureValueNode(value)
	return s
}

func (s *SearchAnalyzer) DocLengthsOrder(value any) *SearchAnalyzer {
	s.docLengthsOrder = ensureValueNode(value)
	return s
}

func (s *SearchAnalyzer) PostingsOrder(value any) *SearchAnalyzer {
	s.postingsOrder = ensureValueNode(value)
	return s
}

func (s *SearchAnalyzer) TermsOrder(value any) *SearchAnalyzer {
	s.termsOrder = ensureValueNode(value)
	return s
}

func (s *SearchAnalyzer) build(b *Builder) {
	b.Write("FULLTEXT ANALYZER ")
	if s.Analyzer != nil {
		s.Analyzer.build(b)
	}
	if s.highlight {
		b.Write(" HIGHLIGHTS")
	}
	if s.scoring != nil {
		b.Write(" ")
		s.scoring.build(b)
	}
	if s.docIDsOrder != nil {
		b.Write(" DOC_IDS_ORDER ")
		s.docIDsOrder.build(b)
	}
	if s.docLengthsOrder != nil {
		b.Write(" DOC_LENGTHS_ORDER ")
		s.docLengthsOrder.build(b)
	}
	if s.postingsOrder != nil {
		b.Write(" POSTINGS_ORDER ")
		s.postingsOrder.build(b)
	}
	if s.termsOrder != nil {
		b.Write(" TERMS_ORDER ")
		s.termsOrder.build(b)
	}
}

type scoringBM25 struct {
	K1 Node
	B  Node
}

func (s scoringBM25) build(b *Builder) {
	b.Write("BM25 ")
	s.K1.build(b)
	b.Write(" ")
	s.B.build(b)
}

type scoringVS struct{}

func (s scoringVS) build(b *Builder) {
	b.Write("VS")
}

package qb

// ForBuilder builds FOR statements.
type ForBuilder struct {
	params   []string
	flowType string
	iterable Node
	body     Node
}

// For starts a FOR statement with one or more params.
func For(params ...string) *ForBuilder {
	return &ForBuilder{params: params, flowType: "IN"}
}

// In sets the iterable expression.
func (f *ForBuilder) In(iterable any) *ForBuilder {
	f.iterable = ensureValueNode(iterable)
	f.flowType = "IN"
	return f
}

// Block sets the loop body.
func (f *ForBuilder) Block(body Node) *ForBuilder {
	f.body = body
	return f
}

func (f *ForBuilder) build(b *Builder) {
	b.Write("FOR ")
	for i, p := range f.params {
		if i > 0 {
			b.Write(", ")
		}
		b.Write("$")
		b.Write(trimParamName(p))
	}
	b.Write(" ")
	b.Write(f.flowType)
	b.Write(" ")
	if f.iterable != nil {
		f.iterable.build(b)
	}
	b.Write(" ")
	BlockOf(f.body).build(b)
}

func (f *ForBuilder) Build() Query {
	return Build(f)
}

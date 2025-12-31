package qb

// Query is a rendered SurrealQL statement with bindings.
type Query struct {
	Text string
	Args map[string]any
}

// RawQuery builds a query from a raw SurrealQL string.
func RawQuery(text string, args map[string]any) Query {
	return Query{Text: text, Args: args}
}

package qb

// Table is a SurrealDB table name.
type Table string

func T(name string) Table {
	return Table(name)
}

func (t Table) Name() string {
	return string(t)
}

func (t Table) build(b *Builder) {
	b.Write(string(t))
}

// Field represents a typed field in a table.
type Field[T any] struct {
	Name string
}

func F[T any](name string) Field[T] {
	return Field[T]{Name: name}
}

func (f Field[T]) build(b *Builder) {
	b.Write(f.Name)
}

func (f Field[T]) As(alias string) Expr[any] {
	return As(f, alias)
}

func (f Field[T]) Expr() Expr[T] {
	return Expr[T]{node: f}
}

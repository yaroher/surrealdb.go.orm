package orm

// LinkSelf is a self-referential relation.
type LinkSelf[T any] struct{}

// LinkOne is a one-to-one relation.
type LinkOne[T any] struct{}

// LinkMany is a one-to-many relation.
type LinkMany[T any] struct{}

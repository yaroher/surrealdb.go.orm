package orm

// ID represents a typed record identifier.
type ID[T any, K comparable] struct {
	Value K
}

func NewID[T any, K comparable](value K) ID[T, K] {
	return ID[T, K]{Value: value}
}

// SimpleID is a string-based identifier.
type SimpleID[T any] struct {
	Value string
}

func NewSimpleID[T any](value string) SimpleID[T] {
	return SimpleID[T]{Value: value}
}

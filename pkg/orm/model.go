package orm

import "github.com/yaroher/surrealdb.go.orm/pkg/qb"

// Model is the minimal surface required by the query builder and migrator.
type Model interface {
	Table() qb.Table
}

// Node/Edge/Object are marker interfaces for generated models.
type Node interface{ Model }
type Edge interface {
	Model
	EdgeIn() qb.Table
	EdgeOut() qb.Table
}
type Object interface{ Model }

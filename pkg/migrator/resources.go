package migrator

import "github.com/yaroher/surrealdb.go.orm/pkg/qb"

// Definition is a named SurrealQL statement.
type Definition struct {
	Name      string
	Scope     string
	Statement qb.Statement
}

// ResourceSet represents the codebase resources for migrations.
type ResourceSet struct {
	Tables       map[string]Definition
	Fields       map[string]map[string]Definition
	Indexes      map[string]map[string]Definition
	Events       map[string]map[string]Definition
	Access       map[string]Definition
	AccessGrants []Definition
}

func NewResourceSet() ResourceSet {
	return ResourceSet{
		Tables:  map[string]Definition{},
		Fields:  map[string]map[string]Definition{},
		Indexes: map[string]map[string]Definition{},
		Events:  map[string]map[string]Definition{},
		Access:  map[string]Definition{},
	}
}

func (r *ResourceSet) AddTable(name string, stmt qb.Statement) {
	if r.Tables == nil {
		r.Tables = map[string]Definition{}
	}
	r.Tables[name] = Definition{Name: name, Statement: stmt}
}

func (r *ResourceSet) AddField(table string, name string, stmt qb.Statement) {
	if r.Fields == nil {
		r.Fields = map[string]map[string]Definition{}
	}
	if r.Fields[table] == nil {
		r.Fields[table] = map[string]Definition{}
	}
	r.Fields[table][name] = Definition{Name: name, Statement: stmt}
}

func (r *ResourceSet) AddIndex(table string, name string, stmt qb.Statement) {
	if r.Indexes == nil {
		r.Indexes = map[string]map[string]Definition{}
	}
	if r.Indexes[table] == nil {
		r.Indexes[table] = map[string]Definition{}
	}
	r.Indexes[table][name] = Definition{Name: name, Statement: stmt}
}

func (r *ResourceSet) AddEvent(table string, name string, stmt qb.Statement) {
	if r.Events == nil {
		r.Events = map[string]map[string]Definition{}
	}
	if r.Events[table] == nil {
		r.Events[table] = map[string]Definition{}
	}
	r.Events[table][name] = Definition{Name: name, Statement: stmt}
}

func (r *ResourceSet) AddAccess(name string, scope string, stmt qb.Statement) {
	if r.Access == nil {
		r.Access = map[string]Definition{}
	}
	r.Access[name] = Definition{Name: name, Scope: scope, Statement: stmt}
}

func (r *ResourceSet) AddAccessGrant(name string, stmt qb.Statement) {
	r.AccessGrants = append(r.AccessGrants, Definition{Name: name, Statement: stmt})
}

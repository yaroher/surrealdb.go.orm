package qb

import "strings"

// CrudType defines CRUD operations for permissions.
type CrudType string

const (
	CrudSelect CrudType = "select"
	CrudCreate CrudType = "create"
	CrudUpdate CrudType = "update"
	CrudDelete CrudType = "delete"
)

// ForPermissionClause defines a permission clause.
type ForPermissionClause struct {
	Crud []CrudType
	Cond Condition
}

// ForPermission starts a FOR permission clause.
func ForPermission(crud ...CrudType) *ForPermissionClause {
	return &ForPermissionClause{Crud: crud}
}

// Where sets the condition for the FOR clause.
func (f *ForPermissionClause) Where(cond Condition) *ForPermissionClause {
	f.Cond = cond
	return f
}

func (f *ForPermissionClause) build(b *Builder) {
	b.Write("FOR ")
	if len(f.Crud) == 0 {
		b.Write(string(CrudSelect))
	} else {
		parts := make([]string, 0, len(f.Crud))
		for _, ct := range f.Crud {
			parts = append(parts, string(ct))
		}
		b.Write(strings.Join(parts, ", "))
	}
	if f.Cond.node != nil {
		b.Write(" WHERE ")
		f.Cond.build(b)
	}
}

// Permissions groups permissions for DEFINE statements.
type Permissions struct {
	None  bool
	Full  bool
	Lines []Node
}

func (p *Permissions) With(lines ...Node) *Permissions {
	p.Lines = append(p.Lines, lines...)
	return p
}

func (p *Permissions) NoneOnly() *Permissions {
	p.None = true
	p.Full = false
	p.Lines = nil
	return p
}

func (p *Permissions) FullOnly() *Permissions {
	p.Full = true
	p.None = false
	p.Lines = nil
	return p
}

func renderPermissions(b *Builder, p *Permissions) {
	if p == nil {
		return
	}
	if p.None {
		b.Write(" PERMISSIONS NONE")
		return
	}
	if p.Full {
		b.Write(" PERMISSIONS FULL")
		return
	}
	if len(p.Lines) == 0 {
		return
	}
	b.Write(" PERMISSIONS\n")
	for i, line := range p.Lines {
		if i > 0 {
			b.Write("\n")
		}
		line.build(b)
	}
}

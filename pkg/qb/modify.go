package qb

// Assignment represents a field assignment.
type Assignment struct {
	Field Node
	Op    string
	Value Node
}

func Set(field Node, value any) Assignment {
	return Assignment{Field: field, Op: "=", Value: ensureValueNode(value)}
}

// CreateBuilder builds CREATE statements.
type CreateBuilder struct {
	target    Node
	content   Node
	set       []Assignment
	returning Node
}

func Create(target Node) *CreateBuilder {
	return &CreateBuilder{target: target}
}

func (c *CreateBuilder) Content(value any) *CreateBuilder {
	c.content = ensureValueNode(value)
	return c
}

func (c *CreateBuilder) Set(assignments ...Assignment) *CreateBuilder {
	c.set = assignments
	return c
}

func (c *CreateBuilder) Return(ret Node) *CreateBuilder {
	c.returning = ret
	return c
}

func (c *CreateBuilder) Build() Query {
	return Build(c)
}

func (c *CreateBuilder) build(b *Builder) {
	b.Write("CREATE ")
	c.target.build(b)
	if c.content != nil {
		b.Write(" CONTENT ")
		c.content.build(b)
	}
	if len(c.set) > 0 {
		b.Write(" SET ")
		renderAssignments(b, c.set)
	}
	if c.returning != nil {
		b.Write(" RETURN ")
		c.returning.build(b)
	}
}

// InsertBuilder builds INSERT statements.
type InsertBuilder struct {
	into      Node
	values    []Node
	returning Node
}

func Insert(into Node) *InsertBuilder {
	return &InsertBuilder{into: into}
}

func (i *InsertBuilder) Values(values ...any) *InsertBuilder {
	for _, v := range values {
		i.values = append(i.values, ensureValueNode(v))
	}
	return i
}

func (i *InsertBuilder) Return(ret Node) *InsertBuilder {
	i.returning = ret
	return i
}

func (i *InsertBuilder) Build() Query {
	return Build(i)
}

func (i *InsertBuilder) build(b *Builder) {
	b.Write("INSERT INTO ")
	i.into.build(b)
	if len(i.values) > 0 {
		b.Write(" ")
		renderNodes(b, i.values)
	}
	if i.returning != nil {
		b.Write(" RETURN ")
		i.returning.build(b)
	}
}

// UpdateBuilder builds UPDATE statements.
type UpdateBuilder struct {
	target    Node
	set       []Assignment
	where     Condition
	returning Node
}

func Update(target Node) *UpdateBuilder {
	return &UpdateBuilder{target: target}
}

func (u *UpdateBuilder) Set(assignments ...Assignment) *UpdateBuilder {
	u.set = assignments
	return u
}

func (u *UpdateBuilder) Where(cond Condition) *UpdateBuilder {
	u.where = cond
	return u
}

func (u *UpdateBuilder) Return(ret Node) *UpdateBuilder {
	u.returning = ret
	return u
}

func (u *UpdateBuilder) Build() Query {
	return Build(u)
}

func (u *UpdateBuilder) build(b *Builder) {
	b.Write("UPDATE ")
	u.target.build(b)
	if len(u.set) > 0 {
		b.Write(" SET ")
		renderAssignments(b, u.set)
	}
	if u.where.node != nil {
		b.Write(" WHERE ")
		u.where.build(b)
	}
	if u.returning != nil {
		b.Write(" RETURN ")
		u.returning.build(b)
	}
}

// DeleteBuilder builds DELETE statements.
type DeleteBuilder struct {
	target    Node
	where     Condition
	returning Node
}

func Delete(target Node) *DeleteBuilder {
	return &DeleteBuilder{target: target}
}

func (d *DeleteBuilder) Where(cond Condition) *DeleteBuilder {
	d.where = cond
	return d
}

func (d *DeleteBuilder) Return(ret Node) *DeleteBuilder {
	d.returning = ret
	return d
}

func (d *DeleteBuilder) Build() Query {
	return Build(d)
}

func (d *DeleteBuilder) build(b *Builder) {
	b.Write("DELETE ")
	d.target.build(b)
	if d.where.node != nil {
		b.Write(" WHERE ")
		d.where.build(b)
	}
	if d.returning != nil {
		b.Write(" RETURN ")
		d.returning.build(b)
	}
}

func renderAssignments(b *Builder, assigns []Assignment) {
	for i, a := range assigns {
		if i > 0 {
			b.Write(", ")
		}
		a.Field.build(b)
		b.Write(" ")
		b.Write(a.Op)
		b.Write(" ")
		a.Value.build(b)
	}
}

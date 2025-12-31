package qb

// Projection represents a select projection.
type Projection = Node

// All selects all fields.
type allProjection struct{}

func (a allProjection) build(b *Builder) {
	b.Write("*")
}

// All is a projection for selecting all fields.
var All Projection = allProjection{}

// SelectBuilder builds a SELECT statement.
type SelectBuilder struct {
	projections []Projection
	from        []Node
	where       Condition
	orders      []Order
	groupBy     []Node
	splitOn     []Node
	limit       Node
	start       Node
	fetch       []Node
	timeout     Node
	parallel    bool
}

// Select starts a SELECT statement.
func Select(projections ...Projection) *SelectBuilder {
	if len(projections) == 0 {
		projections = []Projection{All}
	}
	return &SelectBuilder{projections: projections}
}

// From sets the FROM targets.
func (s *SelectBuilder) From(targets ...Node) *SelectBuilder {
	s.from = targets
	return s
}

// Where sets the WHERE condition.
func (s *SelectBuilder) Where(cond Condition) *SelectBuilder {
	s.where = cond
	return s
}

// OrderBy sets the ORDER BY clause.
func (s *SelectBuilder) OrderBy(orders ...Order) *SelectBuilder {
	s.orders = orders
	return s
}

// GroupBy sets the GROUP BY clause.
func (s *SelectBuilder) GroupBy(fields ...Node) *SelectBuilder {
	s.groupBy = fields
	return s
}

// Split sets the SPLIT clause.
func (s *SelectBuilder) Split(fields ...Node) *SelectBuilder {
	s.splitOn = fields
	return s
}

// Limit sets LIMIT.
func (s *SelectBuilder) Limit(limit any) *SelectBuilder {
	s.limit = ensureValueNode(limit)
	return s
}

// Start sets START.
func (s *SelectBuilder) Start(start any) *SelectBuilder {
	s.start = ensureValueNode(start)
	return s
}

// Fetch sets FETCH fields.
func (s *SelectBuilder) Fetch(fields ...Node) *SelectBuilder {
	s.fetch = fields
	return s
}

// Timeout sets TIMEOUT duration or expression.
func (s *SelectBuilder) Timeout(value any) *SelectBuilder {
	s.timeout = ensureValueNode(value)
	return s
}

// Parallel enables PARALLEL.
func (s *SelectBuilder) Parallel() *SelectBuilder {
	s.parallel = true
	return s
}

// Build renders the query.
func (s *SelectBuilder) Build() Query {
	return Build(s)
}

func (s *SelectBuilder) build(b *Builder) {
	b.Write("SELECT ")
	renderNodes(b, s.projections)
	b.Write(" FROM ")
	renderNodes(b, s.from)

	if s.where.node != nil {
		b.Write(" WHERE ")
		s.where.build(b)
	}

	if len(s.splitOn) > 0 {
		b.Write(" SPLIT ")
		renderNodes(b, s.splitOn)
	}

	if len(s.groupBy) > 0 {
		b.Write(" GROUP BY ")
		renderNodes(b, s.groupBy)
	}

	if len(s.orders) > 0 {
		b.Write(" ORDER BY ")
		for i, order := range s.orders {
			if i > 0 {
				b.Write(", ")
			}
			order.build(b)
		}
	}

	if s.limit != nil {
		b.Write(" LIMIT ")
		s.limit.build(b)
	}

	if s.start != nil {
		b.Write(" START ")
		s.start.build(b)
	}

	if len(s.fetch) > 0 {
		b.Write(" FETCH ")
		renderNodes(b, s.fetch)
	}

	if s.timeout != nil {
		b.Write(" TIMEOUT ")
		s.timeout.build(b)
	}

	if s.parallel {
		b.Write(" PARALLEL")
	}
}

// Order defines ordering.
type Order struct {
	field  Node
	option string
	dir    string
}

func OrderBy(field Node) Order {
	return Order{field: field}
}

func (o Order) Rand() Order {
	o.option = "RAND()"
	return o
}

func (o Order) Collate() Order {
	o.option = "COLLATE"
	return o
}

func (o Order) Numeric() Order {
	o.option = "NUMERIC"
	return o
}

func (o Order) Asc() Order {
	o.dir = "ASC"
	return o
}

func (o Order) Desc() Order {
	o.dir = "DESC"
	return o
}

func (o Order) build(b *Builder) {
	o.field.build(b)
	if o.option != "" {
		b.Write(" ")
		b.Write(o.option)
	}
	if o.dir != "" {
		b.Write(" ")
		b.Write(o.dir)
	}
}

func renderNodes(b *Builder, nodes []Node) {
	for i, n := range nodes {
		if i > 0 {
			b.Write(", ")
		}
		n.build(b)
	}
}

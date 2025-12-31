package qb

// Eq builds "field = value".
func (f Field[T]) Eq(value any) Condition {
	return Expr[bool]{node: Binary{Left: f, Op: "=", Right: ensureValueNode(value)}}
}

func (f Field[T]) Neq(value any) Condition {
	return Expr[bool]{node: Binary{Left: f, Op: "!=", Right: ensureValueNode(value)}}
}

func (f Field[T]) Gt(value any) Condition {
	return Expr[bool]{node: Binary{Left: f, Op: ">", Right: ensureValueNode(value)}}
}

func (f Field[T]) Gte(value any) Condition {
	return Expr[bool]{node: Binary{Left: f, Op: ">=", Right: ensureValueNode(value)}}
}

func (f Field[T]) Lt(value any) Condition {
	return Expr[bool]{node: Binary{Left: f, Op: "<", Right: ensureValueNode(value)}}
}

func (f Field[T]) Lte(value any) Condition {
	return Expr[bool]{node: Binary{Left: f, Op: "<=", Right: ensureValueNode(value)}}
}

func (f Field[T]) Like(value any) Condition {
	return Expr[bool]{node: Binary{Left: f, Op: "~", Right: ensureValueNode(value)}}
}

func (f Field[T]) Contains(value any) Condition {
	return Expr[bool]{node: Binary{Left: f, Op: "CONTAINS", Right: ensureValueNode(value)}}
}

func (f Field[T]) In(values ...any) Condition {
	items := make([]Node, 0, len(values))
	for _, v := range values {
		items = append(items, ensureValueNode(v))
	}
	return Expr[bool]{node: Binary{Left: f, Op: "IN", Right: List{Items: items}}}
}

// And combines conditions with AND.
func And(conds ...Condition) Condition {
	return combine("AND", conds...)
}

// Or combines conditions with OR.
func Or(conds ...Condition) Condition {
	return combine("OR", conds...)
}

// Not negates a condition.
func Not(cond Condition) Condition {
	return Expr[bool]{node: Unary{Op: "!", Expr: cond}}
}

func combine(op string, conds ...Condition) Condition {
	if len(conds) == 0 {
		return Expr[bool]{node: RawExpr{Text: "true"}}
	}
	out := conds[0].node
	for i := 1; i < len(conds); i++ {
		out = Binary{Left: out, Op: op, Right: conds[i].node}
	}
	return Expr[bool]{node: out}
}

func ensureValueNode(v any) Node {
	if n, ok := v.(Node); ok {
		return n
	}
	return Value{Val: v}
}

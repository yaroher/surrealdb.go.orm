package qb

// Comparison helpers on any expression.
func (e Expr[T]) Eq(value any) Condition {
	return Expr[bool]{node: Binary{Left: e, Op: "=", Right: ensureValueNode(value)}}
}

func (e Expr[T]) EqExact(value any) Condition {
	return Expr[bool]{node: Binary{Left: e, Op: "==", Right: ensureValueNode(value)}}
}

func (e Expr[T]) AnyEq(value any) Condition {
	return Expr[bool]{node: Binary{Left: e, Op: "?=", Right: ensureValueNode(value)}}
}

func (e Expr[T]) AllEq(value any) Condition {
	return Expr[bool]{node: Binary{Left: e, Op: "*=", Right: ensureValueNode(value)}}
}

func (e Expr[T]) Neq(value any) Condition {
	return Expr[bool]{node: Binary{Left: e, Op: "!=", Right: ensureValueNode(value)}}
}

func (e Expr[T]) Lt(value any) Condition {
	return Expr[bool]{node: Binary{Left: e, Op: "<", Right: ensureValueNode(value)}}
}

func (e Expr[T]) Lte(value any) Condition {
	return Expr[bool]{node: Binary{Left: e, Op: "<=", Right: ensureValueNode(value)}}
}

func (e Expr[T]) Gt(value any) Condition {
	return Expr[bool]{node: Binary{Left: e, Op: ">", Right: ensureValueNode(value)}}
}

func (e Expr[T]) Gte(value any) Condition {
	return Expr[bool]{node: Binary{Left: e, Op: ">=", Right: ensureValueNode(value)}}
}

func (e Expr[T]) Like(value any) Condition {
	return Expr[bool]{node: Binary{Left: e, Op: "~", Right: ensureValueNode(value)}}
}

func (e Expr[T]) NotLike(value any) Condition {
	return Expr[bool]{node: Binary{Left: e, Op: "!~", Right: ensureValueNode(value)}}
}

func (e Expr[T]) AnyLike(value any) Condition {
	return Expr[bool]{node: Binary{Left: e, Op: "?~", Right: ensureValueNode(value)}}
}

func (e Expr[T]) AllLike(value any) Condition {
	return Expr[bool]{node: Binary{Left: e, Op: "*~", Right: ensureValueNode(value)}}
}

func (e Expr[T]) Is(value any) Condition {
	return Expr[bool]{node: Binary{Left: e, Op: "IS", Right: ensureValueNode(value)}}
}

func (e Expr[T]) IsNot(value any) Condition {
	return Expr[bool]{node: Binary{Left: e, Op: "IS NOT", Right: ensureValueNode(value)}}
}

func (e Expr[T]) Contains(value any) Condition {
	return Expr[bool]{node: Binary{Left: e, Op: "CONTAINS", Right: ensureValueNode(value)}}
}

func (e Expr[T]) In(values ...any) Condition {
	items := make([]Node, 0, len(values))
	for _, v := range values {
		items = append(items, ensureValueNode(v))
	}
	return Expr[bool]{node: Binary{Left: e, Op: "IN", Right: List{Items: items}}}
}

func (e Expr[T]) ContainsNot(value any) Condition {
	return Expr[bool]{node: Binary{Left: e, Op: "CONTAINSNOT", Right: ensureValueNode(value)}}
}

func (e Expr[T]) ContainsAll(value any) Condition {
	return Expr[bool]{node: Binary{Left: e, Op: "CONTAINSALL", Right: ensureValueNode(value)}}
}

func (e Expr[T]) ContainsAny(value any) Condition {
	return Expr[bool]{node: Binary{Left: e, Op: "CONTAINSANY", Right: ensureValueNode(value)}}
}

func (e Expr[T]) ContainsNone(value any) Condition {
	return Expr[bool]{node: Binary{Left: e, Op: "CONTAINSNONE", Right: ensureValueNode(value)}}
}

func (e Expr[T]) Inside(value any) Condition {
	return Expr[bool]{node: Binary{Left: e, Op: "INSIDE", Right: ensureValueNode(value)}}
}

func (e Expr[T]) NotInside(value any) Condition {
	return Expr[bool]{node: Binary{Left: e, Op: "NOTINSIDE", Right: ensureValueNode(value)}}
}

func (e Expr[T]) AllInside(value any) Condition {
	return Expr[bool]{node: Binary{Left: e, Op: "ALLINSIDE", Right: ensureValueNode(value)}}
}

func (e Expr[T]) AnyInside(value any) Condition {
	return Expr[bool]{node: Binary{Left: e, Op: "ANYINSIDE", Right: ensureValueNode(value)}}
}

func (e Expr[T]) NoneInside(value any) Condition {
	return Expr[bool]{node: Binary{Left: e, Op: "NONEINSIDE", Right: ensureValueNode(value)}}
}

func (e Expr[T]) Outside(value any) Condition {
	return Expr[bool]{node: Binary{Left: e, Op: "OUTSIDE", Right: ensureValueNode(value)}}
}

func (e Expr[T]) Intersects(value any) Condition {
	return Expr[bool]{node: Binary{Left: e, Op: "INTERSECTS", Right: ensureValueNode(value)}}
}

// Arithmetic helpers.
func (e Expr[T]) Add(value any) Expr[any] {
	return Expr[any]{node: Binary{Left: e, Op: "+", Right: ensureValueNode(value)}}
}

func (e Expr[T]) Sub(value any) Expr[any] {
	return Expr[any]{node: Binary{Left: e, Op: "-", Right: ensureValueNode(value)}}
}

func (e Expr[T]) Mul(value any) Expr[any] {
	return Expr[any]{node: Binary{Left: e, Op: "*", Right: ensureValueNode(value)}}
}

func (e Expr[T]) Div(value any) Expr[any] {
	return Expr[any]{node: Binary{Left: e, Op: "/", Right: ensureValueNode(value)}}
}

package qb

import "testing"

func TestExprBasics(t *testing.T) {
	assertQuery(t, Return(Raw("1 + 1")), "RETURN 1 + 1")

	q := Build(Return(PWith("id", 10)))
	if q.Text != "RETURN $id" {
		t.Fatalf("unexpected param text: %s", q.Text)
	}
	if q.Args["id"] != 10 {
		t.Fatalf("expected bound param id")
	}

	assertQuery(t, Return(As(I("x"), "alias")), "RETURN x AS alias")
	assertQuery(t, Return(Binary{Left: I("a"), Op: "+", Right: I("b")}), "RETURN (a + b)")
	assertQuery(t, Return(Unary{Op: "!", Expr: I("flag")}), "RETURN !flag")
	assertQuery(t, Return(Fn("fn", I("a"), V(2))), "RETURN fn(a, $p1)")
	assertQuery(t, Return(L(I("a"), I("b"))), "RETURN [a, b]")

	if got := trimParamName("$x"); got != "x" {
		t.Fatalf("expected trim of $x, got %s", got)
	}
	if got := trimParamName("x"); got != "x" {
		t.Fatalf("expected trim of x, got %s", got)
	}

	_ = debugNode(I("dbg"))
}

func TestFieldAndExprOps(t *testing.T) {
	field := F[int]("age")
	conds := []Node{
		field.Eq(1),
		field.Neq(2),
		field.Gt(3),
		field.Gte(4),
		field.Lt(5),
		field.Lte(6),
		field.Like("a"),
		field.Contains("x"),
		field.In(1, 2),
	}
	for _, cond := range conds {
		q := Build(Return(cond))
		if q.Text == "" {
			t.Fatalf("expected condition to build")
		}
	}

	expr := I("score")
	opNodes := []Node{
		expr.Eq(1),
		expr.EqExact(1),
		expr.AnyEq(1),
		expr.AllEq(1),
		expr.Neq(1),
		expr.Lt(1),
		expr.Lte(1),
		expr.Gt(1),
		expr.Gte(1),
		expr.Like("a"),
		expr.NotLike("a"),
		expr.AnyLike("a"),
		expr.AllLike("a"),
		expr.Is("NONE"),
		expr.IsNot("NONE"),
		expr.Contains("a"),
		expr.In(1, 2),
		expr.ContainsNot("a"),
		expr.ContainsAll("a"),
		expr.ContainsAny("a"),
		expr.ContainsNone("a"),
		expr.Inside("a"),
		expr.NotInside("a"),
		expr.AllInside("a"),
		expr.AnyInside("a"),
		expr.NoneInside("a"),
		expr.Outside("a"),
		expr.Intersects("a"),
		expr.Add(1),
		expr.Sub(1),
		expr.Mul(1),
		expr.Div(1),
	}
	for _, node := range opNodes {
		q := Build(Return(node))
		if q.Text == "" {
			t.Fatalf("expected expr to build")
		}
	}

	if got := Build(Return(And())).Text; got != "RETURN true" {
		t.Fatalf("expected empty And to be true, got %s", got)
	}
	_ = Build(Return(And(expr.Eq(1), expr.Eq(2))))
	_ = Build(Return(Or(expr.Eq(1), expr.Eq(2))))
	_ = Build(Return(Not(expr.Eq(1))))
}

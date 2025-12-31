package qbfn

import (
	"testing"

	"github.com/yaroher/surrealdb.go.orm/pkg/qb"
)

func TestModuleCall(t *testing.T) {
	expr := Array.Call2("len", qb.I("items"), 1)
	q := qb.Build(qb.Return(expr))
	if q.Text != "RETURN array::len(items, $p1)" {
		t.Fatalf("unexpected module call: %s", q.Text)
	}

	q = qb.Build(qb.Return(Array.Call1("first", qb.I("items"))))
	if q.Text != "RETURN array::first(items)" {
		t.Fatalf("unexpected module call1: %s", q.Text)
	}

	q = qb.Build(qb.Return(Array.Call3("slice", qb.I("items"), 1, 2)))
	if q.Text != "RETURN array::slice(items, $p1, $p2)" {
		t.Fatalf("unexpected module call3: %s", q.Text)
	}

	q = qb.Build(qb.Return(Fn("upper", "x")))
	if q.Text != "RETURN upper($p1)" {
		t.Fatalf("unexpected fn: %s", q.Text)
	}

	q = qb.Build(qb.Return(Count()))
	if q.Text != "RETURN count()" {
		t.Fatalf("unexpected count: %s", q.Text)
	}

	q = qb.Build(qb.Return(Sleep("1s")))
	if q.Text != "RETURN sleep($p1)" {
		t.Fatalf("unexpected sleep: %s", q.Text)
	}

	q = qb.Build(qb.Return(Rand()))
	if q.Text != "RETURN rand()" {
		t.Fatalf("unexpected rand: %s", q.Text)
	}
}

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

package orm

import (
	"testing"

	"github.com/yaroher/surrealdb.go.orm/pkg/qb"
)

type testModel struct{}

func (testModel) Table() qb.Table { return qb.T("test") }

func TestIDHelpers(t *testing.T) {
	id := NewID[testModel, int](10)
	if id.Value != 10 {
		t.Fatalf("expected id value 10")
	}
	sid := NewSimpleID[testModel]("x")
	if sid.Value != "x" {
		t.Fatalf("expected simple id value x")
	}
}

func TestLinkTypes(t *testing.T) {
	_ = LinkSelf[testModel]{}
	_ = LinkOne[testModel]{}
	_ = LinkMany[testModel]{}
}

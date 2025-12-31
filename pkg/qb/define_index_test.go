package qb

import "testing"

func TestDefineIndexSearchAnalyzer(t *testing.T) {
	stmt := DefineIndex("idx").
		OnTableName("book").
		Fields(I("title")).
		SearchAnalyzer(
			SearchAnalyzerFor("simple").
				Highlight().
				BM25(1.2, 0.7).
				DocIDsOrder(1),
		)

	q := Build(stmt)
	expected := "DEFINE INDEX idx ON TABLE book FIELDS title FULLTEXT ANALYZER simple HIGHLIGHTS BM25 $p1 $p2 DOC_IDS_ORDER $p3"
	if q.Text != expected {
		t.Fatalf("unexpected query:\n%s\nexpected:\n%s", q.Text, expected)
	}
	if len(q.Args) != 3 {
		t.Fatalf("expected 3 bound params, got %d", len(q.Args))
	}
}

func TestDefineIndexUniqueOverridesSearch(t *testing.T) {
	stmt := DefineIndex("idx").
		OnTableName("book").
		Fields(I("title")).
		UniqueOnly().
		SearchAnalyzer(SearchAnalyzerFor("simple"))

	q := Build(stmt)
	expected := "DEFINE INDEX idx ON TABLE book FIELDS title UNIQUE"
	if q.Text != expected {
		t.Fatalf("unexpected query:\n%s\nexpected:\n%s", q.Text, expected)
	}
}

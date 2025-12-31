package qb

import "testing"

func TestRemoveStatements(t *testing.T) {
	cases := []struct {
		stmt     Statement
		expected string
	}{
		{RemoveNamespace("ns"), "REMOVE NAMESPACE ns"},
		{RemoveDatabaseExpr(I("db")), "REMOVE DATABASE db"},
		{RemoveTable("t"), "REMOVE TABLE t"},
		{RemoveField("f").OnTableName("t"), "REMOVE FIELD f ON TABLE t"},
		{RemoveIndex("idx").OnTableName("t"), "REMOVE INDEX idx ON TABLE t"},
		{RemoveEvent("ev").OnTableName("t"), "REMOVE EVENT ev ON TABLE t"},
		{RemoveFunction("fn::x"), "REMOVE FUNCTION fn::x"},
		{RemoveParam("$p"), "REMOVE PARAM $p"},
		{RemoveScopeExpr(I("sc")), "REMOVE SCOPE sc"},
		{RemoveToken("tok").OnScope("sc"), "REMOVE TOKEN tok ON SCOPE sc"},
		{RemoveUser("bob").OnDatabase(), "REMOVE USER bob ON DATABASE"},
		{RemoveAnalyzer("an"), "REMOVE ANALYZER an"},
		{RemoveLogin("lg").OnNamespace(), "REMOVE LOGIN lg ON NAMESPACE"},
		{RemoveModel("m").VersionValue("1"), "REMOVE MODEL ml::m<$p1>"},
	}

	for i, tc := range cases {
		q := Build(tc.stmt)
		if q.Text != tc.expected {
			t.Fatalf("case %d unexpected query:\n%s\nexpected:\n%s", i, q.Text, tc.expected)
		}
	}
}

package qb

import "testing"

func TestAccessGrantAdditional(t *testing.T) {
	stmt := AccessGrantExpr(I("acc")).OnNamespace().Subject(Raw("user")).Token(Raw("tok")).Duration(Raw("1h"))
	assertQuery(t, stmt, "ACCESS acc ON NAMESPACE GRANT FOR user TOKEN tok DURATION 1h")
}

func TestDefineAccessAdditional(t *testing.T) {
	jwt := DefineAccessExpr(I("acc")).OnRoot().TypeJWT().JWTURL("https://example.com")
	q := assertQuery(t, jwt, "DEFINE ACCESS acc ON ROOT TYPE JWT URL $p1")
	assertArgsLen(t, q, 1)

	rec := DefineAccess("rec").TypeRecord().Signup(Raw("signup")).Signin(Raw("signin")).RecordWithJWT()
	assertQuery(t, rec, "DEFINE ACCESS rec TYPE RECORD SIGNUP { signup } SIGNIN { signin } WITH JWT")

	recURL := DefineAccess("rec").TypeRecord().RecordJWTURL("https://jwks")
	q = assertQuery(t, recURL, "DEFINE ACCESS rec TYPE RECORD WITH JWT URL $p1")
	assertArgsLen(t, q, 1)

	bear := DefineAccess("bear").TypeBearer().BearerForRecord()
	assertQuery(t, bear, "DEFINE ACCESS bear TYPE BEARER FOR RECORD")
}

func TestDefineAnalyzerAndEventAdditional(t *testing.T) {
	an := DefineAnalyzerExpr(I("an")).TokenizersList(I("t")).FiltersList(I("f")).CommentExpr(Raw("c"))
	assertQuery(t, an, "DEFINE ANALYZER an TOKENIZERS t FILTERS f COMMENT c")

	ev := DefineEventExpr(I("ev")).OnTable(T("t")).When(RawCond("true")).ThenValue("x")
	q := assertQuery(t, ev, "DEFINE EVENT ev ON TABLE t WHEN true THEN $p1")
	assertArgsLen(t, q, 1)
}

func TestDefineFieldAdditional(t *testing.T) {
	field := DefineFieldName("f", "t").TypeExpr(I("string")).ValueExpr(Raw("val")).DefaultExpr(Raw("def")).PermissionsFull()
	assertQuery(t, field, "DEFINE FIELD f ON TABLE t TYPE string VALUE val DEFAULT def PERMISSIONS FULL")
}

func TestDefineFunctionAndSignatureAdditional(t *testing.T) {
	fn := DefineFunction("fn", "$a").BodyExpr(Raw("return $a"))
	assertQuery(t, fn, "DEFINE FUNCTION fn($a) { return $a }")

	if got := FuncSignature("fn", []string{"$a", "b"}); got != "fn($a, $b)" {
		t.Fatalf("unexpected signature: %s", got)
	}
}

func TestDefineIndexSearchAnalyzerAdditional(t *testing.T) {
	idx := DefineIndexExpr(I("idx")).OnTable(T("t")).Columns(I("a"), I("b")).SearchAnalyzer(
		SearchAnalyzerExpr(I("an")).VS().DocLengthsOrder(1).PostingsOrder(2).TermsOrder(3),
	)
	q := assertQuery(t, idx, "DEFINE INDEX idx ON TABLE t COLUMNS a, b FULLTEXT ANALYZER an VS DOC_LENGTHS_ORDER $p1 POSTINGS_ORDER $p2 TERMS_ORDER $p3")
	assertArgsLen(t, q, 3)
}

func TestDefineModelAdditional(t *testing.T) {
	m := DefineModelExpr(I("m")).VersionExpr(Raw("1")).CommentExpr(Raw("note")).PermissionsNone()
	assertQuery(t, m, "DEFINE MODEL ml::m<1>\nCOMMENT note PERMISSIONS NONE")

	m2 := DefineModelExpr(I("m2")).PermissionsFull()
	assertQuery(t, m2, "DEFINE MODEL ml::m2 PERMISSIONS FULL")
}

func TestDefineParamAndScopeAdditional(t *testing.T) {
	param := DefineParam("p", 1).ValueExpr(Raw("2"))
	assertQuery(t, param, "DEFINE PARAM $p VALUE 2")

	scope := DefineScopeExpr(I("sc")).SessionExpr(Raw("1h")).SignupValue("signup").SigninExpr(Raw("signin"))
	q := assertQuery(t, scope, "DEFINE SCOPE sc SESSION 1h SIGNUP { $p1 } SIGNIN { signin }")
	assertArgsLen(t, q, 1)
}

func TestDefineNamespaceDatabaseTableAdditional(t *testing.T) {
	assertQuery(t, DefineNamespaceExpr(I("ns")), "DEFINE NAMESPACE ns")
	assertQuery(t, DefineDatabaseExpr(I("db")), "DEFINE DATABASE db")
	assertQuery(t, DefineTableName("t").SchemaLess(), "DEFINE TABLE t SCHEMALESS")
}

func TestDefineTokenAndUserAdditional(t *testing.T) {
	tok := DefineTokenExpr(I("tok")).OnNamespace().TypeExpr(I("HS256")).ValueExpr(Raw("secret"))
	assertQuery(t, tok, "DEFINE TOKEN tok ON NAMESPACE TYPE HS256 VALUE secret")

	assertQuery(t, DefineToken("tok2").OnDatabase(), "DEFINE TOKEN tok2 ON DATABASE")
	assertQuery(t, DefineToken("tok3").OnScopeExpr(I("scope")), "DEFINE TOKEN tok3 ON SCOPE scope")

	assertQuery(t, DefineUserExpr(I("u")).OnNamespace(), "DEFINE USER u ON NAMESPACE")
	assertQuery(t, DefineUserExpr(I("u")).OnDatabase(), "DEFINE USER u ON DATABASE")
}

func TestMiscStatementsAdditional(t *testing.T) {
	assertQuery(t, Let("x", 1), "LET $x = $p1")
	assertQuery(t, Return("x"), "RETURN $p1")
	assertQuery(t, Sleep("1s"), "SLEEP $p1")
	assertQuery(t, If(RawCond("true")).Then(Raw("ok")), "IF true { ok } END")
	assertQuery(t, For("i").In(1).Block(Raw("ok")), "FOR $i IN $p1 { ok }")
	assertQuery(t, Throw("err"), "THROW $p1")
	assertQuery(t, InfoForNamespaceExpr(I("ns")), "INFO FOR NS ns")
	assertQuery(t, InfoForDatabase("db"), "INFO FOR DB db")
	assertQuery(t, ShowChangesForTable(T("t")).SinceExpr(Raw("d")), "SHOW CHANGES FOR TABLE t SINCE d")
	assertQuery(t, Use().NamespaceExpr(I("ns")).DatabaseExpr(I("db")), "USE NS ns DB db")
}

func TestModifyBuildersAdditional(t *testing.T) {
	assertQuery(t, Create(T("t")).Set(Set(I("a"), 1)), "CREATE t SET a = $p1")
	assertQuery(t, Insert(T("t")).Values(1), "INSERT INTO t $p1")
	assertQuery(t, Update(T("t")), "UPDATE t")
	assertQuery(t, Delete(T("t")), "DELETE t")
	assertQuery(t, Relate(I("a"), I("e"), I("b")), "RELATE a -> e -> b")
}

func TestRawQueryAdditional(t *testing.T) {
	q := RawQuery("SELECT 1", nil)
	if q.Text != "SELECT 1" {
		t.Fatalf("unexpected raw text: %s", q.Text)
	}
	if q.Args != nil {
		t.Fatalf("expected no args")
	}
}

func TestSelectAndChainAdditional(t *testing.T) {
	assertQuery(t, Select(I("a")).From(T("t")).OrderBy(OrderBy(I("a")).Numeric().Desc()), "SELECT a FROM t ORDER BY a NUMERIC DESC")
	assertQuery(t, QueryChain(DefineNamespace("ns"), DefineDatabase("db")), "DEFINE NAMESPACE ns; DEFINE DATABASE db")
}

func TestTypesHelpersAdditional(t *testing.T) {
	if Table("t").Name() != "t" {
		t.Fatalf("expected table name")
	}
	if got := Build(Return(F[int]("age").As("a"))).Text; got != "RETURN age AS a" {
		t.Fatalf("unexpected field alias: %s", got)
	}
}

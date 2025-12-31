package qb

import "testing"

func TestCoverageExtras(t *testing.T) {
	// Access grant variants
	_ = AccessGrantExpr(I("acc")).OnNamespace().Subject(Raw("user")).Token(Raw("tok")).Duration(Raw("1h")).Build()

	// Access definitions
	_ = DefineAccessExpr(I("acc")).OnRoot().TypeJWT().JWTURL("https://example.com").Build()
	_ = DefineAccess("rec").TypeRecord().Signup(Raw("signup"))
	_ = DefineAccess("rec").TypeRecord().Signin(Raw("signin"))
	_ = DefineAccess("rec").TypeRecord().RecordWithJWT().RecordJWTAlgorithmKey("HS256", "key").RecordIssuerKey("iss").WithRefresh().Build()
	_ = DefineAccess("rec").TypeRecord().RecordJWTURL("https://jwks").Build()
	_ = DefineAccess("bear").TypeBearer().BearerForRecord().Build()

	// Analyzer and events
	_ = DefineAnalyzerExpr(I("an")).TokenizersList(I("t")).FiltersList(I("f")).CommentExpr(Raw("c")).Build()
	_ = DefineEventExpr(I("ev")).OnTable(T("t")).When(RawCond("true")).ThenValue("x").Build()

	// Field helpers
	_ = DefineFieldName("f", "t").TypeExpr(I("string")).ValueExpr(Raw("val")).DefaultExpr(Raw("def")).PermissionsFull().PermissionsFor(ForPermission()).Build()

	// Function signature
	_ = DefineFunction("fn", "$a").BodyExpr(Raw("return $a")).Build()
	_ = FuncSignature("fn", []string{"$a", "b"})

	// Index helpers
	_ = DefineIndexExpr(I("idx")).OnTable(T("t")).Columns(I("a"), I("b")).SearchAnalyzer(
		SearchAnalyzerExpr(I("an")).VS().DocLengthsOrder(1).PostingsOrder(2).TermsOrder(3),
	).Build()

	// Model helpers
	_ = DefineModelExpr(I("m")).VersionExpr(Raw("1")).CommentExpr(Raw("note")).PermissionsNone().Build()
	_ = DefineModelExpr(I("m2")).PermissionsFull().Build()

	// Param helpers
	_ = DefineParam("p", 1).ValueExpr(Raw("2")).Build()

	// Scope helpers
	_ = DefineScopeExpr(I("sc")).SessionExpr(Raw("1h")).SignupValue("signup").SigninExpr(Raw("signin")).Build()

	// Namespace / DB
	_ = DefineNamespaceExpr(I("ns")).Build()
	_ = DefineDatabaseExpr(I("db")).Build()

	// Table helpers
	_ = DefineTableName("t").SchemaLess().Build()

	// Token helpers
	_ = DefineTokenExpr(I("tok")).OnNamespace().TypeExpr(I("HS256")).ValueExpr(Raw("secret")).Build()
	_ = DefineToken("tok2").OnDatabase().Build()
	_ = DefineToken("tok3").OnScopeExpr(I("scope")).Build()

	// User helpers
	_ = DefineUserExpr(I("u")).OnNamespace().Build()
	_ = DefineUserExpr(I("u")).OnDatabase().Build()

	// Param and Let build
	_ = Let("x", 1).Build()
	_ = Return("x").Build()
	_ = Sleep("1s").Build()
	_ = If(RawCond("true")).Then(Raw("ok")).Build()
	_ = For("i").In(1).Block(Raw("ok")).Build()
	_ = Throw("err").Build()
	_ = InfoForNamespaceExpr(I("ns")).Build()
	_ = InfoForDatabase("db").Build()
	_ = ShowChangesForTable(T("t")).SinceExpr(Raw("d")).Build()
	_ = Use().NamespaceExpr(I("ns")).DatabaseExpr(I("db")).Build()

	// Modify builders Build methods
	_ = Create(T("t")).Set(Set(I("a"), 1)).Build()
	_ = Insert(T("t")).Values(1).Build()
	_ = Update(T("t")).Build()
	_ = Delete(T("t")).Build()
	_ = Relate(I("a"), I("e"), I("b")).Build()

	// Raw query
	_ = RawQuery("SELECT 1", nil)

	// Remove variants
	_ = RemoveFunction("fn").Build()
	_ = RemoveNamespaceExpr(I("ns")).Build()
	_ = RemoveDatabase("db").Build()
	_ = RemoveTableExpr(I("t")).Build()
	_ = RemoveFieldExpr(I("f")).OnTable(T("t")).Build()
	_ = RemoveIndexExpr(I("i")).OnTable(T("t")).Build()
	_ = RemoveEventExpr(I("e")).OnTable(T("t")).Build()
	_ = RemoveParam("p").Build()
	_ = RemoveScope("s").Build()
	_ = RemoveTokenExpr(I("t")).OnNamespace().Build()
	_ = RemoveToken("t").OnDatabase().Build()
	_ = RemoveToken("t").OnScopeExpr(I("sc")).Build()
	_ = RemoveUserExpr(I("u")).OnRoot().Build()
	_ = RemoveUser("u").OnNamespace().Build()
	_ = RemoveAnalyzerExpr(I("a")).Build()
	_ = RemoveLoginExpr(I("l")).OnDatabase().Build()
	_ = RemoveModelExpr(I("m")).VersionExpr(Raw("1")).Build()
	_ = RemoveAccessExpr(I("a")).OnNamespace().Build()

	// Select build helpers
	_ = Select(I("a")).From(T("t")).OrderBy(OrderBy(I("a")).Numeric().Desc()).Build()

	// Chain build
	_ = QueryChain(DefineNamespace("ns"), DefineDatabase("db")).Build()

	// Types helpers
	if Table("t").Name() == "" {
		t.Fatalf("expected table name")
	}
	_ = F[int]("age").As("a")
	_ = F[int]("age").Expr()

	// Expr helpers
	_ = P("name")
	_ = ensureNode(I("x"))
	_ = ensureNode(Raw("x"))
	_ = ensureNode(1)
}

package qb

import "testing"

func TestBuilderBindInitializesArgs(t *testing.T) {
	var b Builder
	b.Bind("p1", 1)
	if b.Args() == nil {
		t.Fatalf("expected args to be initialized")
	}
	if got := b.Args()["p1"]; got != 1 {
		t.Fatalf("expected p1=1, got %v", got)
	}
}

func TestDefineAccessEnsurePaths(t *testing.T) {
	jwt := DefineAccess("acc")
	jwt.JWTAlgorithmKey("HS256", "key")
	jwt.JWTURL("https://example.com")
	assertQuery(t, jwt, "DEFINE ACCESS acc TYPE JWT URL $p1")

	rec := DefineAccess("rec").TypeJWT().RecordWithJWT().Signup(Raw("signup"))
	rec.Signin(Raw("signin"))
	assertQuery(t, rec, "DEFINE ACCESS rec TYPE RECORD SIGNUP { signup } SIGNIN { signin } WITH JWT")

	bearer := DefineAccess("bear").BearerForUser().BearerForRecord()
	assertQuery(t, bearer, "DEFINE ACCESS bear TYPE BEARER FOR RECORD")
}

func TestDefineAccessRecordJWTAlgorithmKeyWithoutInit(t *testing.T) {
	rec := DefineAccess("rec").TypeRecord().RecordJWTAlgorithmKey("HS256", "key")
	assertQuery(t, rec, "DEFINE ACCESS rec TYPE RECORD WITH JWT ALGORITHM HS256 KEY $p1")
}

func TestDefineFieldPermissionsForAndEmpty(t *testing.T) {
	stmt := DefineFieldName("f", "t").PermissionsFor(ForPermission(CrudSelect))
	assertQuery(t, stmt, "DEFINE FIELD f ON TABLE t PERMISSIONS\nFOR select")

	empty := DefineFieldName("f", "t").PermissionsFor()
	assertQuery(t, empty, "DEFINE FIELD f ON TABLE t")
}

func TestDefineFunctionMultipleParamsAndTrim(t *testing.T) {
	stmt := DefineFunction("fn", "$a", "b").BodyExpr(Raw("return $a"))
	assertQuery(t, stmt, "DEFINE FUNCTION fn($a, $b) { return $a }")

	if got := FuncSignature("fn", []string{""}); got != "fn($)" {
		t.Fatalf("unexpected signature: %s", got)
	}
}

func TestEnsureNodeExprBranch(t *testing.T) {
	n := ensureNode(V(1))
	b := NewBuilder()
	n.build(b)
	if got := b.String(); got != "$p1" {
		t.Fatalf("expected $p1, got %q", got)
	}
}

func TestEnsureNodeNodeBranch(t *testing.T) {
	n := ensureNode(Ident{Name: "x"})
	b := NewBuilder()
	n.build(b)
	if got := b.String(); got != "x" {
		t.Fatalf("expected x, got %q", got)
	}
}

func TestIfThenNoBranch(t *testing.T) {
	ifb := &IfBuilder{}
	ifb.Then(Raw("noop"))
	if len(ifb.branches) != 0 {
		t.Fatalf("expected no branches to be added")
	}
}

func TestRenderAssignmentsMultiple(t *testing.T) {
	stmt := Update(T("t")).Set(Set(I("a"), 1), Set(I("b"), 2))
	assertQuery(t, stmt, "UPDATE t SET a = $p1, b = $p2")
}

func TestRenderPermissionsMultipleLines(t *testing.T) {
	stmt := DefineTableName("t").PermissionsFor(
		ForPermission(CrudSelect),
		ForPermission(CrudCreate),
	)
	assertQuery(t, stmt, "DEFINE TABLE t PERMISSIONS\nFOR select\nFOR create")
}

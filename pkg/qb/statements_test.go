package qb

import "testing"

func TestSelectBuilder(t *testing.T) {
	stmt := Select().From(T("user")).
		Where(F[int]("age").Gt(18)).
		Split(I("tags")).
		GroupBy(I("role")).
		OrderBy(OrderBy(I("name")).Collate().Asc(), OrderBy(I("id")).Rand()).
		Limit(10).
		Start(5).
		Fetch(I("profile")).
		Timeout("5s").
		Parallel()

	expected := "SELECT * FROM user WHERE (age > $p1) SPLIT tags GROUP BY role ORDER BY name COLLATE ASC, id RAND() LIMIT $p2 START $p3 FETCH profile TIMEOUT $p4 PARALLEL"
	q := assertQuery(t, stmt, expected)
	assertArgsLen(t, q, 4)
}

func TestModifyBuilders(t *testing.T) {
	create := Create(T("user")).Content(map[string]any{"name": "a"}).Return(I("id"))
	q := Build(create)
	if q.Text != "CREATE user CONTENT $p1 RETURN id" {
		t.Fatalf("unexpected create: %s", q.Text)
	}

	insert := Insert(T("user")).Values(map[string]any{"name": "a"}).Return(I("id"))
	q = Build(insert)
	if q.Text != "INSERT INTO user $p1 RETURN id" {
		t.Fatalf("unexpected insert: %s", q.Text)
	}

	update := Update(T("user")).Set(Set(I("name"), "b")).Where(I("id").Eq("u:1")).Return(I("id"))
	q = Build(update)
	if q.Text != "UPDATE user SET name = $p1 WHERE (id = $p2) RETURN id" {
		t.Fatalf("unexpected update: %s", q.Text)
	}

	del := Delete(T("user")).Where(I("id").Eq("u:1")).Return(I("id"))
	q = Build(del)
	if q.Text != "DELETE user WHERE (id = $p1) RETURN id" {
		t.Fatalf("unexpected delete: %s", q.Text)
	}

	relate := Relate(I("user:1"), I("likes"), I("post:1")).Set(Set(I("created_at"), "now")).Return(I("id"))
	q = Build(relate)
	if q.Text != "RELATE user:1 -> likes -> post:1 SET created_at = $p1 RETURN id" {
		t.Fatalf("unexpected relate: %s", q.Text)
	}
}

func TestUseLetReturnSleepShowInfo(t *testing.T) {
	assertQuery(t, Use().Namespace("test").Database("db"), "USE NS test DB db")
	assertQuery(t, Let("$x", 1), "LET $x = $p1")
	assertQuery(t, Return(I("id")), "RETURN id")
	assertQuery(t, Sleep("1s"), "SLEEP $p1")
	assertQuery(t, ShowChangesForTable(T("user")).Since("2020-01-01").Limit(10), "SHOW CHANGES FOR TABLE user SINCE $p1 LIMIT $p2")
	assertQuery(t, InfoForRoot(), "INFO FOR ROOT")
	assertQuery(t, InfoForNamespace("ns"), "INFO FOR NS ns")
	assertQuery(t, InfoForDatabaseExpr(I("db")), "INFO FOR DB db")
	assertQuery(t, InfoForTable(T("user")), "INFO FOR TABLE user")
}

func TestIfForFlowAndTransaction(t *testing.T) {
	ifStmt := If(I("a").Eq(1)).Then(Raw("do()"))
	ifStmt.ElseIf(I("b").Eq(2)).Then(Raw("do2()"))
	ifStmt.Else(Raw("do3()"))
	assertQuery(t, ifStmt, "IF (a = $p1) { do() } ELSE IF (b = $p2) { do2() } ELSE { do3() } END")

	forStmt := For("$item", "idx").In(L(V(1), V(2))).Block(Raw("do()"))
	assertQuery(t, forStmt, "FOR $item, $idx IN [$p1, $p2] { do() }")

	assertQuery(t, Break(), "BREAK")
	assertQuery(t, Continue(), "CONTINUE")
	assertQuery(t, Throw("err"), "THROW $p1")

	assertQuery(t, BeginTransaction(), "BEGIN TRANSACTION")
	assertQuery(t, CommitTransaction(), "COMMIT TRANSACTION")
	assertQuery(t, CancelTransaction(), "CANCEL TRANSACTION")
}

func TestDefineStatements(t *testing.T) {
	assertQuery(t, DefineNamespace("test"), "DEFINE NAMESPACE test")
	assertQuery(t, DefineDatabaseExpr(I("db")), "DEFINE DATABASE db")

	table := DefineTableName("user").DropTable().SchemaFull().AsSelect(Select(I("id")).From(T("user"))).PermissionsFull()
	assertQuery(t, table, "DEFINE TABLE user DROP SCHEMAFULL AS SELECT id FROM user PERMISSIONS FULL")

	field := DefineFieldName("name", "user").Type("string").Value("x").Assert(RawCond("$value != ''")).Default("y").PermissionsNone()
	assertQuery(t, field, "DEFINE FIELD name ON TABLE user TYPE string VALUE $p1 ASSERT $value != '' DEFAULT $p2 PERMISSIONS NONE")

	event := DefineEvent("ev").OnTableName("user").When(RawCond("$event == 'CREATE'")).ThenExpr(Raw("fn()"))
	assertQuery(t, event, "DEFINE EVENT ev ON TABLE user WHEN $event == 'CREATE' THEN fn()")

	scope := DefineScope("scope").SessionValue("1h").SignupExpr(Raw("signup()"))
	scope.SigninValue("signin()")
	assertQuery(t, scope, "DEFINE SCOPE scope SESSION $p1 SIGNUP { signup() } SIGNIN { $p2 }")

	token := DefineToken("tok").OnScope("scope").Type("HS256").Value("secret")
	assertQuery(t, token, "DEFINE TOKEN tok ON SCOPE scope TYPE HS256 VALUE $p1")

	user := DefineUser("alice").OnRoot().Password("pass").RolesList("admin", "reader")
	assertQuery(t, user, "DEFINE USER alice ON ROOT PASSWORD $p1 ROLES admin, reader")

	model := DefineModel("embed").VersionValue("1").CommentValue("note").PermissionsFor(ForPermission(CrudSelect).Where(RawCond("true")))
	assertQuery(t, model, "DEFINE MODEL ml::embed<$p1>\nCOMMENT $p2 PERMISSIONS\nFOR select WHERE true")

	fn := DefineFunction("util::echo", "$x").BodyExpr(Raw("return $x"))
	assertQuery(t, fn, "DEFINE FUNCTION util::echo($x) { return $x }")

	ana := DefineAnalyzer("an").TokenizersList(I("blank")).FiltersList(I("lowercase")).CommentValue("ok")
	assertQuery(t, ana, "DEFINE ANALYZER an TOKENIZERS blank FILTERS lowercase COMMENT $p1")

	param := DefineParam("$val", "x")
	assertQuery(t, param, "DEFINE PARAM $val VALUE $p1")
}

func TestDefineAccessBearer(t *testing.T) {
	stmt := DefineAccess("bear").OnDatabase().TypeBearer().BearerForUser().Authenticate(Raw("auth()"))
	q := Build(stmt)
	if q.Text != "DEFINE ACCESS bear ON DATABASE TYPE BEARER FOR USER AUTHENTICATE auth()" {
		t.Fatalf("unexpected bearer access: %s", q.Text)
	}
}

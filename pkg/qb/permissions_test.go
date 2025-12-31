package qb

import "testing"

func TestPermissions(t *testing.T) {
	perm := &Permissions{}
	perm.With(ForPermission(CrudSelect).Where(RawCond("true")))
	q := Build(DefineTableName("t").PermissionsFor(perm.Lines...))
	if q.Text != "DEFINE TABLE t PERMISSIONS\nFOR select WHERE true" {
		t.Fatalf("unexpected permissions: %s", q.Text)
	}

	perm.FullOnly()
	q = Build(DefineTableName("t").PermissionsFull())
	if q.Text != "DEFINE TABLE t PERMISSIONS FULL" {
		t.Fatalf("unexpected full permissions: %s", q.Text)
	}

	perm.NoneOnly()
	q = Build(DefineTableName("t").PermissionsNone())
	if q.Text != "DEFINE TABLE t PERMISSIONS NONE" {
		t.Fatalf("unexpected none permissions: %s", q.Text)
	}
}

func TestForPermissionDefaults(t *testing.T) {
	clause := ForPermission().Where(RawCond("true"))
	q := Build(Return(clause))
	if q.Text != "RETURN FOR select WHERE true" {
		t.Fatalf("unexpected for permission default: %s", q.Text)
	}
}

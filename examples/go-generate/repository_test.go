package gogenerate

import (
	"context"
	"reflect"
	"testing"

	"github.com/yaroher/surrealdb.go.orm/pkg/qb"
)

type mockDB struct {
	lastExec  QueryCapture
	lastQuery QueryCapture
	execErr   error
	queryErr  error
}

type QueryCapture struct {
	Query qb.Query
	Out   any
}

func (m *mockDB) Exec(ctx context.Context, query qb.Query) error {
	m.lastExec = QueryCapture{Query: query}
	return m.execErr
}

func (m *mockDB) Query(ctx context.Context, query qb.Query, out any) error {
	m.lastQuery = QueryCapture{Query: query, Out: out}
	return m.queryErr
}

func TestRepositoryCreateUser(t *testing.T) {
	db := &mockDB{}
	repo := NewRepository(db)

	user := User{ID: "user:1", FirstName: "Ana", LastName: "Fox", Email: "ana@example.com"}
	if err := repo.CreateUser(context.Background(), user); err != nil {
		t.Fatalf("create user: %v", err)
	}

	if db.lastExec.Query.Text != "CREATE users CONTENT $p1" {
		t.Fatalf("unexpected query: %s", db.lastExec.Query.Text)
	}

	payload, ok := db.lastExec.Query.Args["p1"].(map[string]any)
	if !ok {
		t.Fatalf("expected content payload")
	}

	want := map[string]any{
		"id":         user.ID,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"email":      user.Email,
	}
	if !reflect.DeepEqual(payload, want) {
		t.Fatalf("unexpected payload: %+v", payload)
	}
}

func TestRepositoryListUsers(t *testing.T) {
	db := &mockDB{}
	repo := NewRepository(db)

	_, err := repo.ListUsers(context.Background())
	if err != nil {
		t.Fatalf("list users: %v", err)
	}

	if db.lastQuery.Query.Text != "SELECT id, first_name, last_name, email FROM users" {
		t.Fatalf("unexpected query: %s", db.lastQuery.Query.Text)
	}
	if db.lastQuery.Out == nil {
		t.Fatalf("expected output destination")
	}
}

func TestRepositoryFindUserByEmail(t *testing.T) {
	db := &mockDB{}
	repo := NewRepository(db)

	_, err := repo.FindUserByEmail(context.Background(), "ana@example.com")
	if err != nil {
		t.Fatalf("find by email: %v", err)
	}

	if db.lastQuery.Query.Text != "SELECT id, first_name, last_name, email FROM users WHERE (email = $p1)" {
		t.Fatalf("unexpected query: %s", db.lastQuery.Query.Text)
	}
	if got := db.lastQuery.Query.Args["p1"]; got != "ana@example.com" {
		t.Fatalf("unexpected arg: %v", got)
	}
}

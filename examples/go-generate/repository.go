package gogenerate

import (
	"context"

	"github.com/yaroher/surrealdb.go.orm/pkg/qb"
)

type DB interface {
	Exec(ctx context.Context, query qb.Query) error
	Query(ctx context.Context, query qb.Query, out any) error
}

type Repository struct {
	db DB
}

func NewRepository(db DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateUser(ctx context.Context, user User) error {
	schema := (User{}).Schema()
	payload := map[string]any{
		schema.ID.Name:        user.ID,
		schema.FirstName.Name: user.FirstName,
		schema.LastName.Name:  user.LastName,
		schema.Email.Name:     user.Email,
	}
	q := qb.Create((User{}).Table()).Content(payload).Build()
	return r.db.Exec(ctx, q)
}

func (r *Repository) ListUsers(ctx context.Context) ([]User, error) {
	schema := (User{}).Schema()
	q := qb.Select(schema.ID, schema.FirstName, schema.LastName, schema.Email).
		From((User{}).Table()).
		Build()
	var out []User
	if err := r.db.Query(ctx, q, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *Repository) FindUserByEmail(ctx context.Context, email string) ([]User, error) {
	schema := (User{}).Schema()
	q := qb.Select(schema.ID, schema.FirstName, schema.LastName, schema.Email).
		From((User{}).Table()).
		Where(schema.Email.Eq(email)).
		Build()
	var out []User
	if err := r.db.Query(ctx, q, &out); err != nil {
		return nil, err
	}
	return out, nil
}

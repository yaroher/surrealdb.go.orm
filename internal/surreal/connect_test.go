package surreal

import (
	"context"
	"testing"

	surrealdb "github.com/surrealdb/surrealdb.go"
)

func TestConnectInvalidDSN(t *testing.T) {
	if _, err := Connect(context.Background(), "://", "", "", "", ""); err == nil {
		t.Fatalf("expected error for invalid dsn")
	}
}

func TestConnectWithAuthAndUse(t *testing.T) {
	origFrom := fromEndpointFn
	origSignIn := signInFn
	origUse := useFn
	defer func() {
		fromEndpointFn = origFrom
		signInFn = origSignIn
		useFn = origUse
	}()

	fromEndpointFn = func(ctx context.Context, dsn string) (*surrealdb.DB, error) {
		return &surrealdb.DB{}, nil
	}

	signInCalled := false
	useCalled := false
	signInFn = func(ctx context.Context, db *surrealdb.DB, auth surrealdb.Auth) (any, error) {
		signInCalled = true
		return nil, nil
	}
	useFn = func(ctx context.Context, db *surrealdb.DB, ns, dbName string) error {
		useCalled = true
		return nil
	}

	if _, err := Connect(context.Background(), "http://example", "ns", "db", "user", "pass"); err != nil {
		t.Fatalf("connect: %v", err)
	}
	if !signInCalled {
		t.Fatalf("expected signIn to be called")
	}
	if !useCalled {
		t.Fatalf("expected use to be called")
	}
}

func TestConnectAuthError(t *testing.T) {
	origFrom := fromEndpointFn
	origSignIn := signInFn
	defer func() {
		fromEndpointFn = origFrom
		signInFn = origSignIn
	}()

	fromEndpointFn = func(ctx context.Context, dsn string) (*surrealdb.DB, error) {
		return &surrealdb.DB{}, nil
	}
	signInFn = func(ctx context.Context, db *surrealdb.DB, auth surrealdb.Auth) (any, error) {
		return nil, context.Canceled
	}

	if _, err := Connect(context.Background(), "http://example", "ns", "db", "user", "pass"); err == nil {
		t.Fatalf("expected signIn error")
	}
}

func TestConnectUseError(t *testing.T) {
	origFrom := fromEndpointFn
	origUse := useFn
	defer func() {
		fromEndpointFn = origFrom
		useFn = origUse
	}()

	fromEndpointFn = func(ctx context.Context, dsn string) (*surrealdb.DB, error) {
		return &surrealdb.DB{}, nil
	}
	useFn = func(ctx context.Context, db *surrealdb.DB, ns, dbName string) error {
		return context.Canceled
	}

	if _, err := Connect(context.Background(), "http://example", "ns", "db", "", ""); err == nil {
		t.Fatalf("expected use error")
	}
}

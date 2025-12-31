package migrator

import (
	"context"
	"sync"
)

type fakeDB struct {
	mu    sync.Mutex
	fn    func(sql string) ([]map[string]any, error)
	calls []string
}

func (f *fakeDB) Query(ctx context.Context, sql string, vars map[string]any) ([]map[string]any, error) {
	f.mu.Lock()
	f.calls = append(f.calls, sql)
	f.mu.Unlock()
	if f.fn != nil {
		return f.fn(sql)
	}
	return nil, nil
}

func (f *fakeDB) Calls() []string {
	f.mu.Lock()
	defer f.mu.Unlock()
	out := make([]string, len(f.calls))
	copy(out, f.calls)
	return out
}

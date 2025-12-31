package migrator

import (
	"sort"
	"strings"

	"github.com/yaroher/surrealdb.go.orm/pkg/qb"
)

// DiffOptions control diff behavior.
type DiffOptions struct {
	Prompter       Prompter
	Force          bool
	RenameStrategy string
	RenameExpr     string
}

// DiffResources compares codebase and database resources and returns up/down statements.
func DiffResources(code ResourceSet, db ResourceSet, opts DiffOptions) (up []qb.Statement, down []qb.Statement) {
	up = []qb.Statement{}
	down = []qb.Statement{}

	// Tables
	allTables := unionKeys(code.Tables, db.Tables)
	for _, table := range allTables {
		codeDef, hasCode := code.Tables[table]
		dbDef, hasDB := db.Tables[table]
		switch {
		case hasCode && !hasDB:
			up = append(up, codeDef.Statement)
			down = append(down, qb.RemoveTable(table))
		case !hasCode && hasDB:
			up = append(up, qb.RemoveTable(table))
			down = append(down, dbDef.Statement)
		case hasCode && hasDB:
			if buildText(codeDef.Statement) != buildText(dbDef.Statement) {
				up = append(up, codeDef.Statement)
				down = append(down, dbDef.Statement)
			}
		}
	}

	// Access
	allAccess := unionKeys(code.Access, db.Access)
	for _, name := range allAccess {
		codeDef, hasCode := code.Access[name]
		dbDef, hasDB := db.Access[name]
		switch {
		case hasCode && !hasDB:
			up = append(up, codeDef.Statement)
			down = append(down, removeAccessStmt(codeDef))
		case !hasCode && hasDB:
			up = append(up, removeAccessStmt(dbDef))
			down = append(down, dbDef.Statement)
		case hasCode && hasDB:
			if buildText(codeDef.Statement) != buildText(dbDef.Statement) {
				up = append(up, codeDef.Statement)
				down = append(down, dbDef.Statement)
			}
		}
	}

	// Fields / Indexes / Events
	up = append(up, diffTableResources(code.Fields, db.Fields, opts, true, true, func(table, name string) qb.Statement {
		return qb.RemoveField(name).OnTableName(table)
	})...)
	down = append(down, diffTableResources(db.Fields, code.Fields, opts, true, true, func(table, name string) qb.Statement {
		return qb.RemoveField(name).OnTableName(table)
	})...)

	up = append(up, diffTableResources(code.Indexes, db.Indexes, opts, true, false, func(table, name string) qb.Statement {
		return qb.RemoveIndex(name).OnTableName(table)
	})...)
	down = append(down, diffTableResources(db.Indexes, code.Indexes, opts, true, false, func(table, name string) qb.Statement {
		return qb.RemoveIndex(name).OnTableName(table)
	})...)

	up = append(up, diffTableResources(code.Events, db.Events, opts, true, false, func(table, name string) qb.Statement {
		return qb.RemoveEvent(name).OnTableName(table)
	})...)
	down = append(down, diffTableResources(db.Events, code.Events, opts, true, false, func(table, name string) qb.Statement {
		return qb.RemoveEvent(name).OnTableName(table)
	})...)

	return up, down
}

func removeAccessStmt(def Definition) qb.Statement {
	stmt := qb.RemoveAccess(def.Name)
	switch strings.ToLower(def.Scope) {
	case "namespace":
		stmt.OnNamespace()
	case "database":
		stmt.OnDatabase()
	}
	return stmt
}

func diffTableResources(current map[string]map[string]Definition, previous map[string]map[string]Definition, opts DiffOptions, allowRename bool, allowCopy bool, removeFn func(table, name string) qb.Statement) []qb.Statement {
	var out []qb.Statement
	tables := unionKeysMap(current, previous)
	for _, table := range tables {
		cur := current[table]
		prev := previous[table]

		added := keysOnly(cur, prev)
		removed := keysOnly(prev, cur)

		if allowRename && len(added) == 1 && len(removed) == 1 {
			from := removed[0]
			to := added[0]
			switch pickRenameStrategy(opts) {
			case renameStrategyRename:
				if allowCopy {
					out = append(out, renameCopyStatement(table, from, to, opts.RenameExpr))
				}
				out = append(out, removeFn(table, from))
				out = append(out, cur[to].Statement)
				continue
			case renameStrategyDelete:
				out = append(out, cur[to].Statement)
				out = append(out, removeFn(table, from))
				continue
			case renameStrategyKeep:
				out = append(out, cur[to].Statement)
				continue
			case renameStrategyPrompt:
				if !opts.Force && opts.Prompter != nil && opts.Prompter.ConfirmRename(table, from, to) {
					if allowCopy {
						out = append(out, renameCopyStatement(table, from, to, opts.RenameExpr))
					}
					out = append(out, removeFn(table, from))
					out = append(out, cur[to].Statement)
					continue
				}
			}
		}

		for _, name := range added {
			out = append(out, cur[name].Statement)
		}
		for _, name := range removed {
			out = append(out, removeFn(table, name))
		}

		for name, curDef := range cur {
			if prevDef, ok := prev[name]; ok {
				if buildText(curDef.Statement) != buildText(prevDef.Statement) {
					out = append(out, curDef.Statement)
				}
			}
		}
	}
	return out
}

func renameCopyStatement(table, from, to, expr string) qb.Statement {
	if strings.TrimSpace(expr) == "" {
		expr = "{old}"
	}
	repl := strings.NewReplacer(
		"{old}", from,
		"{new}", to,
		"{table}", table,
	)
	expr = repl.Replace(expr)
	return qb.Update(qb.T(table)).Set(qb.Set(qb.I(to), qb.Raw(expr)))
}

type renameStrategy string

const (
	renameStrategyPrompt renameStrategy = "prompt"
	renameStrategyRename renameStrategy = "rename"
	renameStrategyDelete renameStrategy = "delete"
	renameStrategyKeep   renameStrategy = "keep"
)

func pickRenameStrategy(opts DiffOptions) renameStrategy {
	if opts.Force && (opts.RenameStrategy == "" || strings.ToLower(opts.RenameStrategy) == string(renameStrategyPrompt)) {
		return renameStrategyDelete
	}
	switch strings.ToLower(opts.RenameStrategy) {
	case string(renameStrategyRename):
		return renameStrategyRename
	case string(renameStrategyDelete):
		return renameStrategyDelete
	case string(renameStrategyKeep):
		return renameStrategyKeep
	default:
		return renameStrategyPrompt
	}
}

func buildText(stmt qb.Statement) string {
	if stmt == nil {
		return ""
	}
	return qb.Build(stmt).Text
}

func unionKeys(a, b map[string]Definition) []string {
	seen := map[string]struct{}{}
	for k := range a {
		seen[k] = struct{}{}
	}
	for k := range b {
		seen[k] = struct{}{}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func unionKeysMap(a, b map[string]map[string]Definition) []string {
	seen := map[string]struct{}{}
	for k := range a {
		seen[k] = struct{}{}
	}
	for k := range b {
		seen[k] = struct{}{}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func keysOnly(a, b map[string]Definition) []string {
	if a == nil {
		return nil
	}
	var out []string
	for k := range a {
		if b == nil {
			out = append(out, k)
			continue
		}
		if _, ok := b[k]; !ok {
			out = append(out, k)
		}
	}
	sort.Strings(out)
	return out
}

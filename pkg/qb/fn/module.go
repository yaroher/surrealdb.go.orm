package qbfn

import "github.com/yaroher/surrealdb.go.orm/pkg/qb"

// Module represents a SurrealQL function namespace (e.g. time::, math::, array::).
type Module struct {
	Prefix string
}

// Call builds a namespaced function call.
func (m Module) Call(name string, args ...any) qb.Expr[any] {
	return qb.Fn(m.Prefix+"::"+name, toNodes(args)...)
}

// Call1 builds a namespaced function with one arg.
func (m Module) Call1(name string, arg any) qb.Expr[any] {
	return m.Call(name, arg)
}

// Call2 builds a namespaced function with two args.
func (m Module) Call2(name string, a1 any, a2 any) qb.Expr[any] {
	return m.Call(name, a1, a2)
}

// Call3 builds a namespaced function with three args.
func (m Module) Call3(name string, a1 any, a2 any, a3 any) qb.Expr[any] {
	return m.Call(name, a1, a2, a3)
}

// Fn builds a function call without a namespace.
func Fn(name string, args ...any) qb.Expr[any] {
	return qb.Fn(name, toNodes(args)...)
}

// Count builds count(). If args are empty, it renders count().
func Count(args ...any) qb.Expr[any] {
	return Fn("count", args...)
}

// Sleep builds sleep(duration).
func Sleep(duration any) qb.Expr[any] {
	return Fn("sleep", duration)
}

// Rand builds rand().
func Rand() qb.Expr[any] {
	return Fn("rand")
}

// Namespaced modules.
var (
	Array     = Module{Prefix: "array"}
	Crypto    = Module{Prefix: "crypto"}
	Duration  = Module{Prefix: "duration"}
	Geo       = Module{Prefix: "geo"}
	HTTP      = Module{Prefix: "http"}
	Math      = Module{Prefix: "math"}
	Meta      = Module{Prefix: "meta"}
	Parse     = Module{Prefix: "parse"}
	RandNS    = Module{Prefix: "rand"}
	Script    = Module{Prefix: "script"}
	Scripting = Module{Prefix: "scripting"}
	Search    = Module{Prefix: "search"}
	Session   = Module{Prefix: "session"}
	String    = Module{Prefix: "string"}
	Time      = Module{Prefix: "time"}
	Type      = Module{Prefix: "type"}
	Vector    = Module{Prefix: "vector"}
)

func toNodes(args []any) []qb.Node {
	out := make([]qb.Node, 0, len(args))
	for _, arg := range args {
		if n, ok := arg.(qb.Node); ok {
			out = append(out, n)
			continue
		}
		out = append(out, qb.V(arg))
	}
	return out
}

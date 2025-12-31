package qb

import (
	"fmt"
	"strings"
)

// Builder renders SurrealQL and collects bound parameters.
type Builder struct {
	sb      strings.Builder
	args    map[string]any
	counter int
}

func NewBuilder() *Builder {
	return &Builder{args: map[string]any{}}
}

func (b *Builder) Write(s string) {
	b.sb.WriteString(s)
}

func (b *Builder) Arg(value any) string {
	b.counter++
	name := fmt.Sprintf("p%d", b.counter)
	b.args[name] = value
	return "$" + name
}

func (b *Builder) Bind(name string, value any) {
	if b.args == nil {
		b.args = map[string]any{}
	}
	if _, exists := b.args[name]; !exists {
		b.args[name] = value
	}
}

func (b *Builder) String() string {
	return b.sb.String()
}

func (b *Builder) Args() map[string]any {
	return b.args
}

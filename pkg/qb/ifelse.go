package qb

// IfBuilder builds IF / ELSE IF / ELSE statements.
type IfBuilder struct {
	branches []ifBranch
	elseBody Node
}

type ifBranch struct {
	Cond Condition
	Body Node
}

// If starts an IF statement.
func If(cond Condition) *IfBuilder {
	return &IfBuilder{branches: []ifBranch{{Cond: cond}}}
}

// Then sets the body for the current branch.
func (i *IfBuilder) Then(body Node) *IfBuilder {
	if len(i.branches) == 0 {
		return i
	}
	i.branches[len(i.branches)-1].Body = body
	return i
}

// ElseIf adds an ELSE IF branch.
func (i *IfBuilder) ElseIf(cond Condition) *IfBuilder {
	i.branches = append(i.branches, ifBranch{Cond: cond})
	return i
}

// Else sets the ELSE body.
func (i *IfBuilder) Else(body Node) *IfBuilder {
	i.elseBody = body
	return i
}

func (i *IfBuilder) build(b *Builder) {
	for idx, br := range i.branches {
		if idx == 0 {
			b.Write("IF ")
		} else {
			b.Write(" ELSE IF ")
		}
		br.Cond.build(b)
		b.Write(" ")
		BlockOf(br.Body).build(b)
	}
	if i.elseBody != nil {
		b.Write(" ELSE ")
		BlockOf(i.elseBody).build(b)
	}
	b.Write(" END")
}

func (i *IfBuilder) Build() Query {
	return Build(i)
}

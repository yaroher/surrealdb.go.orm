package qb

// SleepStatement builds SLEEP statements.
type SleepStatement struct {
	Duration Node
}

func Sleep(duration any) *SleepStatement {
	return &SleepStatement{Duration: ensureValueNode(duration)}
}

func (s *SleepStatement) build(b *Builder) {
	b.Write("SLEEP ")
	if s.Duration != nil {
		s.Duration.build(b)
	}
}

func (s *SleepStatement) Build() Query {
	return Build(s)
}

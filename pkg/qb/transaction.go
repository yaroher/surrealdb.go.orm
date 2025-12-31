package qb

// BeginTransaction starts a transaction.
func BeginTransaction() Statement {
	return rawKeywordStatement("BEGIN TRANSACTION")
}

// CommitTransaction commits a transaction.
func CommitTransaction() Statement {
	return rawKeywordStatement("COMMIT TRANSACTION")
}

// CancelTransaction cancels a transaction.
func CancelTransaction() Statement {
	return rawKeywordStatement("CANCEL TRANSACTION")
}

type rawKeywordStatement string

func (r rawKeywordStatement) build(b *Builder) {
	b.Write(string(r))
}

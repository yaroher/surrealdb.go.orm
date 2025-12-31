package codegen

import (
	"strings"

	"github.com/yaroher/surrealdb.go.orm/pkg/migrator"
	"github.com/yaroher/surrealdb.go.orm/pkg/qb"
)

// BuildResourceSet builds migrator.ResourceSet from parsed models.
func BuildResourceSet(models []Model) migrator.ResourceSet {
	res := migrator.NewResourceSet()
	for _, model := range models {
		if model.Kind == "access" && model.Access.Name != "" {
			addAccessResource(&res, model.Access)
			continue
		}
		table := qb.DefineTableName(model.Table)
		if model.Drop {
			table.DropTable()
		}
		if model.SchemaFull {
			table.SchemaFull()
		}
		if model.SchemaLess {
			table.SchemaLess()
		}
		applyPermissionsTable(table, model.Permissions)
		res.AddTable(model.Table, table)

		for _, field := range model.Fields {
			fieldStmt := qb.DefineFieldName(field.DBName, model.Table)
			if t := inferSurrealType(field.Type, field, model); t != "" {
				fieldStmt.Type(t)
			}
			if field.ValueExpr != "" {
				fieldStmt.ValueExpr(qb.Raw(field.ValueExpr))
			}
			if field.AssertExpr != "" {
				fieldStmt.Assert(qb.RawCond(field.AssertExpr))
			}
			if field.DefaultExpr != "" {
				fieldStmt.DefaultExpr(qb.Raw(field.DefaultExpr))
			}
			applyPermissionsField(fieldStmt, field.Permissions)
			res.AddField(model.Table, field.DBName, fieldStmt)
		}

		if model.Kind == "edge" && model.EdgeIn != "" && model.EdgeOut != "" {
			res.AddField(model.Table, "in", qb.DefineFieldName("in", model.Table).Type("record<"+model.EdgeIn+">"))
			res.AddField(model.Table, "out", qb.DefineFieldName("out", model.Table).Type("record<"+model.EdgeOut+">"))
		}
	}
	return res
}

func addAccessResource(res *migrator.ResourceSet, ac AccessConfig) {
	stmt := qb.DefineAccess(ac.Name)
	if ac.Overwrite {
		stmt.Overwrite()
	} else if ac.IfNotExists {
		stmt.IfNotExists()
	}
	switch strings.ToLower(ac.Scope) {
	case "namespace":
		stmt.OnNamespace()
	case "database":
		stmt.OnDatabase()
	}
	switch strings.ToLower(ac.Type) {
	case "jwt":
		stmt.TypeJWT()
		if ac.URL != "" {
			stmt.JWTURL(ac.URL)
		} else if ac.Algorithm != "" && ac.Key != "" {
			stmt.JWTAlgorithmKey(ac.Algorithm, ac.Key)
		}
	case "record":
		stmt.TypeRecord()
		if ac.RecordSignup != "" {
			stmt.Signup(qb.Raw(ac.RecordSignup))
		}
		if ac.RecordSignin != "" {
			stmt.Signin(qb.Raw(ac.RecordSignin))
		}
		if ac.RecordIssuer != "" {
			stmt.RecordIssuerKey(ac.RecordIssuer)
		}
		if ac.URL != "" {
			stmt.RecordJWTURL(ac.URL)
		} else if ac.Algorithm != "" && ac.Key != "" {
			stmt.RecordJWTAlgorithmKey(ac.Algorithm, ac.Key)
		}
		if ac.RecordRefresh {
			stmt.WithRefresh()
		}
	case "bearer":
		stmt.TypeBearer()
		if strings.ToLower(ac.Algorithm) == "record" {
			stmt.BearerForRecord()
		} else if strings.ToLower(ac.Algorithm) == "user" {
			stmt.BearerForUser()
		}
	}
	if ac.Authenticate != "" {
		stmt.Authenticate(qb.Raw(ac.Authenticate))
	}
	if ac.DurationGrant != "" {
		stmt.DurationGrant(qb.Raw(ac.DurationGrant))
	}
	if ac.DurationToken != "" {
		stmt.DurationToken(qb.Raw(ac.DurationToken))
	}
	if ac.DurationSession != "" {
		stmt.DurationSession(qb.Raw(ac.DurationSession))
	}
	res.AddAccess(ac.Name, ac.Scope, stmt)

	if ac.GrantSubject != "" || ac.GrantToken != "" || ac.GrantDuration != "" {
		grant := qb.AccessGrant(ac.Name)
		switch strings.ToLower(ac.Scope) {
		case "namespace":
			grant.OnNamespace()
		case "database":
			grant.OnDatabase()
		}
		if ac.GrantSubject != "" {
			grant.Subject(qb.Raw(ac.GrantSubject))
		}
		if ac.GrantToken != "" {
			grant.Token(qb.Raw(ac.GrantToken))
		}
		if ac.GrantDuration != "" {
			grant.Duration(qb.Raw(ac.GrantDuration))
		}
		res.AddAccessGrant(ac.Name, grant)
	}
}

func applyPermissionsTable(stmt *qb.DefineTableStatement, permissions string) {
	applyPermissions(permissions, func(n qb.Node) {
		stmt.PermissionsFor(n)
	}, func() { stmt.PermissionsFull() }, func() { stmt.PermissionsNone() })
}

func applyPermissionsField(stmt *qb.DefineFieldStatement, permissions string) {
	applyPermissions(permissions, func(n qb.Node) {
		stmt.PermissionsFor(n)
	}, func() { stmt.PermissionsFull() }, func() { stmt.PermissionsNone() })
}

func applyPermissions(permissions string, forFn func(qb.Node), fullFn func(), noneFn func()) {
	p := strings.TrimSpace(permissions)
	if p == "" {
		return
	}
	switch strings.ToLower(p) {
	case "full":
		fullFn()
	case "none":
		noneFn()
	default:
		forFn(qb.Raw(p))
	}
}

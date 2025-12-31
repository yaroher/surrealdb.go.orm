package qb

// AccessScope defines where an access method is defined.
type AccessScope string

const (
	AccessRoot      AccessScope = "ROOT"
	AccessNamespace AccessScope = "NAMESPACE"
	AccessDatabase  AccessScope = "DATABASE"
)

// DefineAccessStatement builds DEFINE ACCESS.
type DefineAccessStatement struct {
	Name            Node
	Scope           AccessScope
	overwrite       bool
	ifNotExists     bool
	AccessType      accessType
	authenticate    Node
	durationGrant   Node
	durationToken   Node
	durationSession Node
}

type accessType interface {
	build(*Builder)
}

type accessJWT struct {
	Algorithm string
	Key       Node
	URL       Node
}

type accessRecord struct {
	Signup      Node
	Signin      Node
	JWT         *accessJWT
	IssuerKey   Node
	WithRefresh bool
}

type accessBearer struct {
	For string
}

// DefineAccess starts a DEFINE ACCESS statement.
func DefineAccess(name string) *DefineAccessStatement {
	return &DefineAccessStatement{Name: Ident{Name: name}}
}

func DefineAccessExpr(expr Node) *DefineAccessStatement {
	return &DefineAccessStatement{Name: expr}
}

func (d *DefineAccessStatement) Overwrite() *DefineAccessStatement {
	d.overwrite = true
	d.ifNotExists = false
	return d
}

func (d *DefineAccessStatement) IfNotExists() *DefineAccessStatement {
	d.ifNotExists = true
	d.overwrite = false
	return d
}

func (d *DefineAccessStatement) OnRoot() *DefineAccessStatement {
	d.Scope = AccessRoot
	return d
}

func (d *DefineAccessStatement) OnNamespace() *DefineAccessStatement {
	d.Scope = AccessNamespace
	return d
}

func (d *DefineAccessStatement) OnDatabase() *DefineAccessStatement {
	d.Scope = AccessDatabase
	return d
}

// TypeJWT sets TYPE JWT.
func (d *DefineAccessStatement) TypeJWT() *DefineAccessStatement {
	d.AccessType = &accessJWT{}
	return d
}

// JWTAlgorithmKey sets JWT ALGORITHM <alg> KEY <key>.
func (d *DefineAccessStatement) JWTAlgorithmKey(alg string, key any) *DefineAccessStatement {
	jwt := ensureJWT(d)
	jwt.Algorithm = alg
	jwt.Key = ensureValueNode(key)
	jwt.URL = nil
	return d
}

// JWTURL sets JWT URL <url>.
func (d *DefineAccessStatement) JWTURL(url any) *DefineAccessStatement {
	jwt := ensureJWT(d)
	jwt.URL = ensureValueNode(url)
	jwt.Key = nil
	jwt.Algorithm = ""
	return d
}

// TypeRecord sets TYPE RECORD.
func (d *DefineAccessStatement) TypeRecord() *DefineAccessStatement {
	d.AccessType = &accessRecord{}
	return d
}

func (d *DefineAccessStatement) Signup(expr Node) *DefineAccessStatement {
	rec := ensureRecord(d)
	rec.Signup = expr
	return d
}

func (d *DefineAccessStatement) Signin(expr Node) *DefineAccessStatement {
	rec := ensureRecord(d)
	rec.Signin = expr
	return d
}

// RecordWithJWT enables WITH JWT and configures JWT options.
func (d *DefineAccessStatement) RecordWithJWT() *DefineAccessStatement {
	rec := ensureRecord(d)
	if rec.JWT == nil {
		rec.JWT = &accessJWT{}
	}
	return d
}

func (d *DefineAccessStatement) RecordJWTAlgorithmKey(alg string, key any) *DefineAccessStatement {
	rec := ensureRecord(d)
	if rec.JWT == nil {
		rec.JWT = &accessJWT{}
	}
	rec.JWT.Algorithm = alg
	rec.JWT.Key = ensureValueNode(key)
	rec.JWT.URL = nil
	return d
}

func (d *DefineAccessStatement) RecordJWTURL(url any) *DefineAccessStatement {
	rec := ensureRecord(d)
	if rec.JWT == nil {
		rec.JWT = &accessJWT{}
	}
	rec.JWT.URL = ensureValueNode(url)
	rec.JWT.Key = nil
	rec.JWT.Algorithm = ""
	return d
}

func (d *DefineAccessStatement) RecordIssuerKey(key any) *DefineAccessStatement {
	rec := ensureRecord(d)
	rec.IssuerKey = ensureValueNode(key)
	return d
}

func (d *DefineAccessStatement) WithRefresh() *DefineAccessStatement {
	rec := ensureRecord(d)
	rec.WithRefresh = true
	return d
}

// TypeBearer sets TYPE BEARER.
func (d *DefineAccessStatement) TypeBearer() *DefineAccessStatement {
	d.AccessType = &accessBearer{}
	return d
}

// BearerForUser sets TYPE BEARER FOR USER.
func (d *DefineAccessStatement) BearerForUser() *DefineAccessStatement {
	b := ensureBearer(d)
	b.For = "USER"
	return d
}

// BearerForRecord sets TYPE BEARER FOR RECORD.
func (d *DefineAccessStatement) BearerForRecord() *DefineAccessStatement {
	b := ensureBearer(d)
	b.For = "RECORD"
	return d
}

func (d *DefineAccessStatement) Authenticate(expr Node) *DefineAccessStatement {
	d.authenticate = expr
	return d
}

func (d *DefineAccessStatement) DurationToken(expr Node) *DefineAccessStatement {
	d.durationToken = expr
	return d
}

func (d *DefineAccessStatement) DurationSession(expr Node) *DefineAccessStatement {
	d.durationSession = expr
	return d
}

func (d *DefineAccessStatement) DurationGrant(expr Node) *DefineAccessStatement {
	d.durationGrant = expr
	return d
}

func (d *DefineAccessStatement) build(b *Builder) {
	b.Write("DEFINE ACCESS ")
	if d.overwrite {
		b.Write("OVERWRITE ")
	} else if d.ifNotExists {
		b.Write("IF NOT EXISTS ")
	}
	d.Name.build(b)
	if d.Scope != "" {
		b.Write(" ON ")
		b.Write(string(d.Scope))
	}
	if d.AccessType != nil {
		b.Write(" TYPE ")
		d.AccessType.build(b)
	}
	if d.authenticate != nil {
		b.Write(" AUTHENTICATE ")
		d.authenticate.build(b)
	}
	if d.durationGrant != nil || d.durationToken != nil || d.durationSession != nil {
		b.Write(" DURATION")
		parts := 0
		if d.durationGrant != nil {
			b.Write(" FOR GRANT ")
			d.durationGrant.build(b)
			parts++
		}
		if d.durationToken != nil {
			if parts > 0 {
				b.Write(",")
			}
			b.Write(" FOR TOKEN ")
			d.durationToken.build(b)
			parts++
		}
		if d.durationSession != nil {
			if parts > 0 {
				b.Write(",")
			}
			b.Write(" FOR SESSION ")
			d.durationSession.build(b)
		}
	}
}

func (d *DefineAccessStatement) Build() Query {
	return Build(d)
}

func (j *accessJWT) build(b *Builder) {
	b.Write("JWT")
	if j.URL != nil {
		b.Write(" URL ")
		j.URL.build(b)
		return
	}
	if j.Algorithm != "" && j.Key != nil {
		b.Write(" ALGORITHM ")
		b.Write(j.Algorithm)
		b.Write(" KEY ")
		j.Key.build(b)
	}
}

func (r *accessRecord) build(b *Builder) {
	b.Write("RECORD")
	if r.Signup != nil {
		b.Write(" SIGNUP ")
		BlockOf(r.Signup).build(b)
	}
	if r.Signin != nil {
		b.Write(" SIGNIN ")
		BlockOf(r.Signin).build(b)
	}
	if r.JWT != nil {
		b.Write(" WITH JWT")
		if r.JWT.URL != nil {
			b.Write(" URL ")
			r.JWT.URL.build(b)
		} else if r.JWT.Algorithm != "" && r.JWT.Key != nil {
			b.Write(" ALGORITHM ")
			b.Write(r.JWT.Algorithm)
			b.Write(" KEY ")
			r.JWT.Key.build(b)
		}
		if r.IssuerKey != nil {
			b.Write(" WITH ISSUER KEY ")
			r.IssuerKey.build(b)
		}
	}
	if r.WithRefresh {
		b.Write(" WITH REFRESH")
	}
}

func (btype *accessBearer) build(b *Builder) {
	b.Write("BEARER")
	if btype.For != "" {
		b.Write(" FOR ")
		b.Write(btype.For)
	}
}

func ensureJWT(d *DefineAccessStatement) *accessJWT {
	jwt, ok := d.AccessType.(*accessJWT)
	if !ok || jwt == nil {
		jwt = &accessJWT{}
		d.AccessType = jwt
	}
	return jwt
}

func ensureRecord(d *DefineAccessStatement) *accessRecord {
	rec, ok := d.AccessType.(*accessRecord)
	if !ok || rec == nil {
		rec = &accessRecord{}
		d.AccessType = rec
	}
	return rec
}

func ensureBearer(d *DefineAccessStatement) *accessBearer {
	b, ok := d.AccessType.(*accessBearer)
	if !ok || b == nil {
		b = &accessBearer{}
		d.AccessType = b
	}
	return b
}

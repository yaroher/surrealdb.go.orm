package codegen

type Options struct {
	Dirs []string
}

type Package struct {
	Name    string
	Models  []Model
	Imports map[string]string
}

type Model struct {
	Name        string
	Kind        string
	Table       string
	Fields      []Field
	RenameAll   string
	EdgeIn      string
	EdgeOut     string
	SchemaFull  bool
	SchemaLess  bool
	Drop        bool
	Permissions string
	Access      AccessConfig
}

type Field struct {
	Name        string
	Type        string
	DBName      string
	Annotations []Annotation
	TypeHint    string
	ValueExpr   string
	AssertExpr  string
	DefaultExpr string
	Permissions string
	LinkOne     string
	LinkMany    string
	LinkSelf    string
}

type AccessConfig struct {
	Name            string
	Scope           string
	Overwrite       bool
	IfNotExists     bool
	Type            string
	Algorithm       string
	Key             string
	URL             string
	RecordSignup    string
	RecordSignin    string
	RecordIssuer    string
	RecordRefresh   bool
	Authenticate    string
	DurationGrant   string
	DurationToken   string
	DurationSession string
	GrantSubject    string
	GrantToken      string
	GrantDuration   string
}

// Annotation describes a parsed orm directive.
type Annotation struct {
	Kind string
	Args map[string]string
}

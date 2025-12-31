package codegen

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io/fs"
	"path/filepath"
	"strings"
)

func ParseDir(dir string) (Package, error) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, func(info fs.FileInfo) bool {
		name := info.Name()
		if strings.HasSuffix(name, "_test.go") {
			return false
		}
		if strings.HasSuffix(name, "_orm_gen.go") {
			return false
		}
		return strings.HasSuffix(name, ".go")
	}, parser.ParseComments)
	if err != nil {
		return Package{}, err
	}

	var pkgName string
	var models []Model
	usedImports := map[string]string{}
	for name, pkg := range pkgs {
		pkgName = name
		for _, file := range pkg.Files {
			fileImports := collectImports(file)
			fileModels := parseFile(file, fset, fileImports, usedImports)
			models = append(models, fileModels...)
		}
		break
	}
	if pkgName == "" {
		return Package{}, fmt.Errorf("no package found in %s", dir)
	}
	return Package{Name: pkgName, Models: models, Imports: usedImports}, nil
}

func parseFile(file *ast.File, fset *token.FileSet, imports map[string]string, usedImports map[string]string) []Model {
	var models []Model
	for _, decl := range file.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok || gen.Tok != token.TYPE {
			continue
		}
		for _, spec := range gen.Specs {
			tspec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			structType, ok := tspec.Type.(*ast.StructType)
			if !ok {
				continue
			}
			anns := collectAnnotations(gen.Doc, tspec.Doc)
			modelAnn, ok := findModelAnnotation(anns)
			if !ok {
				continue
			}

			rest := buildModel(tspec.Name.Name, modelAnn, structType, fset, imports, usedImports)
			models = append(models, rest)
		}
	}
	return models
}

func buildModel(name string, ann Annotation, st *ast.StructType, fset *token.FileSet, imports map[string]string, usedImports map[string]string) Model {
	table := ann.Args["table"]
	if table == "" {
		table = toSnake(name)
	}

	renameAll := ann.Args["rename_all"]
	if renameAll == "" {
		renameAll = ann.Args["renameAll"]
	}
	if renameAll == "" {
		renameAll = "camelCase"
	}

	model := Model{
		Name:        name,
		Kind:        ann.Kind,
		Table:       table,
		RenameAll:   renameAll,
		EdgeIn:      normalizeRef(ann.Args["in"]),
		EdgeOut:     normalizeRef(ann.Args["out"]),
		SchemaFull:  ann.Args["schemafull"] == "true",
		SchemaLess:  ann.Args["schemaless"] == "true",
		Drop:        ann.Args["drop"] == "true",
		Permissions: ann.Args["permissions"],
		Access:      parseAccessConfig(ann),
	}

	for _, field := range st.Fields.List {
		if len(field.Names) == 0 {
			continue
		}
		fieldAnns := collectAnnotations(field.Doc, field.Comment)
		fieldMeta := mergeAnnotationArgs(fieldAnns, "field")
		markTypeImports(field.Type, imports, usedImports)
		fieldType := exprString(field.Type, fset)
		for _, name := range field.Names {
			dbName := applyRenameAll(name.Name, renameAll)
			if override := fieldNameOverride(fieldAnns); override != "" {
				dbName = override
			}
			model.Fields = append(model.Fields, Field{
				Name:        name.Name,
				Type:        fieldType,
				DBName:      dbName,
				Annotations: fieldAnns,
				TypeHint:    fieldMeta["type"],
				ValueExpr:   fieldMeta["value"],
				AssertExpr:  fieldMeta["assert"],
				DefaultExpr: fieldMeta["default"],
				Permissions: fieldMeta["permissions"],
				LinkOne:     fieldMeta["link_one"],
				LinkMany:    fieldMeta["link_many"],
				LinkSelf:    fieldMeta["link_self"],
			})
		}
	}

	return model
}

func parseAccessConfig(ann Annotation) AccessConfig {
	ac := AccessConfig{}
	if ann.Kind != "access" {
		return ac
	}
	ac.Name = ann.Args["name"]
	ac.Scope = ann.Args["on"]
	ac.Type = ann.Args["type"]
	ac.Algorithm = ann.Args["alg"]
	ac.Key = ann.Args["key"]
	ac.URL = ann.Args["url"]
	ac.RecordSignup = ann.Args["signup"]
	ac.RecordSignin = ann.Args["signin"]
	ac.RecordIssuer = ann.Args["issuer"]
	ac.RecordRefresh = ann.Args["refresh"] == "true"
	ac.Authenticate = ann.Args["authenticate"]
	ac.DurationGrant = ann.Args["duration_grant"]
	ac.DurationToken = ann.Args["duration_token"]
	ac.DurationSession = ann.Args["duration_session"]
	ac.GrantSubject = ann.Args["grant_subject"]
	ac.GrantToken = ann.Args["grant_token"]
	ac.GrantDuration = ann.Args["grant_duration"]
	ac.Overwrite = ann.Args["overwrite"] == "true"
	ac.IfNotExists = ann.Args["if_not_exists"] == "true"
	if ac.Name == "" {
		ac.Name = ann.Args["access"]
	}
	return ac
}

func fieldNameOverride(anns []Annotation) string {
	for _, ann := range anns {
		if ann.Kind == "field" {
			if v, ok := ann.Args["name"]; ok {
				return v
			}
		}
	}
	return ""
}

func mergeAnnotationArgs(anns []Annotation, kind string) map[string]string {
	out := map[string]string{}
	for _, ann := range anns {
		if ann.Kind != kind {
			continue
		}
		for k, v := range ann.Args {
			out[k] = v
		}
	}
	return out
}

func findModelAnnotation(anns []Annotation) (Annotation, bool) {
	for _, ann := range anns {
		switch ann.Kind {
		case "node", "edge", "object", "access":
			return ann, true
		}
	}
	return Annotation{}, false
}

func collectAnnotations(groups ...*ast.CommentGroup) []Annotation {
	var out []Annotation
	for _, group := range groups {
		if group == nil {
			continue
		}
		for _, line := range strings.Split(group.Text(), "\n") {
			if ann, ok := ParseAnnotation(line); ok {
				out = append(out, ann)
			}
		}
	}
	return out
}

func exprString(expr ast.Expr, fset *token.FileSet) string {
	var buf bytes.Buffer
	_ = printer.Fprint(&buf, fset, expr)
	return buf.String()
}

func outFileName(dir string) string {
	return filepath.Join(dir, "orm_gen.go")
}

func normalizeRef(name string) string {
	if name == "" {
		return name
	}
	if strings.Contains(name, ":") || strings.Contains(name, "/") || strings.Contains(name, ".") {
		return name
	}
	return toSnake(name)
}

func NormalizeRef(name string) string {
	return normalizeRef(name)
}

func collectImports(file *ast.File) map[string]string {
	out := map[string]string{}
	for _, imp := range file.Imports {
		path := strings.Trim(imp.Path.Value, "\"")
		name := ""
		if imp.Name != nil {
			name = imp.Name.Name
		} else {
			parts := strings.Split(path, "/")
			name = parts[len(parts)-1]
		}
		out[name] = path
	}
	return out
}

func markTypeImports(expr ast.Expr, imports map[string]string, used map[string]string) {
	ast.Inspect(expr, func(node ast.Node) bool {
		sel, ok := node.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		id, ok := sel.X.(*ast.Ident)
		if !ok {
			return true
		}
		path, ok := imports[id.Name]
		if !ok {
			return true
		}
		used[id.Name] = path
		return true
	})
}

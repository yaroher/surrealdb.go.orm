package main

import (
	"log"
	"path"
	"strings"

	"github.com/yaroher/surrealdb.go.orm/internal/codegen"
	"github.com/yaroher/surrealdb.go.orm/protosurrealorm"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

type generator struct {
	messages  map[protoreflect.FullName]*protogen.Message
	enums     map[protoreflect.FullName]*protogen.Enum
	pkgByPath map[protogen.GoImportPath]protogen.GoPackageName
}

func main() {
	opts := protogen.Options{}
	if err := opts.Run(func(gen *protogen.Plugin) error {
		g := newGenerator(gen)
		for _, file := range gen.Files {
			if !file.Generate {
				continue
			}
			models, imports := g.modelsForFile(file)
			if len(models) == 0 {
				continue
			}
			pkg := codegen.Package{
				Name:    string(file.GoPackageName),
				Models:  models,
				Imports: imports,
			}
			out, err := codegen.RenderToBytes(pkg)
			if err != nil {
				return err
			}
			filename := outputName(file)
			outFile := gen.NewGeneratedFile(filename, file.GoImportPath)
			if _, err := outFile.Write(out); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		log.Fatal(err)
	}
}

func newGenerator(gen *protogen.Plugin) *generator {
	g := &generator{
		messages:  map[protoreflect.FullName]*protogen.Message{},
		enums:     map[protoreflect.FullName]*protogen.Enum{},
		pkgByPath: map[protogen.GoImportPath]protogen.GoPackageName{},
	}
	for _, file := range gen.Files {
		if file.GoImportPath != "" {
			g.pkgByPath[file.GoImportPath] = file.GoPackageName
		}
		for _, enum := range file.Enums {
			g.enums[enum.Desc.FullName()] = enum
		}
		g.collectMessages(file.Messages)
	}
	return g
}

func (g *generator) collectMessages(messages []*protogen.Message) {
	for _, msg := range messages {
		g.messages[msg.Desc.FullName()] = msg
		for _, enum := range msg.Enums {
			g.enums[enum.Desc.FullName()] = enum
		}
		if len(msg.Messages) > 0 {
			g.collectMessages(msg.Messages)
		}
	}
}

func outputName(file *protogen.File) string {
	name := file.Desc.Path()
	name = strings.TrimSuffix(name, ".proto")
	return name + "_orm_gen.go"
}

func (g *generator) modelsForFile(file *protogen.File) ([]codegen.Model, map[string]string) {
	imports := map[string]string{}
	var models []codegen.Model
	var walk func(msg *protogen.Message)
	walk = func(msg *protogen.Message) {
		if opts, ok := modelOptions(msg); ok {
			models = append(models, g.modelFromMessage(file, msg, opts, imports))
		}
		for _, nested := range msg.Messages {
			walk(nested)
		}
	}
	for _, msg := range file.Messages {
		walk(msg)
	}
	return models, imports
}

func modelOptions(msg *protogen.Message) (*protosurrealorm.ModelOptions, bool) {
	opts, ok := msg.Desc.Options().(*descriptorpb.MessageOptions)
	if !ok || opts == nil {
		return nil, false
	}
	if !proto.HasExtension(opts, protosurrealorm.E_OrmModel) {
		return nil, false
	}
	ext, err := proto.GetExtension(opts, protosurrealorm.E_OrmModel)
	if err != nil {
		return nil, false
	}
	out, ok := ext.(*protosurrealorm.ModelOptions)
	if !ok || out == nil {
		return nil, false
	}
	return out, true
}

func fieldOptions(field *protogen.Field) *protosurrealorm.FieldOptions {
	opts, ok := field.Desc.Options().(*descriptorpb.FieldOptions)
	if !ok || opts == nil {
		return nil
	}
	if !proto.HasExtension(opts, protosurrealorm.E_OrmField) {
		return nil
	}
	ext, err := proto.GetExtension(opts, protosurrealorm.E_OrmField)
	if err != nil {
		return nil
	}
	out, ok := ext.(*protosurrealorm.FieldOptions)
	if !ok {
		return nil
	}
	return out
}

func (g *generator) modelFromMessage(file *protogen.File, msg *protogen.Message, opts *protosurrealorm.ModelOptions, imports map[string]string) codegen.Model {
	kind := strings.TrimSpace(opts.GetKind())
	if kind == "" {
		if opts.Access != nil {
			kind = "access"
		} else {
			kind = "node"
		}
	}
	table := opts.GetTable()
	if table == "" {
		table = codegen.ToSnake(msg.GoIdent.GoName)
	}
	renameAll := opts.GetRenameAll()
	if renameAll == "" {
		renameAll = "camelCase"
	}

	model := codegen.Model{
		Name:        msg.GoIdent.GoName,
		Kind:        kind,
		Table:       table,
		RenameAll:   renameAll,
		EdgeIn:      codegen.NormalizeRef(opts.GetIn()),
		EdgeOut:     codegen.NormalizeRef(opts.GetOut()),
		SchemaFull:  opts.GetSchemafull(),
		SchemaLess:  opts.GetSchemaless(),
		Drop:        opts.GetDrop(),
		Permissions: opts.GetPermissions(),
		Access:      accessFromOptions(opts.Access),
	}

	for _, field := range msg.Fields {
		fieldOpts := fieldOptions(field)
		dbName := codegen.ApplyRenameAll(field.GoName, renameAll)
		if fieldOpts != nil && fieldOpts.GetName() != "" {
			dbName = fieldOpts.GetName()
		}
		model.Fields = append(model.Fields, codegen.Field{
			Name:        field.GoName,
			Type:        g.goTypeForField(field, file, imports),
			DBName:      dbName,
			TypeHint:    getFieldString(fieldOpts, func(o *protosurrealorm.FieldOptions) string { return o.GetType() }),
			ValueExpr:   getFieldString(fieldOpts, func(o *protosurrealorm.FieldOptions) string { return o.GetValue() }),
			AssertExpr:  getFieldString(fieldOpts, func(o *protosurrealorm.FieldOptions) string { return o.GetAssert() }),
			DefaultExpr: getFieldString(fieldOpts, func(o *protosurrealorm.FieldOptions) string { return o.GetDefault() }),
			Permissions: getFieldString(fieldOpts, func(o *protosurrealorm.FieldOptions) string { return o.GetPermissions() }),
			LinkOne:     getFieldString(fieldOpts, func(o *protosurrealorm.FieldOptions) string { return o.GetLinkOne() }),
			LinkMany:    getFieldString(fieldOpts, func(o *protosurrealorm.FieldOptions) string { return o.GetLinkMany() }),
			LinkSelf:    getFieldString(fieldOpts, func(o *protosurrealorm.FieldOptions) string { return o.GetLinkSelf() }),
		})
	}

	return model
}

func getFieldString(opts *protosurrealorm.FieldOptions, get func(*protosurrealorm.FieldOptions) string) string {
	if opts == nil {
		return ""
	}
	return get(opts)
}

func accessFromOptions(opts *protosurrealorm.AccessOptions) codegen.AccessConfig {
	if opts == nil {
		return codegen.AccessConfig{}
	}
	return codegen.AccessConfig{
		Name:            opts.GetName(),
		Scope:           opts.GetOn(),
		Type:            opts.GetType(),
		Algorithm:       opts.GetAlg(),
		Key:             opts.GetKey(),
		URL:             opts.GetUrl(),
		RecordSignup:    opts.GetSignup(),
		RecordSignin:    opts.GetSignin(),
		RecordIssuer:    opts.GetIssuer(),
		RecordRefresh:   opts.GetRefresh(),
		Authenticate:    opts.GetAuthenticate(),
		DurationGrant:   opts.GetDurationGrant(),
		DurationToken:   opts.GetDurationToken(),
		DurationSession: opts.GetDurationSession(),
		GrantSubject:    opts.GetGrantSubject(),
		GrantToken:      opts.GetGrantToken(),
		GrantDuration:   opts.GetGrantDuration(),
		Overwrite:       opts.GetOverwrite(),
		IfNotExists:     opts.GetIfNotExists(),
	}
}

func (g *generator) goTypeForField(field *protogen.Field, file *protogen.File, imports map[string]string) string {
	desc := field.Desc
	if desc.IsMap() {
		keyType := g.goTypeForDescriptor(desc.MapKey(), file, imports)
		valType := g.goTypeForDescriptor(desc.MapValue(), file, imports)
		return "map[" + keyType + "]" + valType
	}

	base := g.goTypeForDescriptor(desc, file, imports)
	if desc.Cardinality() == protoreflect.Repeated {
		return "[]" + base
	}
	return base
}

func (g *generator) goTypeForDescriptor(desc protoreflect.FieldDescriptor, file *protogen.File, imports map[string]string) string {
	switch desc.Kind() {
	case protoreflect.BoolKind:
		return "bool"
	case protoreflect.EnumKind:
		if enum := g.enums[desc.Enum().FullName()]; enum != nil {
			return g.qualifyIdent(enum.GoIdent, file, imports)
		}
		return string(desc.Enum().Name())
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		return "int32"
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		return "int64"
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		return "uint32"
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		return "uint64"
	case protoreflect.FloatKind:
		return "float32"
	case protoreflect.DoubleKind:
		return "float64"
	case protoreflect.StringKind:
		return "string"
	case protoreflect.BytesKind:
		return "[]byte"
	case protoreflect.MessageKind, protoreflect.GroupKind:
		if msg := g.messages[desc.Message().FullName()]; msg != nil {
			return "*" + g.qualifyIdent(msg.GoIdent, file, imports)
		}
		return "*" + string(desc.Message().Name())
	default:
		return "any"
	}
}

func (g *generator) qualifyIdent(ident protogen.GoIdent, file *protogen.File, imports map[string]string) string {
	if ident.GoImportPath == file.GoImportPath {
		return ident.GoName
	}
	pkg := g.packageName(ident.GoImportPath)
	imports[pkg] = string(ident.GoImportPath)
	return pkg + "." + ident.GoName
}

func (g *generator) packageName(pathName protogen.GoImportPath) string {
	if name, ok := g.pkgByPath[pathName]; ok && name != "" {
		return string(name)
	}
	return path.Base(string(pathName))
}

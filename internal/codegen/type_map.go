package codegen

import "strings"

func inferSurrealType(goType string, field Field, model Model) string {
	if field.TypeHint != "" {
		return field.TypeHint
	}
	if field.LinkOne != "" {
		return "record<" + normalizeRef(field.LinkOne) + ">"
	}
	if field.LinkSelf != "" {
		return "record<" + normalizeRef(field.LinkSelf) + ">"
	}
	if field.LinkMany != "" {
		return "array<record<" + normalizeRef(field.LinkMany) + ">>"
	}
	if strings.Contains(goType, "LinkOne") {
		ref := extractGenericType(goType)
		if ref != "" {
			return "record<" + normalizeRef(ref) + ">"
		}
		return "record"
	}
	if strings.Contains(goType, "LinkSelf") {
		ref := extractGenericType(goType)
		if ref != "" {
			return "record<" + normalizeRef(ref) + ">"
		}
		return "record"
	}
	if strings.Contains(goType, "LinkMany") {
		ref := extractGenericType(goType)
		if ref != "" {
			return "array<record<" + normalizeRef(ref) + ">>"
		}
		return "array"
	}
	if strings.Contains(goType, "SimpleID") || strings.Contains(goType, "ID[") {
		return "record<" + model.Table + ">"
	}
	if strings.HasPrefix(goType, "[]") {
		return "array"
	}
	if strings.HasPrefix(goType, "map[") {
		return "object"
	}
	switch goType {
	case "string":
		return "string"
	case "bool":
		return "bool"
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64", "uintptr":
		return "int"
	case "float32", "float64":
		return "float"
	case "[]byte":
		return "bytes"
	}
	if strings.HasSuffix(goType, ".Time") {
		return "datetime"
	}
	if strings.HasSuffix(goType, ".Duration") {
		return "duration"
	}
	return ""
}

func extractGenericType(typeStr string) string {
	start := strings.Index(typeStr, "[")
	end := strings.LastIndex(typeStr, "]")
	if start == -1 || end == -1 || end <= start+1 {
		return ""
	}
	inner := strings.TrimSpace(typeStr[start+1 : end])
	inner = strings.TrimPrefix(inner, "*")
	parts := strings.Split(inner, ".")
	return parts[len(parts)-1]
}
